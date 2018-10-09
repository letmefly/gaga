package main

import (
	"bufio"
	"io"
	"log"
	"regexp"
	"strings"
)

const (
	TOKEN_MIN              = iota
	TOKEN_COMMENT          // //
	TOKEN_ASSIGN           // =
	TOKEN_SEMICOLON        // ;
	TOKEN_BRACE_LEFT       // {
	TOKEN_BRACE_RIGHT      // }
	TOKEN_BRACKETS_LEFT    // (
	TOKEN_BRACKETS_RIGHT   // )
	TOKEN_QUOTE            // "
	TOKEN_KEYWORD_SERVICE  // service
	TOKEN_KEYWORD_SYNTAX   // syntax
	TOKEN_KEYWORD_PACKAGE  // package
	TOKEN_KEYWORD_RPC      // rpc
	TOKEN_KEYWORD_RETURNS  // returns
	TOKEN_KEYWORD_ENUM     // enum
	TOKEN_KEYWORD_MESSAGE  // message
	TOKEN_KEYWORD_STREAM   // stream
	TOKEN_KEYWORD_DOUBLE   // double
	TOKEN_KEYWORD_FLOAT    // float
	TOKEN_KEYWORD_INT32    // int32
	TOKEN_KEYWORD_INT64    // int64
	TOKEN_KEYWORD_UINT32   // uint32
	TOKEN_KEYWORD_UINT64   // uint64
	TOKEN_KEYWORD_SINT32   // sint32
	TOKEN_KEYWORD_SINT64   // sint64
	TOKEN_KEYWORD_FIXED32  // fixed32
	TOKEN_KEYWORD_FIXED64  // fixed64
	TOKEN_KEYWORD_SFIXED32 // sfixed32
	TOKEN_KEYWORD_SFIXED64 // sfixed64
	TOKEN_KEYWORD_BOOL     // bool
	TOKEN_KEYWORD_STRING   // string
	TOKEN_KEYWORD_BYTES    // bytes
	TOKEN_KEYWORD_OPTION   // option
	TOKEN_KEYWORD_REQUIRED // required
	TOKEN_KEYWORD_REPEATED // repeated
	TOKEN_NUMBER
	TOKEN_SYMBOL
	TOKEN_MAX
)

var token_rules = map[int]string{
	TOKEN_COMMENT:          `\s*//\s*`,
	TOKEN_ASSIGN:           `\s*=\s*`,
	TOKEN_SEMICOLON:        `\s*;\s*`,
	TOKEN_BRACE_LEFT:       `\s*{\s*`,
	TOKEN_BRACE_RIGHT:      `\s*}\s*`,
	TOKEN_BRACKETS_LEFT:    `\s*\(\s*`,
	TOKEN_BRACKETS_RIGHT:   `\s*\)\s*`,
	TOKEN_QUOTE:            `\s*"\s*`,
	TOKEN_KEYWORD_SERVICE:  `^\s*service\s+`,
	TOKEN_KEYWORD_SYNTAX:   `^\s*syntax\s+`,
	TOKEN_KEYWORD_PACKAGE:  `^\s*package\s+`,
	TOKEN_KEYWORD_RPC:      `^\s*rpc\s+`,
	TOKEN_KEYWORD_RETURNS:  `^\s*returns\s+`,
	TOKEN_KEYWORD_ENUM:     `^\s*enum\s+`,
	TOKEN_KEYWORD_MESSAGE:  `^\s*message\s+`,
	TOKEN_KEYWORD_STREAM:   `^\s*stream\s+`,
	TOKEN_KEYWORD_DOUBLE:   `^\s*double\s+`,
	TOKEN_KEYWORD_FLOAT:    `^\s*float\s+`,
	TOKEN_KEYWORD_INT32:    `^\s*int32\s+`,
	TOKEN_KEYWORD_INT64:    `^\s*int64\s+`,
	TOKEN_KEYWORD_UINT32:   `^\s*uint32\s+`,
	TOKEN_KEYWORD_UINT64:   `^\s*uint64\s+`,
	TOKEN_KEYWORD_FIXED32:  `^\s*fixed32\s+`,
	TOKEN_KEYWORD_FIXED64:  `^\s*fixed64\s+`,
	TOKEN_KEYWORD_SFIXED32: `^\s*sfixed32\s+`,
	TOKEN_KEYWORD_SFIXED64: `^\s*sfixed64\s+`,
	TOKEN_KEYWORD_BOOL:     `^\s*bool\s+`,
	TOKEN_KEYWORD_STRING:   `^\s*string\s+`,
	TOKEN_KEYWORD_BYTES:    `^\s*bytes\s+`,
	TOKEN_KEYWORD_OPTION:   `^\s*option\s+`,
	TOKEN_KEYWORD_REQUIRED: `^\s*required\s+`,
	TOKEN_KEYWORD_REPEATED: `^\s*repeated\s+`,
	TOKEN_SYMBOL:           `\s*[\w]+\s*`,
	TOKEN_NUMBER:           `\s*[\d]+\s*`,
}

type Token struct {
	lineno    int
	tokenType int
	text      string
}

type Lexer struct {
	lines        []string
	rules        map[int]*regexp.Regexp
	tokens       []*Token
	currTokenIdx int
}

func newLexer(tplReader io.Reader) *Lexer {
	lexer := &Lexer{}
	lexer.init(tplReader)
	return lexer
}

func (lexer *Lexer) init(tplReader io.Reader) {
	lexer.currTokenIdx = 0
	lexer.rules = map[int]*regexp.Regexp{}
	for k, v := range token_rules {
		reg := regexp.MustCompile(v)
		lexer.rules[k] = reg
	}

	scanner := bufio.NewScanner(tplReader)
	for scanner.Scan() {
		line := scanner.Text()
		lexer.lines = append(lexer.lines, line)
	}
	for k, v := range lexer.lines {
		lexer.parseLine(k+1, v)
	}

	//for _, v := range lexer.tokens {
	//	log.Println("line ", v.lineno, v.text)
	//}
}

func (lexer *Lexer) parseLine(lineno int, lineText string) {
	line := lineText
	for len(line) > 0 {
		isMatch := false
		for tokenType := TOKEN_MIN + 1; tokenType < TOKEN_MAX; tokenType++ {
			reg := lexer.rules[tokenType]
			if reg == nil {
				continue
			}
			ret := reg.FindStringIndex(line)
			if len(ret) == 2 && ret[0] == 0 {
				// do not process comment words
				if tokenType == TOKEN_COMMENT {
					return
				}
				text := strings.TrimSpace(line[ret[0]:ret[1]])
				lexer.tokens = append(lexer.tokens, &Token{lineno: lineno, tokenType: tokenType, text: text})
				line = line[ret[1]:]
				isMatch = true
				break
			}
		}
		if isMatch == false {
			log.Fatalf("There is some error in proto file with line %v: %v", lineno, line)
			break
		}
	}
}

func (lexer *Lexer) takeToken() *Token {
	if lexer.currTokenIdx >= len(lexer.tokens) {
		return nil
	}
	ret := lexer.tokens[lexer.currTokenIdx]
	lexer.currTokenIdx++
	//fmt.Println("takeToken: ", ret.text)
	return ret
}

func (lexer *Lexer) nextTokenType() int {
	if lexer.currTokenIdx >= len(lexer.tokens) {
		return -1
	}
	return lexer.tokens[lexer.currTokenIdx].tokenType
}
