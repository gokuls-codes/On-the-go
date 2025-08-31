package git

import "github.com/labstack/echo/v4"

type Handler struct {

}

func (h *Handler) RegisterRoutes(group *echo.Group) {
	group.POST("", h.gitPush)
}
