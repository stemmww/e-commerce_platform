package main

import (
	"log"
	"net"

	"user-service/infrastructure"
	"user-service/internal/handler"

	pb "user-service/pb/user"

	"google.golang.org/grpc"
)

func main() {
	// Connect to DB
	db := infrastructure.NewPostgres()

	// Create a gRPC server
	grpcServer := grpc.NewServer()

	// Create handler that implements pb.UserServiceServer
	userHandler := handler.NewUserHandler(db)

	// Register the handler with gRPC
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	// Start listening
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("User gRPC server running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
