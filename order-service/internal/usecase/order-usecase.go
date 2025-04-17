package usecase

import (
	"order-service/internal/model"
	"order-service/internal/repository"
)

type OrderUsecase interface {
	CreateOrder(order *model.Order) error
	GetAllOrders() ([]model.Order, error)
	GetOrderByID(id int) (*model.Order, error)
	UpdateStatus(id int, status string) error
}

type orderUsecase struct {
	orderRepo repository.OrderRepository
}

func NewOrderUsecase(repo repository.OrderRepository) OrderUsecase {
	return &orderUsecase{
		orderRepo: repo,
	}
}

func (u *orderUsecase) CreateOrder(order *model.Order) error {
	return u.orderRepo.Create(order)
}

func (u *orderUsecase) GetAllOrders() ([]model.Order, error) {
	return u.orderRepo.GetAll()
}

func (u *orderUsecase) GetOrderByID(id int) (*model.Order, error) {
	return u.orderRepo.GetByID(id)
}

func (u *orderUsecase) UpdateStatus(id int, status string) error {
	return u.orderRepo.UpdateStatus(id, status)
}
