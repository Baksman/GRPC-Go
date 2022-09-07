package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/baksman/grpc/average/averagepb"
	"google.golang.org/grpc"
)

type server struct{}

// Average(CalculatorService_AverageServer) error
func (*server) Average(stream averagepb.CalculatorService_AverageServer) error {
	result := 0
	counter := 0

	for {
		
		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&averagepb.AverageResponse{
				Result: int32(result / counter),
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream %v", err)
		}

		number := req.GetAverageRequest()
		result += int(number)
		counter++

	}
}

func main() {
	fmt.Println("started average server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	s := grpc.NewServer()

	averagepb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}
