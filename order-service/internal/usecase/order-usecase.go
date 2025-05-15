package usecase

import (
	"order-service/internal/inventory"
	"order-service/internal/model"
	"order-service/internal/repository"
)

type OrderUsecase interface {
	CreateOrder(order *model.Order) error
	GetAllOrders() ([]model.Order, error)
	GetOrderByID(id int) (*model.Order, error)
	UpdateStatus(id int, status string) error
	DeleteOrder(id int) error
}

type orderUsecase struct {
	orderRepo repository.OrderRepository
	invClient inventory.Client
}

func NewOrderUsecase(repo repository.OrderRepository, inv inventory.Client) OrderUsecase {
	return &orderUsecase{
		orderRepo: repo,
		invClient: inv,
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

func (u *orderUsecase) DeleteOrder(id int) error {
	return u.orderRepo.Delete(id)
}
