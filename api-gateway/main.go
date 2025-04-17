package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
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

	router.GET("/products", func(c *gin.Context) {
		resp, err := http.Get("http://inventory_service:8081/products")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reach inventory"})
			return
		}
		defer resp.Body.Close()

		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	router.GET("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		target := "http://inventory_service:8081/products/" + id

		resp, err := http.Get(target)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get product by ID"})
			return
		}
		defer resp.Body.Close()

		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	router.GET("/orders/:id", func(c *gin.Context) {
		id := c.Param("id")
		target := "http://order_service:8082/orders/" + id

		resp, err := http.Get(target)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get order"})
			return
		}
		defer resp.Body.Close()

		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	router.GET("/orders", func(c *gin.Context) {
		resp, err := http.Get("http://order_service:8082/orders")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reach orders"})
			return
		}
		defer resp.Body.Close()

		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	router.POST("/orders", func(c *gin.Context) {
		var order map[string]interface{}
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order data"})
			return
		}

		resp, err := http.Post("http://order_service:8082/orders", "application/json", toJSONReader(order))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
			return
		}
		defer resp.Body.Close()

		c.Status(resp.StatusCode)
		io.Copy(c.Writer, resp.Body)
	})

	router.PATCH("/orders/:id", func(c *gin.Context) {
		target := "http://order_service:8082/orders/" + c.Param("id")

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read body"})
			return
		}

		req, err := http.NewRequest(http.MethodPatch, target, bytes.NewReader(bodyBytes))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create PATCH request"})
			return
		}
		req.Header.Set("Content-Type", c.ContentType())

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to patch order"})
			return
		}
		defer resp.Body.Close()

		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	router.POST("/products", func(c *gin.Context) {
		var product map[string]interface{}
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product data"})
			return
		}

		resp, err := http.Post("http://inventory_service:8081/products", "application/json", toJSONReader(product))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
			return
		}
		defer resp.Body.Close()

		c.Status(resp.StatusCode)
		io.Copy(c.Writer, resp.Body)
	})

	router.PATCH("/products/:id", func(c *gin.Context) {
		target := "http://inventory_service:8081/products/" + c.Param("id")

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read body"})
			return
		}

		req, err := http.NewRequest(http.MethodPatch, target, bytes.NewReader(bodyBytes))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create PATCH request"})
			return
		}
		req.Header.Set("Content-Type", c.ContentType())

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to patch product"})
			return
		}
		defer resp.Body.Close()

		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	router.DELETE("/products/:id", func(c *gin.Context) {
		target := "http://inventory_service:8081/products/" + c.Param("id")
		req, err := http.NewRequest(http.MethodDelete, target, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create DELETE request"})
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	router.Run(":8080") // Gateway listens on 8080
}

func toJSONReader(data interface{}) *bytes.Reader {
	jsonBytes, _ := json.Marshal(data)
	return bytes.NewReader(jsonBytes)
}
