package main

import (
	"database/sql"
	"log"
	"net"

	"inventory/config"
	"inventory/internal/handler"
	"inventory/internal/repository"
	"inventory/internal/usecase"

	inventorypb "inventory/proto/inventorypb"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	// Connect to PostgreSQL
	db, err := sql.Open("postgres", config.GetDBConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// Set up Repositories and Usecases
	productRepo := repository.NewProductRepository(db)
	productUsecase := usecase.NewProductUsecase(productRepo)

	categoryRepo := repository.NewCategoryRepository(db)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)

	// Set up gRPC Handlers

	productGRPCHandler := handler.NewInventoryGRPCHandler(productUsecase, categoryUsecase)

	// Start gRPC server
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	grpcServer := grpc.NewServer()
	inventorypb.RegisterInventoryServiceServer(grpcServer, productGRPCHandler)

	log.Println("🚀 InventoryService gRPC running on :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
