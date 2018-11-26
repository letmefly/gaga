// Ref: https://github.com/letmefly/ant/blob/master/cmd/ant/create/generator.go
package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	srcRoot := "./"
	dirs, _ := ioutil.ReadDir(srcRoot + "pb")
	for _, dir := range dirs {
		if dir.IsDir() {
			files, _ := ioutil.ReadDir(srcRoot + "pb/" + dir.Name())
			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".proto") {
					//log.Println(file.Name())
					strs := strings.Split(file.Name(), "_")
					serverName := strs[0]
					filePath := srcRoot + "pb/" + dir.Name() + "/" + file.Name()
					log.Println(filePath)
					protoFile, err := os.Open(filePath)
					if err != nil {
						log.Fatal(err)
					}
					lexer := newLexer(protoFile)
					defer protoFile.Close()

					parser := newParser(lexer)
					parser.parse()
					genType := "rpc"
					genFile := ""
					if strings.Contains(file.Name(), "_msg") {
						genType = "msg"
						genFile = srcRoot + "pb/" + dir.Name() + "/" + serverName + "_msg.gen.go"
					} else if strings.Contains(file.Name(), "_rpc") {
						genType = "rpc"
						genFile = srcRoot + "pb/" + dir.Name() + "/" + serverName + "_rpc.gen.go"
					}

					log.Println("genFile:", genFile)
					generator := newServerGen(genType, genFile, serverName, parser)
					generator.gen()
				}
			}
		}
	}
}
