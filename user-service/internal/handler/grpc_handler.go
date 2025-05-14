package handler

import (
	"context"
	"log"
	"user-service/internal/model"
	"user-service/internal/usecase"
	userpb "user-service/proto/userpb"
)

type UserGRPCHandler struct {
	userpb.UnimplementedUserServiceServer
	usecase usecase.UserUsecase
}

func NewUserGRPCHandler(uc usecase.UserUsecase) *UserGRPCHandler {
	return &UserGRPCHandler{usecase: uc}
}

func (h *UserGRPCHandler) RegisterUser(ctx context.Context, req *userpb.UserRequest) (*userpb.UserResponse, error) {
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	err := h.usecase.Register(user)
	if err != nil {
		return nil, err
	}

	return &userpb.UserResponse{Message: "User registered successfully"}, nil
}

func (h *UserGRPCHandler) AuthenticateUser(ctx context.Context, req *userpb.AuthRequest) (*userpb.AuthResponse, error) {
	log.Printf("Authenticating user: %s", req.Email)

	user, err := h.usecase.Authenticate(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &userpb.AuthResponse{
		Token:   "dummy-token", // Replace with real token generation later
		Message: "Authentication successful for user: " + user.Username,
	}, nil
}

func (h *UserGRPCHandler) GetUserProfile(ctx context.Context, req *userpb.UserID) (*userpb.UserProfile, error) {
	user, err := h.usecase.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &userpb.UserProfile{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
