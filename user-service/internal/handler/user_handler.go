package handler

import (
	"context"
	"user-service/internal/model"
	"user-service/internal/repository"
	"user-service/internal/usecase"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	pb "user-service/pb/user" // Adjust this path to match your actual proto module
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	usecase usecase.UserUsecase
}

func NewUserHandler(db *sqlx.DB) *UserHandler {
	repo := repository.NewUserRepository(db)
	uc := usecase.NewUserUsecase(repo)
	return &UserHandler{usecase: uc}
}

func (h *UserHandler) RegisterUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	user := &model.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	}

	if err := h.usecase.Register(user); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return &pb.UserResponse{
		Id:    int64(user.ID),
		Email: user.Email,
		Name:  user.Name,
	}, nil
}

func (h *UserHandler) AuthenticateUser(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	user, err := h.usecase.Login(req.Email, req.Password)
	if err != nil {
		return &pb.AuthResponse{
			Success: false,
			Message: "Invalid credentials",
		}, status.Errorf(codes.Unauthenticated, "authentication failed")
	}

	return &pb.AuthResponse{
		Success: true,
		Message: "Login successful",
		User: &pb.UserResponse{
			Id:    int64(user.ID),
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}

func (h *UserHandler) GetUserProfile(ctx context.Context, req *pb.UserID) (*pb.UserProfile, error) {
	user, err := h.usecase.GetProfile(int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &pb.UserProfile{
		Id:    int64(user.ID),
		Email: user.Email,
		Name:  user.Name,
	}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	user := model.User{
		ID:       int(req.Id),
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password, // optional: ignore if not updating password
	}

	if err := h.usecase.UpdateUser(user); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &pb.UserResponse{
		Id:    int64(user.ID),
		Email: user.Email,
		Name:  user.Name,
	}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.UserID) (*emptypb.Empty, error) {
	if err := h.usecase.DeleteUser(int(req.Id)); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}
	return &emptypb.Empty{}, nil
}
