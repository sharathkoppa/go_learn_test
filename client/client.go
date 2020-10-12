package main

import (
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

}
