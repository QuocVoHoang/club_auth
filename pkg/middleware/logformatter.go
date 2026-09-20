package middleware

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func LogFormatterJSON(params gin.LogFormatterParams) string {
	latency := params.Latency
	if latency > time.Minute {
		latency = latency.Truncate(time.Second)
	}

	requestMeta := map[string]interface{}{
		"response_time": params.TimeStamp.Format("2006/01/02 - 15:04:05"),
		"status_code":   params.StatusCode,
		"latency":       latency.String(),
		"client_ip":     params.ClientIP,
		"path":          params.Path,
		"method":        params.Method,
		"error_message": params.ErrorMessage,
	}
	if params.Keys != nil {
		if requestID, ok := params.Keys[RequestIDKey]; ok {
			requestMeta[RequestIDKey] = requestID
		}
	}

	raw := map[string]interface{}{
		"request_meta": requestMeta,
		"msg":          "finish route",
	}

	//nolint:errchkjson
	bytes, _ := json.Marshal(&raw)

	return fmt.Sprintf("%s\n", string(bytes))
}
