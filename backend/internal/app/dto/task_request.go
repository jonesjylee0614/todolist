package dto

import (
	"encoding/json"
	"time"
)

type CreateTaskRequest struct {
	Title      string  `json:"title" binding:"required,min=1,max=255"`
	Notes      *string `json:"notes"`
	Deadline   *string `json:"deadline"`
	Status     *string `json:"status" binding:"omitempty,oneof=now future history"`
	SortWeight *int64  `json:"sortWeight"`
	ParentUUID *string `json:"parentUuid"`
}

type UpdateTaskRequest struct {
	Title    *string        `json:"title"`
	Notes    NullableString `json:"notes"`
	Deadline NullableDate   `json:"deadline"`
}

type StatusUpdateRequest struct {
	Status      string  `json:"status" binding:"required,oneof=now future history"`
	SortWeight  *int64  `json:"sortWeight"`
	CompletedAt *string `json:"completedAt"`
}

type BulkOperationRequest struct {
	IDs []string `json:"ids" binding:"required,min=1,dive,required"`
}

type BulkMoveRequest struct {
	IDs    []string `json:"ids" binding:"required,min=1,dive,required"`
	Status string   `json:"targetStatus" binding:"required,oneof=now future history"`
}

type OrderUpdateRequest struct {
	Status     string   `json:"status" binding:"required,oneof=now future history"`
	OrderedIDs []string `json:"orderedIds" binding:"required,min=1,dive,required"`
}

type UndoRequest struct {
	Token string `json:"token" binding:"required"`
}

type NullableString struct {
	Value *string
	Set   bool
}

func (ns *NullableString) UnmarshalJSON(data []byte) error {
	ns.Set = true
	if string(data) == "null" {
		ns.Value = nil
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	ns.Value = &s
	return nil
}

type NullableDate struct {
	Value *time.Time
	Set   bool
}

func (nd *NullableDate) UnmarshalJSON(data []byte) error {
	nd.Set = true
	if string(data) == "null" {
		nd.Value = nil
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		nd.Value = nil
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	nd.Value = &t
	return nil
}

type ListQuery struct {
	Status   string `form:"status"`
	Keyword  string `form:"keyword"`
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
}
