package main

import (
	"context"
	"fmt"
	"grpc-study/sum/sumpb"
	"io"
	"log"
	"math"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (*server) FindMaximum(stream sumpb.SumService_FindMaximumServer) error {
	fmt.Printf("FindMaximum function was invoked with a streaming request\n")
	maximum := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			fmt.Printf("Error while reading client streaming %v", err)
			return err
		}
		result := req.GetNum()
		if result > maximum {
			sendErr := stream.Send(&sumpb.FindMaximumResponse{
				Maximum: result,
			})
			if sendErr != nil {
				fmt.Printf("Error while sending data to client %v", err)
				return err
			}
			maximum = result
		}

	}
}

func (*server) SquareRoot(ctx context.Context, req *sumpb.SquareRootRequest) (*sumpb.SquareRootResponse, error) {
	fmt.Println("Received SquareRoot RPC")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v", number),
		)
	}
	return &sumpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
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
