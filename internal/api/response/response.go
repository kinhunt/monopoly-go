// internal/api/response/response.go
package response

import (
	"encoding/json"
	"monopoly/pkg/utils"
	"net/http"
)

type Response struct {
	Success bool                `json:"success"`
	Data    interface{}         `json:"data,omitempty"`
	Error   utils.ErrorResponse `json:"error,omitempty"` // 移除指针
}

func JSON(w http.ResponseWriter, statusCode int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func Success(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

func Error(err error) Response {
	return Response{
		Success: false,
		Error:   utils.NewErrorResponse(err), // 直接使用返回值
	}
}

// JsonError 是一个便捷函数，用于返回错误响应
func JsonError(w http.ResponseWriter, err error) {
	statusCode := utils.HTTPStatusFromError(err)
	JSON(w, statusCode, Error(err))
}
