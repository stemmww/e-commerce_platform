package nats

import (
	"consumer-service/internal/handler"
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

type Subscriber struct {
	nc               *nats.Conn
	inventoryHandler *handler.InventoryHandler
}

func NewSubscriber(natsURL string) (*Subscriber, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}

	invHandler := handler.NewInventoryHandler()

	return &Subscriber{
		nc:               nc,
		inventoryHandler: invHandler,
	}, nil
}

func (s *Subscriber) Subscribe(subject string) error {
	_, err := s.nc.Subscribe(subject, func(msg *nats.Msg) {
		log.Printf("📩 Received message on [%s]: %s", subject, string(msg.Data))
		s.HandleOrderCreated(msg.Data)
	})
	return err
}

func (s *Subscriber) HandleOrderCreated(data []byte) {
	var order struct {
		ID     int    `json:"id"`
		UserID int    `json:"user_id"`
		Status string `json:"status"`
		Items  []struct {
			ProductID int `json:"product_id"`
			Quantity  int `json:"quantity"`
		} `json:"items"`
	}

	if err := json.Unmarshal(data, &order); err != nil {
		log.Printf("❌ Failed to unmarshal order: %v", err)
		return
	}

	log.Printf("📦 Processing Order ID: %d", order.ID)

	for _, item := range order.Items {
		log.Printf("🔎 Checking stock for Product ID: %d", item.ProductID)

		err := s.inventoryHandler.DecreaseStock(context.Background(), int64(item.ProductID), int32(item.Quantity))
		if err != nil {
			log.Printf("⚠️ Failed to decrease stock for product %d: %v", item.ProductID, err)
			continue
		}

		log.Printf("✅ Successfully decreased stock for product %d by %d units", item.ProductID, item.Quantity)
	}
}

func (s *Subscriber) Close() {
	s.nc.Close()
}
