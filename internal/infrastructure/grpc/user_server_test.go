package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/your-org/go-base/internal/domain/entity"
	grpcUser "github.com/your-org/go-base/internal/infrastructure/grpc/generated/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// stubUserRepo implements repository.IUserRepo for testing.
type stubUserRepo struct {
	getResult *entity.User
	getErr    error
}

func (s *stubUserRepo) Create(ctx context.Context, e *entity.User) error { return nil }
func (s *stubUserRepo) Get(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return s.getResult, s.getErr
}
func (s *stubUserRepo) List(ctx context.Context, cursor *uuid.UUID, limit int) ([]entity.User, error) {
	return nil, nil
}
func (s *stubUserRepo) Count(ctx context.Context) (int64, error)         { return 0, nil }
func (s *stubUserRepo) Update(ctx context.Context, e *entity.User) error { return nil }
func (s *stubUserRepo) Delete(ctx context.Context, id uuid.UUID) error   { return nil }
func (s *stubUserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return nil, nil
}
func (s *stubUserRepo) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error { return nil }

func TestUserServer_GetUser(t *testing.T) {
	validID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	tests := []struct {
		name     string
		userID   string
		stubUser *entity.User
		stubErr  error
		wantCode codes.Code
		wantID   string
	}{
		{
			name:     "valid user",
			userID:   validID.String(),
			stubUser: &entity.User{ID: validID, Email: "test@example.com", FullName: "Test User", Role: 3, Status: "active"},
			wantCode: codes.OK,
			wantID:   validID.String(),
		},
		{
			name:     "not found",
			userID:   validID.String(),
			stubUser: nil,
			wantCode: codes.NotFound,
		},
		{
			name:     "invalid uuid",
			userID:   "not-a-uuid",
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "db error",
			userID:   validID.String(),
			stubErr:  errors.New("db down"),
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewUserServer(&stubUserRepo{
				getResult: tt.stubUser,
				getErr:    tt.stubErr,
			})

			resp, err := srv.GetUser(context.Background(), &grpcUser.GetUserRequest{
				UserId: tt.userID,
			})

			st, _ := status.FromError(err)

			if st.Code() != tt.wantCode {
				t.Errorf("code = %v, want %v", st.Code(), tt.wantCode)
			}

			if tt.wantCode == codes.OK {
				if resp == nil {
					t.Fatal("expected non-nil response")
				}
				if resp.GetId() != tt.wantID {
					t.Errorf("id = %v, want %v", resp.GetId(), tt.wantID)
				}
			}
		})
	}
}
