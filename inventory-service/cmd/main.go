package main

import (
	"database/sql"
	"fmt"
	"inventory/config"
	"inventory/internal/handler"
	"inventory/internal/repository"
	"inventory/internal/usecase"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", config.GetDBConnectionString())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	categoryRepo := repository.NewCategoryRepository(db)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryUsecase)

	repo := repository.NewProductRepository(db)
	usecase := usecase.NewProductUsecase(repo)
	handler := handler.NewProductHandler(usecase)

	router := gin.Default()
	categoryHandler.RegisterRoutes(router)
	handler.RegisterRoutes(router)

	fmt.Println("🚀 Inventory service running on :8081")
	router.Run(":8081")
}
