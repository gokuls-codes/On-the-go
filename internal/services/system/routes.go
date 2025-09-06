package system

import (
	"github.com/labstack/echo/v4"
)

type Handler struct {
}

func (h *Handler) RegisterRoutes(group *echo.Group) {
	group.GET("/info", h.systemInfo)
	group.GET("/sse", h.memoryInfoSSe)
}
