package docker

import "github.com/labstack/echo/v4"

type Handler struct{

}

func (h *Handler) RegisterRoutes(group *echo.Group) {
	group.GET("/containers", h.listContainers)
	group.GET("/containers/create", h.createContainerPage)
	group.POST("/containers", h.createContainer)
	group.GET("/images", h.listImages)

	group.GET("/projects/new", h.newProjectPage)
	group.POST("/projects", h.createProject)
}