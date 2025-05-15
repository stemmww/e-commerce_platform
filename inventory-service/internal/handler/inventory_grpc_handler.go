package handler

import (
	"context"
	"inventory/internal/model"
	"inventory/internal/usecase"
	inventorypb "inventory/proto/inventorypb"
)

type InventoryGRPCHandler struct {
	inventorypb.UnimplementedInventoryServiceServer
	ProductUC  usecase.ProductUsecase
	CategoryUC usecase.CategoryUsecase
}

func NewInventoryGRPCHandler(puc usecase.ProductUsecase, cuc usecase.CategoryUsecase) *InventoryGRPCHandler {
	return &InventoryGRPCHandler{
		ProductUC:  puc,
		CategoryUC: cuc,
	}
}

func (h *InventoryGRPCHandler) CreateProduct(ctx context.Context, req *inventorypb.Product) (*inventorypb.ProductResponse, error) {
	product := &model.Product{
		ID:       req.Id,
		Name:     req.Name,
		Category: int(req.Category),
		Price:    float64(req.Price), // convert if needed
		Stock:    req.Stock,
	}

	err := h.ProductUC.Create(product)

	if err != nil {
		return nil, err
	}

	return &inventorypb.ProductResponse{Message: "Product created"}, nil
}

func (h *InventoryGRPCHandler) GetProductByID(ctx context.Context, req *inventorypb.ProductID) (*inventorypb.Product, error) {
	product, err := h.ProductUC.GetByID(req.Id)
	if err != nil {
		return nil, err
	}
	return &inventorypb.Product{
		Id:       product.ID,
		Name:     product.Name,
		Category: int32(product.Category),
		Price:    float32(product.Price),
		Stock:    product.Stock,
	}, nil
}

func (h *InventoryGRPCHandler) ListProducts(ctx context.Context, _ *inventorypb.Empty) (*inventorypb.ProductList, error) {
	products, err := h.ProductUC.List()
	if err != nil {
		return nil, err
	}

	var pbProducts []*inventorypb.Product
	for _, p := range products {
		pbProducts = append(pbProducts, &inventorypb.Product{
			Id:       p.ID,
			Name:     p.Name,
			Category: int32(p.Category),
			Price:    float32(p.Price),
			Stock:    p.Stock,
		})
	}

	return &inventorypb.ProductList{Products: pbProducts}, nil
}

func (h *InventoryGRPCHandler) UpdateProduct(ctx context.Context, req *inventorypb.Product) (*inventorypb.ProductResponse, error) {
	product := &model.Product{
		ID:       req.Id,
		Name:     req.Name,
		Category: int(req.Category),
		Price:    float64(req.Price),
		Stock:    req.Stock,
	}
	err := h.ProductUC.Update(product)

	if err != nil {
		return nil, err
	}
	return &inventorypb.ProductResponse{Message: "Product updated"}, nil
}

func (h *InventoryGRPCHandler) DeleteProduct(ctx context.Context, req *inventorypb.ProductID) (*inventorypb.ProductResponse, error) {
	err := h.ProductUC.Delete(req.Id)
	if err != nil {
		return nil, err
	}
	return &inventorypb.ProductResponse{Message: "Product deleted"}, nil
}
