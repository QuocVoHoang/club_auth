package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	infracontext "github.com/your-org/go-base/internal/infrastructure/context"
)

func TestRequireRoles(t *testing.T) {
	tests := []struct {
		name       string
		role       int
		allowed    []int
		wantStatus int
		wantCalled bool
	}{
		{name: "allowed", role: 2, allowed: []int{1, 2}, wantStatus: 200, wantCalled: true},
		{name: "denied", role: 3, allowed: []int{1, 2}, wantStatus: 403},
		{name: "missing role", allowed: []int{1, 2}, wantStatus: 403},
		{name: "empty permissions", role: 1, wantStatus: 403},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()
			called := false
			router.GET("/test", func(ctx *gin.Context) {
				infracontext.SetRole(ctx, tt.role)
				RequireRoles(tt.allowed...)(ctx)
				if !ctx.IsAborted() {
					called = true
				}
			})

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			router.ServeHTTP(recorder, req)
			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.wantStatus)
			}
			if called != tt.wantCalled {
				t.Fatalf("called = %t, want %t", called, tt.wantCalled)
			}
		})
	}
}
