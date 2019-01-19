package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "github.com/linshk/grpc-learning/examples/grpc/Arith/Arith"
)

const (
	address     = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewArithClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	args := pb.Args{A: 9, B: 4}
	r1, err := c.Multiply(ctx, &args)
	if err != nil {
		log.Fatalf("could not call Multiply: %v", err)
	}
	log.Printf("%d * %d = %d", args.A, args.B, r1.Value)

	r2, err := c.Divide(ctx, &args)
	if err != nil {
		log.Fatalf("could not call Divide: %v", err)
	}
	log.Printf("%d / %d = %d remains %d", args.A, args.B, r2.Quo, r2.Rem)
}