package main

import (
	"log"
	"strconv"
)

/*
BNF Design:
	<message_type> ::= message <identifier> { <variable_declare> } |
					message <identifier> { <message_type> <variable_declare> } |
					message <identifier> { <variable_declare> <message_type> } |
					message <identifier> { <enum_type> <variable_declare> } |
					message <identifier> { <variable_declare> <enum_type> } |

	<enum_type> :: = enum <identifier> { <variable_declare> }

	<identifier> ::= "\w+"

	<variable_declare> ::= <type> <identifier> = <number>; 			|
						repeated <type> <identifier> = <number>; 	|
						option <type> <identifier> = <number>; 		|
						required <type> <identifier> = <number>; 	|

	<number> ::= "\d"

	<service_define> ::= service <identifier> {<rpc_define>}

	<rpc_define> ::= rpc <identifier> (<message_type>) returns (<message_type>) {};	|
					rpc <identifier> (stream <message_type>) returns (stream <message_type>) {};

	<type> ::= int32 | int64 | uint32 | uint64 | sint32 | sint64 | fixed32 | fixed64 | sfixed32 | sfixed64 |
			   float | double | bool | bytes | string | <message_type> | <enum_type>
*/
type DefRpc struct {
	isParamStream bool
	isRetStream   bool
	name          string
	param         string
	ret           string
}
type DefService struct {
	name    string
	rpcList []*DefRpc
}
type DefType struct {
	parentType string // "" or "xxx"
	def        string // "enum" or "message"
	name       string
	members    []*DefMember
}
type DefMember struct {
	typeName string
	name     string
	no       int
	tag      string // "" or "array"
}

type Parser struct {
	lexer         *Lexer
	v             string // "proto2" or "proto3"
	pkg           string
	types         []*DefType
	services      []*DefService
	typeCheckList []*Token
}

func newParser(lexer *Lexer) *Parser {
	p := &Parser{}
	p.lexer = lexer
	p.types = make([]*DefType, 0)
	p.services = make([]*DefService, 0)
	p.typeCheckList = make([]*Token, 0)
	return p
}

func ParseError(token *Token, m string) {
	log.Panicln(token.lineno, token.text, "ERR:", m)
}

func (p *Parser) isBasicType(tokenType int) bool {
	switch tokenType {
	case TOKEN_KEYWORD_DOUBLE, TOKEN_KEYWORD_FLOAT, TOKEN_KEYWORD_INT32,
		TOKEN_KEYWORD_INT64, TOKEN_KEYWORD_UINT32, TOKEN_KEYWORD_UINT64,
		TOKEN_KEYWORD_SINT32, TOKEN_KEYWORD_SINT64, TOKEN_KEYWORD_FIXED32,
		TOKEN_KEYWORD_FIXED64, TOKEN_KEYWORD_SFIXED32, TOKEN_KEYWORD_SFIXED64,
		TOKEN_KEYWORD_BOOL, TOKEN_KEYWORD_STRING, TOKEN_KEYWORD_BYTES:
		return true
	}
	return false
}

func (p *Parser) getDefType(name string) *DefType {
	for _, v := range p.types {
		if v.name == name {
			return v
		}
	}
	return nil
}

func (p *Parser) checkUnknownTypes() {
	for _, v := range p.typeCheckList {
		if p.isBasicType(v.tokenType) == false {
			if p.getDefType(v.text) == nil {
				ParseError(v, "unkonown type")
			}
		}
	}
}

func (p *Parser) printAll() {
	// print types
	for _, defType := range p.types {
		log.Printf("%s %s %s\n", defType.parentType, defType.def, defType.name)
		for _, defMember := range defType.members {
			log.Printf("	%s %s %s\n", defMember.tag, defMember.typeName, defMember.name)
		}
	}
	// print services
	for _, defService := range p.services {
		log.Printf("service %s\n", defService.name)
		for _, rpc := range defService.rpcList {
			if rpc.isParamStream && rpc.isRetStream {
				log.Printf("  %s (stream %s) (stream %s)\n", rpc.name, rpc.param, rpc.ret)
			} else if rpc.isParamStream {
				log.Printf("  %s (stream %s) (%s)\n", rpc.name, rpc.param, rpc.ret)
			} else if rpc.isRetStream {
				log.Printf("  %s (%s) (stream %s)\n", rpc.name, rpc.param, rpc.ret)
			} else {
				log.Printf("  %s (%s) (%s)\n", rpc.name, rpc.param, rpc.ret)
			}
		}
	}
}

func (p *Parser) checkToken(tokenType int) *Token {
	token := p.lexer.takeToken()
	if token.tokenType != tokenType {
		ParseError(token, "invalid syntax")
	}
	//log.Println("checkToken", token.lineno, token.text)
	return token
}

func (p *Parser) parse() {
	// 1. syntax token
	p.parse_syntax()

	// 2. package token
	p.parse_package()

	// 3. other tokens
	for {
		tokenType := p.lexer.nextTokenType()
		if tokenType == -1 {
			break
		}
		switch tokenType {
		case TOKEN_KEYWORD_ENUM:
			p.parse_enum("")
		case TOKEN_KEYWORD_MESSAGE:
			p.parse_message("")
		case TOKEN_KEYWORD_SERVICE:
			p.parse_service()
		default:
			ParseError(p.lexer.takeToken(), "proto error here")
		}
	}
	p.checkUnknownTypes()
	//p.printAll()
}

func (p *Parser) parse_syntax() {
	p.checkToken(TOKEN_KEYWORD_SYNTAX)   // syntax
	p.checkToken(TOKEN_ASSIGN)           // =
	p.checkToken(TOKEN_QUOTE)            // "
	syntax := p.checkToken(TOKEN_SYMBOL) // xxx
	p.checkToken(TOKEN_QUOTE)            // "
	p.checkToken(TOKEN_SEMICOLON)        // ;
	p.v = syntax.text
}

func (p *Parser) parse_package() {
	p.checkToken(TOKEN_KEYWORD_PACKAGE) // package
	pkg := p.checkToken(TOKEN_SYMBOL)   // xxx
	p.checkToken(TOKEN_SEMICOLON)       // ;
	p.pkg = pkg.text
}

func (p *Parser) parse_enum(parentType string) {
	//log.Println("parse_enum", parentType)
	p.checkToken(TOKEN_KEYWORD_ENUM)       // enum
	enumName := p.checkToken(TOKEN_SYMBOL) // xxx
	defType := &DefType{
		parentType: parentType,
		def:        "enum",
		name:       enumName.text,
		members:    make([]*DefMember, 0),
	}
	p.types = append(p.types, defType)
	p.checkToken(TOKEN_BRACE_LEFT) // {
	p.parse_enum_members(defType)
	p.checkToken(TOKEN_BRACE_RIGHT) // }
}

func (p *Parser) parse_enum_members(defType *DefType) {
	member := p.checkToken(TOKEN_SYMBOL) // xxx
	p.checkToken(TOKEN_ASSIGN)           // =
	no := p.checkToken(TOKEN_NUMBER)     // 1,2,3
	p.checkToken(TOKEN_SEMICOLON)        // ;
	num, _ := strconv.Atoi(no.text)
	defMember := &DefMember{
		typeName: "",
		name:     member.text,
		no:       num,
		tag:      "",
	}
	defType.members = append(defType.members, defMember)
	if p.lexer.nextTokenType() == TOKEN_SYMBOL {
		p.parse_enum_members(defType)
	}
}

func (p *Parser) parse_message(parentType string) {
	//log.Println("parse_message", parentType)
	p.checkToken(TOKEN_KEYWORD_MESSAGE)       // message
	messageName := p.checkToken(TOKEN_SYMBOL) // xxx
	defType := &DefType{
		parentType: parentType,
		def:        "message",
		name:       messageName.text,
		members:    make([]*DefMember, 0),
	}
	p.types = append(p.types, defType)
	p.checkToken(TOKEN_BRACE_LEFT) // {
	p.parse_message_members(defType)
	p.checkToken(TOKEN_BRACE_RIGHT) // }
}

func (p *Parser) parse_message_members(defType *DefType) {
	nextTokenType := p.lexer.nextTokenType()
	if nextTokenType == -1 {
		return
	}
	if p.v == "proto3" {
		switch nextTokenType {
		case TOKEN_KEYWORD_DOUBLE, TOKEN_KEYWORD_FLOAT, TOKEN_KEYWORD_INT32,
			TOKEN_KEYWORD_INT64, TOKEN_KEYWORD_UINT32, TOKEN_KEYWORD_UINT64,
			TOKEN_KEYWORD_SINT32, TOKEN_KEYWORD_SINT64, TOKEN_KEYWORD_FIXED32,
			TOKEN_KEYWORD_FIXED64, TOKEN_KEYWORD_SFIXED32, TOKEN_KEYWORD_SFIXED64,
			TOKEN_KEYWORD_BOOL, TOKEN_KEYWORD_STRING, TOKEN_KEYWORD_BYTES,
			TOKEN_SYMBOL, TOKEN_KEYWORD_REPEATED:
			tag := ""
			if nextTokenType == TOKEN_KEYWORD_REPEATED {
				tag = "array"
				p.checkToken(TOKEN_KEYWORD_REPEATED)
			}
			memberType := p.checkToken(p.lexer.nextTokenType()) // xxx
			p.typeCheckList = append(p.typeCheckList, memberType)
			member := p.checkToken(TOKEN_SYMBOL) // xxx
			p.checkToken(TOKEN_ASSIGN)           // =
			no := p.checkToken(TOKEN_NUMBER)     // 1,2,3
			p.checkToken(TOKEN_SEMICOLON)        // ;
			num, _ := strconv.Atoi(no.text)
			defMember := &DefMember{
				typeName: memberType.text,
				name:     member.text,
				no:       num,
				tag:      tag,
			}
			defType.members = append(defType.members, defMember)

		case TOKEN_KEYWORD_ENUM:
			p.parse_enum(defType.name)
		case TOKEN_KEYWORD_MESSAGE:
			p.parse_message(defType.name)
		// function end
		default:
			return
		}

	} else if p.v == "proto2" {

	}

	p.parse_message_members(defType)
}

func (p *Parser) parse_service() {
	//log.Println("parse_service")
	p.checkToken(TOKEN_KEYWORD_SERVICE)       // service
	serviceName := p.checkToken(TOKEN_SYMBOL) // xxx
	service := &DefService{
		name:    serviceName.text,
		rpcList: make([]*DefRpc, 0),
	}
	p.services = append(p.services, service)
	p.checkToken(TOKEN_BRACE_LEFT) // {
	p.parse_service_rpcs(service)
	p.checkToken(TOKEN_BRACE_RIGHT) // }
}

func (p *Parser) parse_service_rpcs(service *DefService) {
	if p.lexer.nextTokenType() != TOKEN_KEYWORD_RPC {
		return
	}
	p.checkToken(TOKEN_KEYWORD_RPC)       // rpc
	rpcName := p.checkToken(TOKEN_SYMBOL) //xxx
	p.checkToken(TOKEN_BRACKETS_LEFT)     // (

	isParamStream := false
	if p.lexer.nextTokenType() == TOKEN_KEYWORD_STREAM {
		p.checkToken(TOKEN_KEYWORD_STREAM) // stream
		isParamStream = true
	}
	paramType := p.checkToken(TOKEN_SYMBOL) // xxx
	p.checkToken(TOKEN_BRACKETS_RIGHT)      // )

	p.checkToken(TOKEN_KEYWORD_RETURNS) // returns

	p.checkToken(TOKEN_BRACKETS_LEFT) // (
	isRetStream := false
	if p.lexer.nextTokenType() == TOKEN_KEYWORD_STREAM {
		p.checkToken(TOKEN_KEYWORD_STREAM) // stream
		isRetStream = true
	}
	retType := p.checkToken(TOKEN_SYMBOL) // xxx
	p.checkToken(TOKEN_BRACKETS_RIGHT)    // )

	if p.lexer.nextTokenType() == TOKEN_BRACE_LEFT {
		p.checkToken(TOKEN_BRACE_LEFT)
		p.checkToken(TOKEN_BRACE_RIGHT)
	}

	p.checkToken(TOKEN_SEMICOLON) // ;

	p.typeCheckList = append(p.typeCheckList, paramType)
	p.typeCheckList = append(p.typeCheckList, retType)
	service.rpcList = append(service.rpcList, &DefRpc{
		name:          rpcName.text,
		isParamStream: isParamStream,
		isRetStream:   isRetStream,
		param:         paramType.text,
		ret:           retType.text,
	})

	p.parse_service_rpcs(service)
}
