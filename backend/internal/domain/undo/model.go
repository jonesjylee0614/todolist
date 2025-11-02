package undo

import "time"

type Action string

const (
    ActionCreate       Action = "create"
    ActionUpdate       Action = "update"
    ActionMove         Action = "move"
    ActionComplete     Action = "complete"
    ActionDelete       Action = "delete"
    ActionBulkMove     Action = "bulk_move"
    ActionBulkComplete Action = "bulk_complete"
    ActionBulkDelete   Action = "bulk_delete"
    ActionResort       Action = "resort"
)

type Scope string

const (
    ScopeSingle Scope = "single"
    ScopeBulk   Scope = "bulk"
)

type TaskOperation struct {
    ID          uint64    `gorm:"primaryKey;autoIncrement"`
    Token       string    `gorm:"type:char(26);uniqueIndex"`
    Action      Action    `gorm:"size:32;not null"`
    Scope       Scope     `gorm:"size:16;not null"`
    TaskIDs     string    `gorm:"type:json;not null"`
    BeforeState string    `gorm:"type:json"`
    AfterState  string    `gorm:"type:json"`
    ExpireAt    time.Time `gorm:"not null"`
    ConsumedAt  *time.Time
    CreatedAt   time.Time `gorm:"not null;autoCreateTime"`
}

func (op TaskOperation) IsExpired(now time.Time) bool {
    return now.After(op.ExpireAt)
}

func (op TaskOperation) IsConsumed() bool {
    return op.ConsumedAt != nil
}

