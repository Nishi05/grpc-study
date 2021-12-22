package main

import (
	"context"
	"fmt"
	"grpc-study/sum/sumpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := sumpb.NewSumServiceClient(cc)
	// doUnary(c)
	doServerStreaming(c)
}

func doUnary(c sumpb.SumServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &sumpb.SumRequest{
		Total: &sumpb.Sum{
			FirstNum:  3,
			SecondNum: 10,
		},
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Great RPC: %v", err)
	}
	fmt.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c sumpb.SumServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")
	req := &sumpb.PrimeNumberManyTimesRequest{
		PrimeNumber: &sumpb.PrimeNumber{
			Num: 120,
		},
	}
	resStream, err := c.PrimeNumberManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumberManyTimes RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from PrimeNumberManyTimes: %v", msg.GetResult())
	}
}
