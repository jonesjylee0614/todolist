package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    "todolist/backend/internal/pkg/response"
)

func Recovery(log *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if r := recover(); r != nil {
                log.Error("panic recovered", zap.Any("error", r), zap.String("request_id", GetRequestID(c)))
                response.InternalServerError(c, "internal server error")
                c.AbortWithStatus(http.StatusInternalServerError)
            }
        }()
        c.Next()
    }
}

