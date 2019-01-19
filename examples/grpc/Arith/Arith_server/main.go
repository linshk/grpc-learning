package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "github.com/linshk/grpc-learning/examples/grpc/Arith/Arith"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// server is used to implement Arith.ArithServer.
type server struct{}

// Multiply implements Arith.ArithServer
func (s *server) Multiply(ctx context.Context, in *pb.Args) (*pb.Production, error) {
	log.Printf("Received: A = %v, B = %v", in.A, in.B)
	return &pb.Production{Value: in.A * in.B}, nil
}

// Divide implements Arith.ArithServer
func (s *server) Divide(ctx context.Context, in *pb.Args) (*pb.Quotient, error) {
	log.Printf("Received: A = %v, B = %v", in.A, in.B)
	return &pb.Quotient{Quo: in.A / in.B, Rem: in.A % in.B}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Arith server listening at localhost%s", port)
	s := grpc.NewServer()
	pb.RegisterArithServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}