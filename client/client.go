package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/sharathkoppa/go_learn_test/protos"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("i am in client")
	cc, err := grpc.Dial("0.0.0.0:1234", grpc.WithInsecure())
	if err != nil {
		fmt.Println("error while connecting to port", err)
	}
	defer cc.Close()
	c := protos.NewGreetServiceClient(cc)
	fmt.Println("connection successful", c)
	GreetUnary(c)
	ArithmeticUnary(c)
	GreetStreamServer(c)
	PrimeDecomposition(c)
	LongGreet(c)

}

func GreetUnary(c protos.GreetServiceClient) {
	fmt.Println("In unary")
	req := protos.GreetingRequest{FirstName: "sharath", LastName: "koppa"}
	ctx := context.Background()
	resp, err := c.Greet(ctx, &req)

	if err != nil {
		fmt.Println("grpc response error ", err)
	}

	respJson, _ := json.Marshal(resp)
	fmt.Println(string(respJson))
}

func ArithmeticUnary(c protos.GreetServiceClient) {
	fmt.Println("In unary")
	req := protos.SumRequest{FirstNumber: 10, SecondNumber: 20}
	ctx := context.Background()
	resp, err := c.Sum(ctx, &req)

	if err != nil {
		fmt.Println("grpc response error ", err)
	}

	respJson, _ := json.Marshal(resp)
	fmt.Println(string(respJson))
}

func GreetStreamServer(c protos.GreetServiceClient) {
	fmt.Println("In stream server")
	req := protos.GreetingRequest{FirstName: "sharath", LastName: "koppa"}
	ctx := context.Background()
	respStream, err := c.GreetManyTimes(ctx, &req)

	if err != nil {
		fmt.Println("grpc response error ", err)
	}

	for {
		resp, err := respStream.Recv()
		if err == io.EOF {
			fmt.Println("end of stream ", err)
			break
		} else if err != nil {
			fmt.Println("error in stream response", err)
		}
		respJson, _ := json.Marshal(resp)
		fmt.Println(string(respJson))
	}
}

func PrimeDecomposition(c protos.GreetServiceClient) {
	fmt.Println("In stream server prime")
	req := protos.PrimeDecompostionRequest{Number: 120}
	ctx := context.Background()
	respStream, err := c.PrimeCheck(ctx, &req)
	if err != nil {
		fmt.Println("grpc response error for prime ", err)
	}
	for {
		resp, err := respStream.Recv()
		if err == io.EOF {
			fmt.Println("end of stream prime", err)
			break
		} else if err != nil {
			fmt.Println("error in stream response prime", err)
		}
		respJson, _ := json.Marshal(resp)
		fmt.Println(string(respJson))
	}

}

func LongGreet(c protos.GreetServiceClient) {
	ctx := context.Background()
	var request []*protos.GreetingRequest
	lst := [][]string{{"Sharath", "Koppa"}, {"Sowmya", "Bhat"}, {"Chidu", "Koppa"}, {"Asha", "Devi"}}
	for _, ls := range lst {
		request = append(request, &protos.GreetingRequest{
			FirstName: ls[0],
			LastName:  ls[1],
		})
	}

	stream, err := c.LongGreet(ctx)
	if err != nil {
		fmt.Println("grpc rpc error for long greet ", err)
	}

	for _, req := range request {
		fmt.Println("sending client stream", req)
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}

	resp , err := stream.CloseAndRecv()

	if err != nil {
		fmt.Println("grpc response error for long greet ", err)
	}
	respJson, _ := json.Marshal(resp)
	fmt.Println("resp:", string(respJson))
}