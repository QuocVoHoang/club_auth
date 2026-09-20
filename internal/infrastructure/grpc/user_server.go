package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/your-org/go-base/internal/domain/repository"
	grpcUser "github.com/your-org/go-base/internal/infrastructure/grpc/generated/user"
	"github.com/your-org/go-base/pkg/http_error"
	"github.com/your-org/go-base/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserServer implements the generated UserServiceServer interface.
type UserServer struct {
	grpcUser.UnimplementedUserServiceServer
	userRepo repository.IUserRepo
}

func NewUserServer(userRepo repository.IUserRepo) *UserServer {
	return &UserServer{userRepo: userRepo}
}

func (s *UserServer) GetUser(ctx context.Context, req *grpcUser.GetUserRequest) (*grpcUser.GetUserResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	user, err := s.userRepo.Get(ctx, userID)
	if err != nil {
		log.Errorf("get user by id %s: %v", userID, err)
		return nil, status.Error(codes.Internal, "internal error")
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, http_error.NotFoundError("user not found").Error())
	}

	return &grpcUser.GetUserResponse{
		Id:       user.ID.String(),
		Email:    user.Email,
		FullName: user.FullName,
		Role:     int32(user.Role),
		Status:   user.Status,
	}, nil
}
