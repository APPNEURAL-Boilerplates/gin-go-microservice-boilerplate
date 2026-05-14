package common

import (
	"github.com/gin-gonic/gin"
)

const RequestIDKey = "request_id"

type SuccessResponse struct {
	OK        bool   `json:"ok"`
	Data      any    `json:"data,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	OK        bool         `json:"ok"`
	Error     ErrorPayload `json:"error"`
	RequestID string       `json:"request_id,omitempty"`
}

func JSON(c *gin.Context, status int, data any) {
	c.JSON(status, SuccessResponse{
		OK:        true,
		Data:      data,
		RequestID: RequestID(c),
	})
}

func Error(c *gin.Context, err error) {
	appError := ToAppError(err)
	c.AbortWithStatusJSON(appError.Status, ErrorResponse{
		OK: false,
		Error: ErrorPayload{
			Code:    appError.Code,
			Message: appError.Message,
		},
		RequestID: RequestID(c),
	})
}

func RequestID(c *gin.Context) string {
	value, exists := c.Get(RequestIDKey)
	if !exists {
		return ""
	}

	requestID, ok := value.(string)
	if !ok {
		return ""
	}

	return requestID
}
