package sse

import (
	"github.com/gokuls-codes/on-the-go/internal/messageq"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	MessageQ *messageq.MessageQ
}

func (h *Handler) RegisterRoutes(group *echo.Group) {
	group.GET("/project/:id/logs", h.getProjectLogsSSE)
}
