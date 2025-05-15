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
	orderID, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, err
	}

	userID, err := strconv.Atoi(req.UserId)
	if err != nil {
		return nil, err
	}

	order := &model.Order{
		ID:         orderID,
		UserID:     userID,
		TotalPrice: float64(req.TotalPrice),
		Status:     req.Status,
	}

	for _, item := range req.Items {
		productID, err := strconv.Atoi(item.ProductId)
		if err != nil {
			return nil, err
		}

		order.Items = append(order.Items, model.OrderItem{
			ProductID: productID,
			Quantity:  int(item.Quantity),
		})
	}

	if err := h.OrderUC.CreateOrder(order); err != nil {
		return nil, err
	}

	return &orderpb.OrderResponse{Message: "Order created"}, nil
}

func (h *OrderGRPCHandler) GetOrderById(ctx context.Context, req *orderpb.OrderID) (*orderpb.Order, error) {
	id, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, err
	}

	order, err := h.OrderUC.GetOrderByID(id)
	if err != nil {
		return nil, err
	}

	var items []*orderpb.OrderItem
	for _, item := range order.Items {
		items = append(items, &orderpb.OrderItem{
			ProductId: strconv.Itoa(item.ProductID),
			Quantity:  int32(item.Quantity),
		})
	}

	return &orderpb.Order{
		Id:         strconv.Itoa(order.ID),
		UserId:     strconv.Itoa(order.UserID),
		TotalPrice: float32(order.TotalPrice),
		Status:     order.Status,
		Items:      items,
	}, nil
}

func (h *OrderGRPCHandler) ListOrders(ctx context.Context, _ *orderpb.Empty) (*orderpb.OrderList, error) {
	orders, err := h.OrderUC.GetAllOrders()
	if err != nil {
		return nil, err
	}

	var orderList []*orderpb.Order
	for _, o := range orders {
		var items []*orderpb.OrderItem
		for _, item := range o.Items {
			items = append(items, &orderpb.OrderItem{
				ProductId: strconv.Itoa(item.ProductID),
				Quantity:  int32(item.Quantity),
			})
		}

		orderList = append(orderList, &orderpb.Order{
			Id:         strconv.Itoa(o.ID),
			UserId:     strconv.Itoa(o.UserID),
			TotalPrice: float32(o.TotalPrice),
			Status:     o.Status,
			Items:      items,
		})
	}

	return &orderpb.OrderList{Orders: orderList}, nil
}

func (h *OrderGRPCHandler) UpdateStatus(ctx context.Context, req *orderpb.OrderStatusUpdate) (*orderpb.OrderResponse, error) {
	id, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, err
	}

	err = h.OrderUC.UpdateStatus(id, req.Status)
	if err != nil {
		return nil, err
	}

	return &orderpb.OrderResponse{Message: "Order status updated"}, nil
}

func (h *OrderGRPCHandler) DeleteOrder(ctx context.Context, req *orderpb.OrderID) (*orderpb.OrderResponse, error) {
	id, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, err
	}

	err = h.OrderUC.DeleteOrder(id)
	if err != nil {
		return nil, err
	}

	return &orderpb.OrderResponse{Message: "Order deleted"}, nil
}
