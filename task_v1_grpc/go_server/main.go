package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "task_v1_grpc/pb"
)

// greeterServer implements pb.GreeterServer.
type greeterServer struct {
	pb.UnimplementedGreeterServer
}

// SayHello returns a greeting for the given name.
func (s *greeterServer) SayHello(_ context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	if req.GetName() == "" {
		return &pb.HelloReply{Message: "Hello, stranger!"}, nil
	}
	return &pb.HelloReply{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

// NewGreeterServer returns a fresh greeterServer ready to register.
func NewGreeterServer() *greeterServer {
	return &greeterServer{}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, NewGreeterServer())

	log.Printf("gRPC server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
