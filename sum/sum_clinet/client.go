package main

import (
	"context"
	"fmt"
	"grpc-study/sum/sumpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBiDiStreaming(c)

	doErrorUnary(c)

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

func doClientStreaming(c sumpb.SumServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")
	stream, err := c.ComputeAverage(context.Background())
	requests := []*sumpb.ComputeAverageRequest{
		&sumpb.ComputeAverageRequest{
			Num: 1,
		},
		&sumpb.ComputeAverageRequest{
			Num: 2,
		},
		&sumpb.ComputeAverageRequest{
			Num: 3,
		},
		&sumpb.ComputeAverageRequest{
			Num: 4,
		},
	}
	if err != nil {
		log.Fatalf("error while calling ComputeAverage RPC: %v", err)
	}
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		fmt.Printf("error while receiving response from ComputeAverage: %v", err)
	}
	fmt.Printf("ComputeAverage: %v\n", res)

}

func doBiDiStreaming(c sumpb.SumServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("error while calling FindMaximum RPC: %v", err)
	}
	requests := []int32{1, 5, 3, 6, 2, 20}

	waitc := make(chan struct{})
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending number: %v\n", req)
			stream.Send(&sumpb.FindMaximumRequest{
				Num: req,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v\n", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetMaximum())
		}
		close(waitc)
	}()
	<-waitc
}

func doErrorUnary(c sumpb.SumServiceClient) {
	println("Starting to do a Error Unary RPC...")

	doErrorCall(c, 10)
	doErrorCall(c, -1)

}

func doErrorCall(c sumpb.SumServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &sumpb.SquareRootRequest{
		Number: int32(n),
	})
	if err != nil {
		resErr, ok := status.FromError(err)
		if ok {
			fmt.Printf("Error message from server: %v\n", resErr.Message())
			fmt.Println(resErr.Code())
			if resErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
			}
		} else {
			log.Fatal("Big Error calling SquareRoot: %v", err)
		}
	}
	fmt.Printf("Result of square root of %v: %v\n", n, res.GetNumberRoot())
}
