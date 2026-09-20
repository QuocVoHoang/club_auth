package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-org/go-base/internal/domain/usecase"
	"github.com/your-org/go-base/internal/domain/usecase/dto"
	infracontext "github.com/your-org/go-base/internal/infrastructure/context"
	"github.com/your-org/go-base/internal/infrastructure/payload"
	"github.com/your-org/go-base/internal/infrastructure/response"
	"github.com/your-org/go-base/pkg/http_error"
)

type UserHandler struct {
	getCurrentUserUsecase    usecase.GetCurrentUserUsecase
	updateCurrentUserUsecase usecase.UpdateCurrentUserUsecase
}

func NewUserHandler(
	getCurrentUserUsecase usecase.GetCurrentUserUsecase,
	updateCurrentUserUsecase usecase.UpdateCurrentUserUsecase,
) *UserHandler {
	return &UserHandler{
		getCurrentUserUsecase:    getCurrentUserUsecase,
		updateCurrentUserUsecase: updateCurrentUserUsecase,
	}
}

// Get current user
// @Summary Get current user
// @Description Get current authenticated user's profile
// @Tags users
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=response.AuthUser} "Success"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /api/v1/users/me [get]
func (h *UserHandler) GetMe(ctx *gin.Context) {
	result, err := h.getCurrentUserUsecase.Do(
		ctx.Request.Context(),
		dto.GetCurrentUserRequest{
			UserID: infracontext.GetUserID(ctx),
		},
	)
	if err != nil {
		response.ErrorResponse(ctx, err)
		return
	}

	response.DataResponse(ctx, http.StatusOK, response.NewAuthUser(*result))
}

// Update current user
// @Summary Update current user
// @Description Update current authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body payload.UpdateCurrentUserRequest true "User update data"
// @Success 200 {object} response.Response{data=response.AuthUser} "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 409 {object} response.Response "Conflict"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /api/v1/users/me [patch]
func (h *UserHandler) UpdateMe(ctx *gin.Context) {
	var req payload.UpdateCurrentUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(ctx, http_error.BadRequestError("invalid request body"))
		return
	}

	result, err := h.updateCurrentUserUsecase.Do(
		ctx.Request.Context(),
		dto.UpdateCurrentUserRequest{
			UserID:    infracontext.GetUserID(ctx),
			FullName:  req.FullName,
			Phone:     req.Phone,
			Avatar:    req.Avatar,
			Birthdate: req.Birthdate,
		},
	)
	if err != nil {
		response.ErrorResponse(ctx, err)
		return
	}

	response.DataResponse(ctx, http.StatusOK, response.NewAuthUser(*result))
}
