package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/baksman/grpc/average/averagepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("hello  from average client")

	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("eror occured while dialing: %v", err)
	}
	defer cc.Close()
	c := averagepb.NewCalculatorServiceClient(cc)
	doClientStreaming(c)

}

func doClientStreaming(c averagepb.CalculatorServiceClient) {
	fmt.Printf("started to do steaming rpc")

	stream, err := c.Average(context.Background())

	if err != nil {
		log.Fatalf("%v while calling long average times", err)
	}

	requests := []*averagepb.AverageRequest{
		{
			AverageRequest: 10,
		},
		{
			AverageRequest: 20,
		},
		{
			AverageRequest: 90,
		},
		{
			AverageRequest: 100,
		},
	}

	for _, v := range requests {
		stream.Send(v)
		time.Sleep(time.Second)
	}

	response, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("error while receiving long average %v", err)
	}

	log.Printf("Response rom %v", response)

}
