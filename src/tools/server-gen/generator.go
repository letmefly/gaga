package main

import (
	"log"
	"os"
	"strings"
)

type ServerGen struct {
	genType          string
	genFile          string
	serverName       string
	parser           *Parser
	header           string
	protoUseList     string
	grpc             string
	register         string
	sendFunc         string
	pbServerRegister string
}

func newServerGen(genType string, genFile string, name string, parser *Parser) *ServerGen {
	ret := &ServerGen{
		genType:          genType,
		genFile:          genFile,
		serverName:       name,
		parser:           parser,
		header:           "",
		protoUseList:     "",
		grpc:             "",
		register:         "",
		sendFunc:         "",
		pbServerRegister: "",
	}
	return ret
}

func (g *ServerGen) gen() {
	if g.genType == "rpc" {
		g.gen_header()
		g.gen_grpc()
		g.gen_clientApi()
		g.gen_protoUseList()
		g.gen_pbServerRegister()
		g.saveFile(g.genFile, g.header+g.grpc+g.pbServerRegister)
	} else {
		g.gen_header()
		g.gen_grpc()
		g.gen_clientApi()
		g.gen_protoUseList()
		g.gen_pbServerRegister()
		g.saveFile(g.genFile, g.header+g.grpc+g.register+g.sendFunc+g.protoUseList)
	}
}

func (g *ServerGen) gen_header() {
	if g.genType == "rpc" {
		g.header = server_gen_rpc_header_tpl
	} else {
		g.header = server_gen_header_tpl
	}
	g.header = strings.Replace(g.header, "#{server_name}", g.serverName, -1)
}

func (g *ServerGen) gen_grpc() {
	for _, defService := range g.parser.services {
		//log.Printf("service %s\n", defService.name)
		for _, rpc := range defService.rpcList {
			//log.Printf("  %s (%s) (%s) isStream %d %d\n", rpc.name, rpc.param, rpc.ret, rpc.isParamStream, rpc.isRetStream)
			if rpc.isParamStream && rpc.isRetStream {
				//rpcStr := strings.Replace(server_gen_grpc_tpl, "#{rpc_name}", rpc.name, -1)
				//rpcParam := fmt.Sprintf("%s_%sServer", defService.name, defService.name)
				//rpcStr = strings.Replace(rpcStr, "#{rpc_param}", rpcParam, -1)
				//rpcRet := fmt.Sprintf("%s_%sServer", defService.name, defService.name)
				//rpcStr = strings.Replace(rpcStr, "#{rpc_ret}", rpcRet, -1)
				//g.grpc += rpcStr + "\n"
				//log.Printf("  %s (stream %s) (stream %s)\n", rpc.name, rpc.param, rpc.ret)
			} else if rpc.isParamStream {
				//log.Printf("  %s (stream %s) (%s)\n", rpc.name, rpc.param, rpc.ret)
				//rpcStr := strings.Replace(server_gen_grpc_tpl, "#{rpc_name}", rpc.name, -1)
				//rpcStr = strings.Replace(rpcStr, "#{rpc_param}", rpc.param, -1)
				//rpcStr = strings.Replace(rpcStr, "#{rpc_ret}", rpc.ret, -1)
				//g.grpc += rpcStr + "\n"
			} else if rpc.isRetStream {
				//log.Printf("  %s (%s) (stream %s)\n", rpc.name, rpc.param, rpc.ret)
			} else {
				//log.Printf("  %s (%s) (%s)\n", rpc.name, rpc.param, rpc.ret)
				rpcStr := strings.Replace(server_gen_grpc_tpl, "#{rpc_name}", rpc.name, -1)
				rpcStr = strings.Replace(rpcStr, "#{rpc_param}", rpc.param, -1)
				rpcStr = strings.Replace(rpcStr, "#{rpc_ret}", rpc.ret, -1)
				g.grpc += rpcStr + "\n"
			}
		}
	}
}

func (g *ServerGen) gen_clientApi() {
	g.register += server_gen_register_part1_tpl
	for _, defType := range g.parser.types {
		//log.Printf("%s %s %s\n", defType.parentType, defType.def, defType.name)
		if defType.def == "message" && defType.parentType == "" {
			registerStr := server_gen_register_part2_tpl
			registerStr = strings.Replace(registerStr, "#{message_name}", defType.name, -1)
			sendFuncStr := server_gen_sendfunc_tpl
			sendFuncStr = strings.Replace(sendFuncStr, "#{message_name}", defType.name, -1)
			//log.Println(defType.name, sendFuncStr)
			g.register += registerStr
			g.sendFunc += sendFuncStr

			ntfFuncStr := server_gen_ntffunc_tpl
			ntfFuncStr = strings.Replace(ntfFuncStr, "#{message_name}", defType.name, -1)
			g.sendFunc += ntfFuncStr
		}
	}
	g.register += server_gen_register_part3_tpl
}

func (g *ServerGen) gen_protoUseList() {
	messageList := ""
	for _, defType := range g.parser.types {
		if defType.def == "message" && defType.parentType == "" {
			messageList += "\n		\"" + defType.name + "\"" + ","
		}
	}
	protoUseList := server_gen_protouse_tpl
	protoUseList = strings.Replace(protoUseList, "#{proto_use_list}", messageList, -1)
	g.protoUseList = protoUseList

}

func (g *ServerGen) gen_pbServerRegister() {
	pbServerRegister := server_gen_pbregister_tpl
	for _, defService := range g.parser.services {
		pbServerRegister = strings.Replace(pbServerRegister, "#{service_name}", defService.name, -1)
	}
	g.pbServerRegister = pbServerRegister
}

func (g *ServerGen) saveFile(fileName string, txt string) {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("[ServerGen] Create files error: %v", err)
	}
	defer f.Close()
	f.WriteString(txt)
}
