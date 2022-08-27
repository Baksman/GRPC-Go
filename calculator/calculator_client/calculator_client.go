package main

import (
	"context"
	"fmt"
	"github.com/baksman/grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	fmt.Println("hello im a client")

	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error occured %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	doUnary(c)

}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("starting to do unary rpc")
	calculatorResponse, err := c.Sum(context.Background(), &calculatorpb.SumRequest{
		FirstNumber:  10,
		SecondNumber: 20,
	})

	if err != nil {
		log.Fatalf("%v while calling sum", err)
	}

	log.Printf("Response from sum ")
	fmt.Printf("created client %v", calculatorResponse.SumResponse)
}
