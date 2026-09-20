package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-org/go-base/internal/domain/usecase"
	"github.com/your-org/go-base/internal/domain/usecase/dto"
	"github.com/your-org/go-base/internal/infrastructure/payload"
	"github.com/your-org/go-base/internal/infrastructure/response"
	"github.com/your-org/go-base/pkg/http_error"
)

type AuthHandler struct {
	registerUsecase usecase.RegisterUsecase
	loginUsecase    usecase.LoginUsecase
}

func NewAuthHandler(
	registerUsecase usecase.RegisterUsecase,
	loginUsecase usecase.LoginUsecase,
) *AuthHandler {
	return &AuthHandler{
		registerUsecase: registerUsecase,
		loginUsecase:    loginUsecase,
	}
}

// Register user
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body payload.RegisterRequest true "User registration data"
// @Success 200 {object} response.Response{data=response.AuthUser} "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 409 {object} response.Response "Conflict"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(ctx *gin.Context) {
	var req payload.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(ctx, http_error.BadRequestError("invalid request body"))
		return
	}

	result, err := h.registerUsecase.Do(
		ctx.Request.Context(),
		dto.RegisterRequest{
			Email:     req.Email,
			Phone:     req.Phone,
			FullName:  req.FullName,
			Role:      req.Role,
			Password:  req.Password,
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

// Login user
// @Summary Login user
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body payload.LoginRequest true "Login data"
// @Success 200 {object} response.Response{data=response.LoginResponse} "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(ctx *gin.Context) {
	var req payload.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(ctx, http_error.BadRequestError("invalid request body"))
		return
	}

	result, err := h.loginUsecase.Do(
		ctx.Request.Context(),
		dto.LoginRequest{
			Email:    req.Email,
			Password: req.Password,
		},
	)
	if err != nil {
		response.ErrorResponse(ctx, err)
		return
	}

	response.DataResponse(ctx, http.StatusOK, response.NewLoginResponse(*result))
}
