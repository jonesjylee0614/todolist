package handler

import (
    "github.com/gin-gonic/gin"

    "todolist/backend/internal/app/dto"
    "todolist/backend/internal/domain/undo"
    "todolist/backend/internal/pkg/response"
)

type UndoHandler struct {
    service *undo.Service
}

func NewUndoHandler(service *undo.Service) *UndoHandler {
    return &UndoHandler{service: service}
}

func (h *UndoHandler) Undo(c *gin.Context) {
    var req dto.UndoRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, err.Error())
        return
    }

    ids, nextToken, err := h.service.Undo(c.Request.Context(), req.Token)
    if err != nil {
        switch err {
        case undo.ErrTokenNotFound:
            response.Gone(c, "undo token not found")
        case undo.ErrTokenExpired:
            response.Gone(c, "undo token expired")
        case undo.ErrTokenConsumed:
            response.Gone(c, "undo token consumed")
        default:
            response.InternalServerError(c, err.Error())
        }
        return
    }

    response.Success(c, gin.H{"affectedIds": ids}, nextToken)
}

