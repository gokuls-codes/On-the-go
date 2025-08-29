package docker

import (
	"github.com/gokuls-codes/on-the-go/internal/utils"
	"github.com/gokuls-codes/on-the-go/internal/web/templates/pages"
	"github.com/labstack/echo/v4"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

func (h *Handler) listContainers(c echo.Context) error {

	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to create Docker client"})
	}

	defer apiClient.Close()

	containers, err := apiClient.ContainerList(c.Request().Context(), container.ListOptions{All: true})
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	// resp := new([]string)
	// for _, container := range containers {
	// 	*resp = append(*resp, fmt.Sprintf("%s %s (status: %s)\n", container.Names[0], container.Image, container.Status))
	// }

	// return c.JSON(200, resp)
	return utils.Render(c, pages.Containers(containers))
}

func (h *Handler) createContainer(c echo.Context) error {

	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to create Docker client"})
	}

	defer apiClient.Close()

	_, err = apiClient.ContainerCreate(
		c.Request().Context(),
		&container.Config{
			Image: "hello-world",
		},
		&container.HostConfig{},
		&network.NetworkingConfig{},
		nil, // platform (use nil for default)
		"",  // container name (empty for auto-generated)
	)

	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	return c.JSON(200, map[string]string{"message": "Container created successfully"})
}