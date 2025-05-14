package handler

import (
	"context"
	"order-service/internal/model"
	"order-service/internal/usecase"
	orderpb "order-service/proto/orderpb"
	"strconv"
)

type OrderGRPCHandler struct {
	orderpb.UnimplementedOrderServiceServer
	OrderUC usecase.OrderUsecase
}

func NewOrderGRPCHandler(uc usecase.OrderUsecase) *OrderGRPCHandler {
	return &OrderGRPCHandler{OrderUC: uc}
}

func (h *OrderGRPCHandler) CreateOrder(ctx context.Context, req *orderpb.Order) (*orderpb.OrderResponse, error) {
	// Convert string ID/UserID to int
	id, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, err
	}
	userId, err := strconv.Atoi(req.UserId)
	if err != nil {
		return nil, err
	}

	order := model.Order{
		ID:         id,
		UserID:     userId,
		Total:      float64(req.TotalPrice),
		Status:     req.Status,
		OrderItems: []model.OrderItem{},
	}

	for _, item := range req.Items {
		productId, err := strconv.Atoi(item.ProductId)
		if err != nil {
			return nil, err
		}

		order.OrderItems = append(order.OrderItems, model.OrderItem{
			ProductID: productId,
			Quantity:  int(item.Quantity), // Already int32 to int, fine
		})
	}

	if err := h.OrderUC.CreateOrder(&order); err != nil {
		return nil, err
	}

	return &orderpb.OrderResponse{Message: "Order created"}, nil
}
