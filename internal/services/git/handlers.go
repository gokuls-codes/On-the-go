package git

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/docker/go-connections/nat"
	"github.com/gokuls-codes/on-the-go/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/moby/moby/api/types/build"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
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

	headers := c.Request().Header

	fmt.Println("Headers:", headers)

	cmd := exec.Command("git", "pull")
    cmd.Dir = "../test-docker-project"

    output, err := cmd.CombinedOutput()
    if err != nil {
        return c.JSON(500, map[string]string{"error": fmt.Sprintf("Error executing git pull: %s", string(output))})
    }

	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to create Docker client"})
	}

	defer apiClient.Close()

	buildContext, err := utils.TarDirectory("../test-docker-project")
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	buildOptions := build.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags: []string{"test-docker-project"},
		Remove: true,
	}

	response, err := apiClient.ImageBuild(c.Request().Context(), buildContext, buildOptions)

	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	decoder := json.NewDecoder(response.Body)
	for decoder.More() {
		var msg map[string]interface{}
		if err := decoder.Decode(&msg); err != nil {
			return c.JSON(500, map[string]string{"error": err.Error()})
		}
		if stream, ok := msg["stream"].(string); ok {
			fmt.Fprintf(c.Response().Writer, "<div>%s</div>", stream)
			c.Response().Flush()
		}
		// time.Sleep(200 * time.Millisecond)
	}

	containerResp, err := apiClient.ContainerCreate(
		c.Request().Context(),
		&container.Config{
			Image: "test-docker-project",
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				nat.Port("3000/tcp"): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: "3000",
					},
				},
			},
		},
		&network.NetworkingConfig{
		},
		nil, "",
	)

	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	err = apiClient.ContainerStart(c.Request().Context(), containerResp.ID, container.StartOptions{})
	
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	return c.JSON(200, map[string]string{"message": "Webhook successful"})
}