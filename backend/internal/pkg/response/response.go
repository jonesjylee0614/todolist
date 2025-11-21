package response

import "github.com/gin-gonic/gin"

type Envelope struct {
    Code      int         `json:"code"`
    Message   string      `json:"message"`
    Data      interface{} `json:"data,omitempty"`
    UndoToken string      `json:"undoToken,omitempty"`
}

func Success(c *gin.Context, data interface{}, undoToken ...string) {
    resp := Envelope{
        Code:    0,
        Message: "ok",
        Data:    data,
    }
    if len(undoToken) > 0 && undoToken[0] != "" {
        resp.UndoToken = undoToken[0]
    }
    c.JSON(200, resp)
}

func Created(c *gin.Context, data interface{}, undoToken ...string) {
    resp := Envelope{
        Code:    0,
        Message: "ok",
        Data:    data,
    }
    if len(undoToken) > 0 && undoToken[0] != "" {
        resp.UndoToken = undoToken[0]
    }
    c.JSON(201, resp)
}

func BadRequest(c *gin.Context, msg string) {
    c.JSON(400, Envelope{Code: 40001, Message: msg})
}

func NotFound(c *gin.Context, msg string) {
    c.JSON(404, Envelope{Code: 40400, Message: msg})
}

func Conflict(c *gin.Context, msg string) {
    c.JSON(409, Envelope{Code: 40900, Message: msg})
}

func Gone(c *gin.Context, msg string) {
    c.JSON(410, Envelope{Code: 41000, Message: msg})
}

func InternalServerError(c *gin.Context, msg string) {
    c.JSON(500, Envelope{Code: 50000, Message: msg})
}


func Error(c *gin.Context, err error) {
	if err == nil {
		Success(c, nil)
		return
	}

	// Check for specific error types or messages
	// This is a simple implementation; in a real app, you might use errors.Is or custom error types
	msg := err.Error()
	switch msg {
	case "task not found":
		NotFound(c, msg)
	case "invalid status", "invalid deadline format", "invalid completed time", "empty ids", "ordered list empty":
		BadRequest(c, msg)
	case "undo token not found", "undo token expired", "undo token consumed":
		Gone(c, msg)
	default:
		InternalServerError(c, msg)
	}
}
