package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/sharathkoppa/go_learn_test/protos"
	"google.golang.org/grpc"
)

type server struct {
}

func (s *server) Greet(ctx context.Context, req *protos.GreetingRequest) (*protos.GreetingResponse, error) {
	fmt.Println("unary greet server")
	firstName := req.FirstName
	lastName := req.LastName
	responseStr := "Hello " + firstName + " " + lastName
	resp := &protos.GreetingResponse{Response: responseStr}
	return resp, nil
}

func (s *server) Sum(ctx context.Context, req *protos.SumRequest) (*protos.SumResponse, error) {
	summation := req.FirstNumber + req.SecondNumber
	return &protos.SumResponse{Response: summation}, nil
}

func (s *server) GreetManyTimes(req *protos.GreetingRequest, stream protos.GreetService_GreetManyTimesServer) error {
	fmt.Println("stream greet server")
	firstName := req.FirstName
	lastName := req.LastName
	for i := 0; i < 10; i++ {
		responseStr := "Hello " + firstName + " " + lastName
		resp := &protos.GreetingResponse{Response: responseStr}
		stream.Send(resp)
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (s *server) PrimeCheck(req *protos.PrimeDecompostionRequest, stream protos.GreetService_PrimeCheckServer) error {
	number := req.Number
	k := int32(2)
	for number > 1 {
		fmt.Println(number)
		if (number % k) == 0 {
			resp := &protos.PrimeDecompostionResponse{Response: k}
			stream.Send(resp)
			time.Sleep(1 * time.Second)
			number = number / k
		} else {
			k = k + 1
		}
	}

	return nil
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
