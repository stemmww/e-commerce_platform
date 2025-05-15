package main

import (
	"errors"
	"log"
	"net"

	"user-service/internal/handler"
	"user-service/internal/model"
	"user-service/internal/repository"
	"user-service/internal/usecase"
	"user-service/proto/userpb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

// InMemoryRepo is a temporary in-memory implementation of UserRepository.
type InMemoryRepo struct {
	data map[string]*model.User
}

func NewInMemoryRepo() repository.UserRepository {
	return &InMemoryRepo{data: make(map[string]*model.User)}
}

func (r *InMemoryRepo) Register(user *model.User) error {
	user.ID = uuid.New().String()
	r.data[user.ID] = user
	return nil
}

func (r *InMemoryRepo) Authenticate(email, password string) (*model.User, error) {
	for _, u := range r.data {
		if u.Email == email && u.Password == password {
			return u, nil
		}
	}
	return nil, errors.New("invalid email or password")
}

func (r *InMemoryRepo) GetByID(id string) (*model.User, error) {
	user, exists := r.data[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("❌ Failed to listen: %v", err)
	}

	repo := NewInMemoryRepo()
	usecase := usecase.NewUserUsecase(repo)
	grpcHandler := handler.NewUserGRPCHandler(usecase)

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, grpcHandler)

	log.Println("✅ User Service is running on port :8083")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("❌ Failed to serve: %v", err)
	}
}
