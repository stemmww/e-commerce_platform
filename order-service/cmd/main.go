package main

import (
	"database/sql"
	"fmt"
	"log"
	"order-service/internal/handler"
	"order-service/internal/inventory"
	"order-service/internal/repository"
	"order-service/internal/usecase"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Database connection
	db, err := sql.Open("postgres", getDSN())
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize layers
	orderRepo := repository.NewOrderRepository(db)
	orderUsecase := usecase.NewOrderUsecase(orderRepo)

	inventoryClient := inventory.NewInventoryClient("http://inventory_service:8081") // replace with service name if needed
	orderHandler := handler.NewOrderHandler(orderUsecase, inventoryClient)

	// Router
	router := gin.Default()
	handler.RegisterOrderRoutes(router, orderHandler)

	fmt.Println("🚀 Order service running on :8082")
	router.Run(":8082")
}

func getDSN() string {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		// fallback/default
		dsn = "postgres://postgres:admin@localhost:5432/orderdb?sslmode=disable"
	}
	return dsn
}
