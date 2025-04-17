package inventory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	CategoryID  int     `json:"category_id"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Description string  `json:"description"`
}

type InventoryClient struct {
	BaseURL string
}

func NewInventoryClient(baseURL string) *InventoryClient {
	return &InventoryClient{BaseURL: baseURL}
}

func (c *InventoryClient) GetProduct(id int) (*Product, error) {
	url := fmt.Sprintf("%s/products/%d", c.BaseURL, id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("inventory service returned status %d", resp.StatusCode)
	}

	var product Product
	err = json.NewDecoder(resp.Body).Decode(&product)
	return &product, err
}

func (c *InventoryClient) UpdateStock(id int, newStock int) error {
	// Step 1: Fetch the current product details
	product, err := c.GetProduct(id)
	if err != nil {
		return fmt.Errorf("failed to fetch product: %v", err)
	}

	// Step 2: Update only the stock
	product.Stock = newStock

	// Step 3: Send the full product in the PATCH request
	url := fmt.Sprintf("%s/products/%d", c.BaseURL, id)
	body, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("failed to marshal updated product: %v", err)
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create PATCH request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send PATCH request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update stock, status %d", resp.StatusCode)
	}

	return nil
}
