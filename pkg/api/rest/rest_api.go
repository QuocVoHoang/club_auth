package rest

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/your-org/go-base/pkg/api"
)

type API interface {
	Call(ctx context.Context, method, url string, body io.Reader, reqOptions ...api.RequestOption) ([]byte, map[string][]string, error)
}

type restAPI struct {
	client api.Client
}

func NewAPI(timeoutSeconds int) API {
	return &restAPI{
		client: &http.Client{Timeout: time.Duration(timeoutSeconds) * time.Second},
	}
}

func (r *restAPI) Call(ctx context.Context, method, url string, body io.Reader, reqOptions ...api.RequestOption) ([]byte, map[string][]string, error) {
	return api.Call(ctx, r.client, method, url, body, reqOptions...)
}
