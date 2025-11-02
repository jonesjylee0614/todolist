package task

import (
    "time"

    "gorm.io/gorm"
)

type Status string

const (
    StatusNow     Status = "now"
    StatusFuture  Status = "future"
    StatusHistory Status = "history"
)

type Task struct {
    ID          uint64         `gorm:"primaryKey;autoIncrement"`
    UUID        string         `gorm:"type:char(36);uniqueIndex"`
    Title       string         `gorm:"size:255;not null"`
    Notes       *string        `gorm:"type:text"`
    Deadline    *time.Time     `gorm:"type:date"`
    Status      Status         `gorm:"type:enum('now','future','history');not null"`
    SortWeight  int64          `gorm:"not null"`
    CreatedAt   time.Time      `gorm:"not null;autoCreateTime"`
    UpdatedAt   time.Time      `gorm:"not null;autoUpdateTime"`
    CompletedAt *time.Time     `gorm:"type:datetime"`
    DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Snapshot struct {
    UUID        string     `json:"uuid"`
    Title       string     `json:"title"`
    Notes       *string    `json:"notes"`
    Deadline    *time.Time `json:"deadline"`
    Status      Status     `json:"status"`
    SortWeight  int64      `json:"sortWeight"`
    CreatedAt   time.Time  `json:"createdAt"`
    UpdatedAt   time.Time  `json:"updatedAt"`
    CompletedAt *time.Time `json:"completedAt"`
}

type ListFilter struct {
    Status    *Status
    Keyword   string
    Page      int
    PageSize  int
    OrderDesc bool
}

type ActivityLog struct {
    ID        uint64    `gorm:"primaryKey;autoIncrement"`
    TaskUUID  string    `gorm:"type:char(36);index"`
    Action    string    `gorm:"size:32;not null"`
    Payload   string    `gorm:"type:json"`
    Actor     string    `gorm:"size:64;not null"`
    CreatedAt time.Time `gorm:"not null;autoCreateTime"`
}

func FromSnapshot(s Snapshot) *Task {
    return &Task{
        UUID:        s.UUID,
        Title:       s.Title,
        Notes:       s.Notes,
        Deadline:    s.Deadline,
        Status:      s.Status,
        SortWeight:  s.SortWeight,
        CreatedAt:   s.CreatedAt,
        UpdatedAt:   s.UpdatedAt,
        CompletedAt: s.CompletedAt,
    }
}

func (t *Task) ToSnapshot() Snapshot {
    return Snapshot{
        UUID:        t.UUID,
        Title:       t.Title,
        Notes:       t.Notes,
        Deadline:    t.Deadline,
        Status:      t.Status,
        SortWeight:  t.SortWeight,
        CreatedAt:   t.CreatedAt,
        UpdatedAt:   t.UpdatedAt,
        CompletedAt: t.CompletedAt,
    }
}

func IsValidStatus(status Status) bool {
    switch status {
    case StatusNow, StatusFuture, StatusHistory:
        return true
    default:
        return false
    }
}

