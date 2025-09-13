package docker

import (
	"github.com/gokuls-codes/on-the-go/internal/db"
	"github.com/gokuls-codes/on-the-go/internal/messageq"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Store    *db.Store
	MessageQ *messageq.MessageQ
}

func (h *Handler) RegisterRoutes(group *echo.Group) {
	group.GET("/containers", h.listContainers)
	group.GET("/containers/create", h.createContainerPage)
	group.POST("/containers", h.createContainer)
	group.GET("/images", h.listImages)

	group.GET("/projects/new", h.newProjectPage)
	group.POST("/projects", h.createProject)

	group.GET("/projects/new/env-var-row", h.getEnvVarRow)
}
