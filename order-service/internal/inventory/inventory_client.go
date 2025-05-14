package inventory

import (
	"context"
	"log"
	"time"

	inventorypb "order-service/proto/inventorypb"

	"google.golang.org/grpc"
)

type Client interface {
	GetProduct(id string) (*inventorypb.Product, error)
	UpdateStock(id string, newStock int32) error
}

type InventoryClient struct {
	conn   *grpc.ClientConn
	client inventorypb.InventoryServiceClient
}

func NewInventoryClient(address string) Client {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Fatalf("❌ Failed to connect to InventoryService at %s: %v", address, err)
	}

	client := inventorypb.NewInventoryServiceClient(conn)

	log.Println("✅ Connected to InventoryService:", address)
	return &InventoryClient{ // Match this with the struct name
		conn:   conn,
		client: client,
	}
}

// Get product details by ID (used for stock check)
func (c *InventoryClient) GetProduct(id string) (*inventorypb.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := c.client.GetProductByID(ctx, &inventorypb.ProductID{Id: id})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Update product stock (optional - if you want to support stock decrement)
func (c *InventoryClient) UpdateStock(id string, newStock int32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	product, err := c.GetProduct(id)
	if err != nil {
		return err
	}

	product.Stock = newStock

	_, err = c.client.UpdateProduct(ctx, product)
	return err
}
