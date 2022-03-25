package main

import "C"

import (
	"jsonrpcdemo/grpc/grpcserver"
)

func main() {}

//export RunGrpc
func RunGrpc() {
	g := grpcserver.NewGreeter("0.0.0.0:37399")
	g.RunGrpcServer()
}
