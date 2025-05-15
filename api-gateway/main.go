package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	inventorypb "api-gateway/proto/inventorypb"
	orderpb "api-gateway/proto/orderpb"
	userpb "api-gateway/proto/userpb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func reverseProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, err := url.Parse(target)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to parse proxy URL"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, c.FullPath())
		c.Request.Host = remote.Host
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	router := gin.Default()

	// Connect to User Service
	userConn, err := grpc.Dial("user_service:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to User Service: %v", err)
	}
	defer userConn.Close()
	userClient := userpb.NewUserServiceClient(userConn)

	// Connect to Inventory Service
	invConn, err := grpc.Dial("inventory_service:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Inventory Service: %v", err)
	}
	defer invConn.Close()
	invClient := inventorypb.NewInventoryServiceClient(invConn)

	orderConn, err := grpc.Dial("order_service:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Order Service: %v", err)
	}
	defer orderConn.Close()
	orderClient := orderpb.NewOrderServiceClient(orderConn)

	// USER ROUTES ------------------------------------------------------------------
	router.POST("/users/register", func(c *gin.Context) {
		var req userpb.UserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := userClient.RegisterUser(ctx, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, res)
	})

	router.POST("/users/login", func(c *gin.Context) {
		var req userpb.AuthRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := userClient.AuthenticateUser(ctx, &req)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	router.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := userClient.GetUserProfile(ctx, &userpb.UserID{Id: id})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	// Optional health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	//SER ROUTES ENDING ------------------------------------------------------------------

	// REPLACING HTTP ROUTES WITH GRPC ----------------------------------------------------------------
	router.GET("/grpc-products", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		resp, err := invClient.ListProducts(ctx, &inventorypb.Empty{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	})

	router.GET("/products", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := invClient.ListProducts(ctx, &inventorypb.Empty{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res.Products)
	})

	router.GET("/products/:id", func(c *gin.Context) {
		id := c.Param("id")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := invClient.GetProductByID(ctx, &inventorypb.ProductID{Id: id})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	router.GET("/orders/:id", func(c *gin.Context) {
		id := c.Param("id")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := orderClient.GetOrderById(ctx, &orderpb.OrderID{Id: id})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	router.GET("/orders", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := orderClient.ListOrders(ctx, &orderpb.Empty{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res.Orders)
	})

	router.POST("/orders", func(c *gin.Context) {
		var order orderpb.Order
		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := orderClient.CreateOrder(ctx, &order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, res)
	})

	router.PATCH("/orders/:id", func(c *gin.Context) {
		id := c.Param("id")

		var payload struct {
			Status string `json:"status"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil || payload.Status == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := orderClient.UpdateStatus(ctx, &orderpb.OrderStatusUpdate{
			Id:     id,
			Status: payload.Status,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order status updated"})
	})

	router.POST("/products", func(c *gin.Context) {
		var req inventorypb.Product
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := invClient.CreateProduct(ctx, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, res)
	})

	router.PATCH("/products/:id", func(c *gin.Context) {
		id := c.Param("id")

		var updated inventorypb.Product
		if err := c.ShouldBindJSON(&updated); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Ensure ID in URL is used
		updated.Id = id

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := invClient.UpdateProduct(ctx, &updated)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	router.DELETE("/products/:id", func(c *gin.Context) {
		id := c.Param("id")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := invClient.DeleteProduct(ctx, &inventorypb.ProductID{Id: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	router.GET("/categories", func(c *gin.Context) {
		res, err := invClient.ListCategories(context.Background(), &inventorypb.Empty{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res.Categories)
	})

	router.POST("/categories", func(c *gin.Context) {
		var req inventorypb.Category
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := invClient.CreateCategory(ctx, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, res)
	})

	router.PATCH("/categories/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}

		var body struct {
			Name string `json:"name"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := invClient.UpdateCategory(ctx, &inventorypb.UpdateCategoryRequest{
			Id:   int32(id),
			Name: body.Name,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	router.DELETE("/categories/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := invClient.DeleteCategory(ctx, &inventorypb.CategoryID{Id: int32(id)})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	router.GET("/categories/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := invClient.GetCategoryByID(ctx, &inventorypb.CategoryID{Id: int32(id)})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	router.Run(":8080") // Gateway listens on 8080
}

func toJSONReader(data interface{}) *bytes.Reader {
	jsonBytes, _ := json.Marshal(data)
	return bytes.NewReader(jsonBytes)
}

// Products to gRPC: GET(All and by ID, also i have grpc-products GET), DELETE, PATCH, POST
// Orders to gRPC: GET(All and by ID),
