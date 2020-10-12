package main

import (
	"context"
	"encoding/json"
	"fmt"

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