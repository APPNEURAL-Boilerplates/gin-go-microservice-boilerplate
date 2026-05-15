package middleware

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/example/gin-go-microservice/internal/common"
)

func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		logger.Error("panic recovered",
			"request_id", common.RequestID(c),
			"error", fmt.Sprint(recovered),
		)

		common.Error(c, common.Internal("internal server error", nil))
	})
}
