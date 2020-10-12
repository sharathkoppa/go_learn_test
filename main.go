package main

import (
	"fmt"
	"net"
	"os"

	"github.com/sharathkoppa/go_learn_test/protos"
	"google.golang.org/grpc"
)

type server struct{}

func main() {
	lsi, err := net.Listen("tcp", "0.0.0.0:1234")
	if err != nil {
		fmt.Println("error while creating port", err.Error())
		os.Exit(1)
	}

	serv := grpc.NewServer()
	protos.RegisterGreetServiceServer(serv, &server{})
	err = serv.Serve(lsi)
	if err != nil {
		fmt.Println("error while creating service", err.Error())
	}

}
