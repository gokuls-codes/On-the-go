package git

import (
	"github.com/gokuls-codes/on-the-go/internal/db"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Store *db.Store
}

func (h *Handler) RegisterRoutes(group *echo.Group) {
	group.POST("/webhooks", h.gitPush)
}
