package main

import (
	"context"
	"fmt"
	"grpc-study/sum/sumpb"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *sumpb.SumRequest) (*sumpb.SumResponse, error) {
	fmt.Printf("Greet function was invoked with %v", req)
	firstNum := req.GetTotal().GetFirstNum()
	secondNum := req.GetTotal().GetSecondNum()

	total := firstNum + secondNum
	res := &sumpb.SumResponse{
		Result: total,
	}
	return res, nil
}

func (*server) PrimeNumberManyTimes(req *sumpb.PrimeNumberManyTimesRequest, stream sumpb.SumService_PrimeNumberManyTimesServer) error {
	fmt.Printf("PrimeNumberManyTimes function was invoked with %v", req)
	num := req.GetPrimeNumber().GetNum()
	divisor := int64(2)
	for num > 1 {
		if num%divisor == 0 {
			stream.Send(&sumpb.PrimeNumberManyTimesResponse{
				Result: divisor,
			})
			num = num / divisor
		} else {
			divisor++
			fmt.Printf("Divisor has increased to %v", divisor)
		}
	}

	return nil
}

func (*server) ComputeAverage(stream sumpb.SumService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage function was invoked with a streaming request\n")
	result := int32(0)
	count := 0
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			average := float64(result) / float64(count)
			return stream.SendAndClose(&sumpb.ComputeAverageResponse{
				Result: average,
			})
		}
		result += req.GetNum()
		count++
	}

}

func main() {
	fmt.Println("Hello world")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Failed to listen: %v", err)
	}
	s := grpc.NewServer()

	sumpb.RegisterSumServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
