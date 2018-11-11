#!/bin/bash
p=`pwd`
export PATH=$PATH:$p/tools/protoc/mac/bin
#echo $PATH
rm -fr ./protos/*.go
protoc3 --go_out=plugins=grpc:. ./pb/*.proto
protoc3 --go_out=plugins=grpc:. ./pb/template/*.proto
#./tools/server-gen/server-gen
