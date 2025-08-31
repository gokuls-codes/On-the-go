package git

import (
	"fmt"

	"github.com/labstack/echo/v4"
)


func (h *Handler) gitPush(c echo.Context) error {
	var payload interface{}
	err := (&echo.DefaultBinder{}).BindBody(c, &payload)

	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}

	fmt.Println("Git Push Payload:", payload)
	return c.JSON(200, map[string]string{"message": "Git push successful"})
}