package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

const requestIDKey = "X-Request-ID"

func RequestID() gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.GetHeader(requestIDKey)
        if id == "" {
            id = uuid.NewString()
        }
        c.Set(requestIDKey, id)
        c.Writer.Header().Set(requestIDKey, id)
        c.Next()
    }
}

func GetRequestID(c *gin.Context) string {
    if v, ok := c.Get(requestIDKey); ok {
        if id, ok := v.(string); ok {
            return id
        }
    }
    return ""
}

