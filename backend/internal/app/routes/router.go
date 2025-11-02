package routes

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "gorm.io/gorm"

    "todolist/backend/internal/app/handler"
    "todolist/backend/internal/app/middleware"
    "todolist/backend/internal/domain/task"
    "todolist/backend/internal/domain/undo"
    "todolist/backend/internal/infra/config"
    "todolist/backend/internal/pkg/response"
    "todolist/backend/internal/repository"
)

func SetupRouter(cfg *config.Config, log *zap.Logger, db *gorm.DB) *gin.Engine {
    if cfg.App.Env == "production" {
        gin.SetMode(gin.ReleaseMode)
    }

    engine := gin.New()
    engine.Use(middleware.RequestID(), middleware.Logger(log), middleware.Recovery(log), middleware.CORS(cfg.CORS))

    engine.GET("/healthz", func(c *gin.Context) {
        sqlDB, err := db.DB()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
            return
        }
        if err := sqlDB.Ping(); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })

    taskRepo := repository.NewTaskRepository(db)
    undoRepo := repository.NewUndoRepository(db)

    undoService := undo.NewService(undoRepo, taskRepo, cfg.Undo.TTL, log)
    taskService := task.NewService(taskRepo, undoService, log)

    taskHandler := handler.NewTaskHandler(taskService)
    undoHandler := handler.NewUndoHandler(undoService)

    api := engine.Group("/api/v1")
    {
        api.GET("/tasks", taskHandler.List)
        api.POST("/tasks", taskHandler.Create)
        api.GET("/tasks/:uuid", taskHandler.Get)
        api.PATCH("/tasks/:uuid", taskHandler.Update)
        api.PATCH("/tasks/:uuid/status", taskHandler.UpdateStatus)
        api.POST("/tasks/:uuid/complete", taskHandler.Complete)
        api.DELETE("/tasks/:uuid", taskHandler.Delete)

        api.POST("/tasks/bulk/move", taskHandler.BulkMove)
        api.POST("/tasks/bulk/complete", taskHandler.BulkComplete)
        api.POST("/tasks/bulk/delete", taskHandler.BulkDelete)
        api.POST("/tasks/order", taskHandler.UpdateOrder)

        api.POST("/undo", undoHandler.Undo)
    }

    engine.NoRoute(func(c *gin.Context) {
        response.NotFound(c, "route not found")
    })

    return engine
}

