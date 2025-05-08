package handler

import (
	pbInventory "consumer-service/pb/inventory"
	"context"
	"log"

	"google.golang.org/grpc"
)

type InventoryHandler struct {
	client pbInventory.InventoryServiceClient
}

func NewInventoryHandler() *InventoryHandler {
	conn, err := grpc.Dial("inventory-service:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to inventory service: %v", err)
	}

	client := pbInventory.NewInventoryServiceClient(conn)

	return &InventoryHandler{
		client: client,
	}
}

func (h *InventoryHandler) DecreaseStock(ctx context.Context, productID int64, quantity int32) error {
	// Get current product
	product, err := h.client.GetProduct(ctx, &pbInventory.ProductID{Id: productID})
	if err != nil {
		return err
	}

	// Update stock
	newStock := product.Stock - quantity
	if newStock < 0 {
		newStock = 0
	}

	_, err = h.client.UpdateProduct(ctx, &pbInventory.Product{
		Id:       product.Id,
		Name:     product.Name,
		Category: product.Category,
		Stock:    newStock,
		Price:    product.Price,
	})

	return err
}
