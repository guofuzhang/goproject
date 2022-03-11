package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"net"
)

const (
	Address = ":9988"
)

type animalService struct{}

var AnimalService animalService

func main() {
	listener, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatal(err)
	}
	s := grpc.NewServer()
	proto_demo
}
