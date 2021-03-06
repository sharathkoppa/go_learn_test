package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io"
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
	fmt.Println("stream prime server")
	number := req.Number
	k := int32(2)
	if number < 0 {
		err := status.Error(codes.InvalidArgument, "Not valid for negative numbers")
		fmt.Println("---------------------", err.Error())
		return err
	}
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

func (s *server) LongGreet(stream protos.GreetService_LongGreetServer) error {
	fmt.Println("stream greet server")
	resp := &protos.GreetingResponse{}
	str := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			resp.Response = str
			return stream.SendAndClose(resp)
		} else if err != nil {
			fmt.Println("error while receiving stream from client", err.Error())
			return stream.SendAndClose(resp)
		}
		str += "hello " + req.FirstName + " " + req.LastName + "\n"
	}

}

func (s *server) GreetEveryOne(stream protos.GreetService_GreetEveryOneServer) error {
	fmt.Println("stream greet many server")

	str := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			fmt.Println("error while receiving stream from client", err.Error())
			return nil
		}

		str = "hello " + req.FirstName + " " + req.LastName
		resp := &protos.GreetingResponse{Response: str}
		stream.Send(resp)
	}

}

func (s *server) MaxNumber(stream protos.GreetService_MaxNumberServer) error {
	fmt.Println("stream greet max number server")
	maxNumber := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			fmt.Println("error while receiving stream for max number", err.Error())
			return nil
		}
		number := req.GetNumber()
		fmt.Println(number, maxNumber)
		if number > maxNumber {
			maxNumber = number
			stream.Send(&protos.MaxNumberResponse{Response: number})
		}

	}
}

func (s *server) GreetWithDeadLine(ctx context.Context, req *protos.GreetingRequest) (*protos.GreetingResponse, error) {
	fmt.Println("unary greet server")

	// client set the deadline time and server exits after that.
	// here server takes to response more than 3 seconds and if client sets lower time
	// deadline exceeded error will be returned.
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			return nil, status.Error(codes.DeadlineExceeded, "stopped by client")
		}
		time.Sleep(1 * time.Second)
	}

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

	// server key and server certificate for ssl authentication
	// set in server startup.
	certFile := "ssl/server.crt"
	keyFile := "ssl/server.pem"
	transportCreds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	for sslErr != nil {
		fmt.Println("ssl error in server", sslErr)
	}
	creds := grpc.Creds(transportCreds)

	serv := grpc.NewServer(creds)
	protos.RegisterGreetServiceServer(serv, &server{})
	err = serv.Serve(lsi)
	if err != nil {
		fmt.Println("error while creating service", err.Error())
	}

}
