package git

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo/v4"
)


func (h *Handler) gitPush(c echo.Context) error {
	var payload interface{}
	err := (&echo.DefaultBinder{}).BindBody(c, &payload)

	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}

	jsonBytes, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to marshal payload"})
	}

	jsonString := string(jsonBytes)

	fmt.Println("Git Push Payload:", jsonString)
	return c.JSON(200, map[string]string{"message": "Git push successful"})
}