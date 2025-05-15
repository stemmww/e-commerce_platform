package handler

import (
	"context"
	"inventory/internal/model"
	inventorypb "inventory/proto/inventorypb"
)

func (h *InventoryGRPCHandler) ListCategories(ctx context.Context, _ *inventorypb.Empty) (*inventorypb.CategoryList, error) {
	categories, err := h.CategoryUC.GetAll()
	if err != nil {
		return nil, err
	}

	var pbCategories []*inventorypb.Category
	for _, c := range categories {
		pbCategories = append(pbCategories, &inventorypb.Category{
			Id:   int32(c.ID),
			Name: c.Name,
		})
	}

	return &inventorypb.CategoryList{Categories: pbCategories}, nil
}

func (h *InventoryGRPCHandler) CreateCategory(ctx context.Context, req *inventorypb.Category) (*inventorypb.CategoryResponse, error) {
	category := model.Category{
		Name: req.Name,
	}

	err := h.CategoryUC.Create(category)
	if err != nil {
		return nil, err
	}

	return &inventorypb.CategoryResponse{Message: "Category created"}, nil
}

func (h *InventoryGRPCHandler) UpdateCategory(ctx context.Context, req *inventorypb.UpdateCategoryRequest) (*inventorypb.CategoryResponse, error) {
	category := model.Category{
		ID:   int(req.Id),
		Name: req.Name,
	}

	err := h.CategoryUC.Update(category.ID, category)
	if err != nil {
		return nil, err
	}

	return &inventorypb.CategoryResponse{Message: "Category updated"}, nil
}

func (h *InventoryGRPCHandler) DeleteCategory(ctx context.Context, req *inventorypb.CategoryID) (*inventorypb.CategoryResponse, error) {
	err := h.CategoryUC.Delete(int(req.Id))
	if err != nil {
		return nil, err
	}
	return &inventorypb.CategoryResponse{Message: "Category deleted"}, nil
}

func (h *InventoryGRPCHandler) GetCategoryByID(ctx context.Context, req *inventorypb.CategoryID) (*inventorypb.Category, error) {
	category, err := h.CategoryUC.GetByID(int(req.Id))
	if err != nil {
		return nil, err
	}

	return &inventorypb.Category{
		Id:   int32(category.ID),
		Name: category.Name,
	}, nil
}
