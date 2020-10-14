package main

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io"
	"time"

	"github.com/sharathkoppa/go_learn_test/protos"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("i am in client")

	// certficate to send in client call and server verifies the certificate.
	certFile := "../ssl/ca.crt"
	creds, certErr := credentials.NewClientTLSFromFile(certFile, "")
	fmt.Println(creds.Info())
	if certErr != nil {
		fmt.Println("certificate error in client", certErr)
	}
	credsToPass := grpc.WithTransportCredentials(creds)

	cc, err := grpc.Dial("0.0.0.0:1234", credsToPass)
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
	LongGreetClientStream(c)
	GreetBiDiStream(c)
	GetMaxNumber(c)
	GreetDeadLine(c, 5 *time.Second)
	GreetDeadLine(c, 1 *time.Second)

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
	req := protos.PrimeDecompostionRequest{Number: -120}
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
			respErr, ok := status.FromError(err)
			if ok {
				if respErr.Code() == codes.InvalidArgument {
					fmt.Println("invalid argument", respErr.Message())
				} else {
					fmt.Println("error in stream response prime", respErr.Message())
				}
			} else {
				fmt.Println("error in stream response prime not ok", respErr.Message())
			}
			break
		}
		respJson, _ := json.Marshal(resp)
		fmt.Println(string(respJson))
	}

}

func LongGreetClientStream(c protos.GreetServiceClient) {
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


func GreetBiDiStream(c protos.GreetServiceClient) {
	var request []*protos.GreetingRequest
	lst := [][]string{{"Sharath", "Koppa"}, {"Sowmya", "Bhat"}, {"Chidu", "Koppa"}, {"Asha", "Devi"}}
	for _, ls := range lst {
		request = append(request, &protos.GreetingRequest{
			FirstName: ls[0],
			LastName:  ls[1],
		})
	}

	ctx := context.Background()
	stream, err := c.GreetEveryOne(ctx)
	if err != nil {
		fmt.Println("grpc rpc error for bidi greet ", err)
	}

	for _, req := range request {
		fmt.Println("sending client stream", req)
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}
	stream.CloseSend()

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("error in stream bidi", err)
		}
		respJson, _ := json.Marshal(resp)
		fmt.Println("resp:", string(respJson))
	}
}

func GetMaxNumber(c protos.GreetServiceClient) {
	numbers := []int32{2, 4, 7, 3, 5, 19, 35}

	ctx := context.Background()
	stream, err := c.MaxNumber(ctx)
	if err != nil {
		fmt.Println("grpc rpc error for max number ", err)
	}

	for _, n := range numbers {
		stream.Send(&protos.MaxNumberRequest{Number:n})
		fmt.Println("number sent :", n)
		time.Sleep(1 * time.Second)
	}
	stream.CloseSend()

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("error in stream max number", err)
		}
		fmt.Println("max number is :", resp.Response)
	}

}

func GreetDeadLine(c protos.GreetServiceClient, n time.Duration) {
	fmt.Println("In unary deadline")
	req := protos.GreetingRequest{FirstName: "sharath", LastName: "koppa"}

	// Client sets the duration within server has to send the response.
	// else client terminate the call.
	ctx, cancel := context.WithTimeout(context.Background(), n)
	defer cancel()

	resp, err := c.GreetWithDeadLine(ctx, &req)

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			if respErr.Code() == codes.DeadlineExceeded {
				fmt.Println("dead line occured ", err)
			} else {
				fmt.Println("some other error", err)
			}
		}
		fmt.Println("grpc response error ", err)
		return
	}

	respJson, _ := json.Marshal(resp)
	fmt.Println(string(respJson))
}

