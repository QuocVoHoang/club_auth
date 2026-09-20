package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Total   interface{} `json:"total,omitempty"`
}

func DataResponse(ctx *gin.Context, status int, data interface{}) {
	ctx.JSON(status, Response{
		Code: status,
		Data: data,
	})
}

func ErrorResponse(ctx *gin.Context, err error) {
	var statusErr interface {
		StatusCode() int
		Error() string
	}
	if errors.As(err, &statusErr) {
		ctx.JSON(statusErr.StatusCode(), Response{
			Code:    statusErr.StatusCode(),
			Message: statusErr.Error(),
		})
		return
	}

	ctx.JSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: "internal server error",
	})
}
