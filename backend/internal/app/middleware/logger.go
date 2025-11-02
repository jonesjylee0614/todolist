package middleware

import (
    "time"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

func Logger(log *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        method := c.Request.Method

        c.Next()

        status := c.Writer.Status()
        latency := time.Since(start)

        log.Info("http request",
            zap.String("method", method),
            zap.String("path", path),
            zap.Int("status", status),
            zap.Duration("latency", latency),
            zap.String("request_id", GetRequestID(c)),
        )
    }
}

