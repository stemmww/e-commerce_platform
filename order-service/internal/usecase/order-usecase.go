package usecase

import (
	"errors"
	"order-service/internal/inventory"
	"order-service/internal/model"
	"order-service/internal/repository"
	"strconv"
)

type OrderUsecase interface {
	CreateOrder(order *model.Order) error
	GetAllOrders() ([]model.Order, error)
	GetOrderByID(id int) (*model.Order, error)
	UpdateStatus(id int, status string) error
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
	var total float64

	for _, item := range order.OrderItems {
		product, err := u.invClient.GetProduct(strconv.Itoa(item.ProductID))
		if err != nil {
			return errors.New("failed to fetch product info from inventory")
		}

		if product.Stock < int32(item.Quantity) {
			return errors.New("not enough stock for product: " + product.Name)
		}

		total += float64(item.Quantity) * float64(product.Price)
	}

	order.Total = total
	order.Status = "PENDING"

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
