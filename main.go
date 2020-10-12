package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/sharathkoppa/go_learn_test/protos"
	"google.golang.org/grpc"
)

type server struct{
}

func (s *server) Greet(ctx context.Context, req *protos.GreetingRequest) (*protos.GreetingResponse, error) {
	firstName := req.FirstName
	lastName := req.LastName
	responseStr := "Hello " + firstName + " " + lastName
	resp := &protos.GreetingResponse{Response: responseStr}
	return resp, nil
}

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
