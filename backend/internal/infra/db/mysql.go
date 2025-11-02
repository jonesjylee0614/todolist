package db

import (
    "fmt"

    "go.uber.org/zap"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"

    "todolist/backend/internal/domain/task"
    "todolist/backend/internal/domain/undo"
    "todolist/backend/internal/infra/config"
)

func Connect(cfg *config.Config, log *zap.Logger) (*gorm.DB, error) {
    gormCfg := &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    }

    db, err := gorm.Open(mysql.Open(cfg.Database.DSN), gormCfg)
    if err != nil {
        return nil, fmt.Errorf("gorm open: %w", err)
    }

    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("sql db: %w", err)
    }

    sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
    sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
    sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

    return db, nil
}

func AutoMigrate(db *gorm.DB, log *zap.Logger) error {
    if err := db.AutoMigrate(&task.Task{}, &undo.TaskOperation{}, &task.ActivityLog{}); err != nil {
        return fmt.Errorf("auto migrate: %w", err)
    }
    return nil
}

