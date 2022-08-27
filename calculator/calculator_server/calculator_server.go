package main

import (
	// "context"
	"context"
	"fmt"
	"log"
	"net"

	"github.com/baksman/grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct {
	// calculatorpb.UnimplementedCalculatorServiceServer
}

func (s *server) Sum(ctx context.Context, req *(calculatorpb.SumRequest)) (*(calculatorpb.SumResponse), error) {
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber

	fmt.Printf("receive form of %v: %v  %v\n ", firstNumber, secondNumber, firstNumber+secondNumber)
	result := firstNumber + secondNumber

	return &calculatorpb.SumResponse{
		SumResponse: result,
	}, nil

}

// Sum(context.Context, *SumRequest) (*SumResponse, error)

func main() {
	fmt.Println("hello word")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("failed to lsiten %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve %v", err)
	}

}
