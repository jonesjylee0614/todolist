package middleware

import (
    "time"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"

    "todolist/backend/internal/infra/config"
)

func CORS(cfg config.CORSConfig) gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     cfg.AllowOrigins,
        AllowMethods:     cfg.AllowMethods,
        AllowHeaders:     cfg.AllowHeaders,
        ExposeHeaders:    []string{"X-Request-ID"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    })
}

