package main

import (
	"database/sql"
	"log"
	"net"
	"order-service/internal/handler"
	"order-service/internal/inventory"
	"order-service/internal/repository"
	"order-service/internal/usecase"
	orderpb "order-service/proto/orderpb"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	// Connect to PostgreSQL
	db, err := sql.Open("postgres", getDSN())
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Setup Repository and Inventory gRPC Client
	orderRepo := repository.NewOrderRepository(db)
	invClient := inventory.NewInventoryClient("localhost:50051") // InventoryService gRPC address

	// Usecase with inventory dependency
	orderUsecase := usecase.NewOrderUsecase(orderRepo, invClient)

	// Setup gRPC Handler
	orderGRPCHandler := handler.NewOrderGRPCHandler(orderUsecase)

	// gRPC Server
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("❌ Failed to listen on :50052: %v", err)
	}

	grpcServer := grpc.NewServer()
	orderpb.RegisterOrderServiceServer(grpcServer, orderGRPCHandler)

	log.Println("🚀 OrderService gRPC running on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("❌ Failed to serve: %v", err)
	}
}

// getDSN reads DB connection string from env or uses default fallback
func getDSN() string {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:admin@localhost:5432/orderdb?sslmode=disable"
	}
	return dsn
}
