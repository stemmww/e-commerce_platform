package handler

import (
	"net/http"
	"order-service/internal/inventory"
	"order-service/internal/model"
	"order-service/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderUsecase    usecase.OrderUsecase
	inventoryClient *inventory.InventoryClient
}

func NewOrderHandler(orderUC usecase.OrderUsecase, invClient *inventory.InventoryClient) *OrderHandler {
	return &OrderHandler{
		orderUsecase:    orderUC,
		inventoryClient: invClient,
	}
}

func RegisterOrderRoutes(router *gin.Engine, h *OrderHandler) {
	orderGroup := router.Group("/orders")
	{
		orderGroup.POST("/", h.CreateOrder)
		orderGroup.GET("/", h.GetAllOrders)
		orderGroup.GET("/:id", h.GetOrderByID)
		orderGroup.PATCH("/:id", h.UpdateOrderStatus)

	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req model.Order
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	product, err := h.inventoryClient.GetProduct(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
		return
	}

	if product.Stock < req.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough stock"})
		return
	}

	newStock := product.Stock - req.Quantity
	err = h.inventoryClient.UpdateStock(req.ProductID, newStock)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock"})
		return
	}

	req.Status = "pending"

	err = h.orderUsecase.CreateOrder(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Order created"})
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	orders, err := h.orderUsecase.GetAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.orderUsecase.GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var payload struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil || payload.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status required"})
		return
	}

	if err := h.orderUsecase.UpdateStatus(id, payload.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated"})
}
