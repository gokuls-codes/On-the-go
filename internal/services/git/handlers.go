package git

import (
	"context"
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
	var payload struct{
		Repository struct{
			Name string `json:"name"`
		} `json:"repository"`
	}
	err := (&echo.DefaultBinder{}).BindBody(c, &payload)

	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}

	fmt.Println("Repository: ", payload.Repository.Name)

	// jsonBytes, err := json.MarshalIndent(payload, "", "  ")
	// if err != nil {
	// 	return c.JSON(500, map[string]string{"error": "Failed to marshal payload"})
	// }

	// jsonString := string(jsonBytes)

	// fmt.Println("Git Push Payload:", jsonBytes)

	headers := c.Request().Header

	fmt.Println("Headers:", headers.Get("X-Hub-Signature"))

	cmd := exec.Command("git", "pull")
    cmd.Dir = "../test-docker-project"

    output, err := cmd.CombinedOutput()
    if err != nil {
		fmt.Println("Error executing git pull:", string(output))
        return c.JSON(500, map[string]string{"error": fmt.Sprintf("Error executing git pull: %s", string(output))})
    }

	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to create Docker client"})
	}

	defer apiClient.Close()

	containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	var containerId string

	for _, container := range containers {
		if container.Image == "test-docker-project" {
			containerId = container.ID
			break
		}
	}

	fmt.Println("Found existing container ID:", containerId)

	if containerId != "" {
		err = apiClient.ContainerStop(context.Background(), containerId, container.StopOptions{})
		if err != nil {
			fmt.Println("Error stopping container:", err)
			return c.JSON(500, map[string]string{"error": err.Error()})
		}

		err = apiClient.ContainerRemove(context.Background(), containerId, container.RemoveOptions{})
		if err != nil {
			fmt.Println("Error removing container:", err)
			return c.JSON(500, map[string]string{"error": err.Error()})
		}
	}

	buildContext, err := utils.TarDirectory("../test-docker-project")
	if err != nil {
		fmt.Println("Error creating build context:", err)
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	buildOptions := build.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags: []string{"test-docker-project"},
		Remove: true,
	}

	_, err = apiClient.ImageBuild(context.Background(), buildContext, buildOptions)

	if err != nil {
		fmt.Println("Error building image:", err)
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	// decoder := json.NewDecoder(response.Body)
	// for decoder.More() {
	// 	var msg map[string]interface{}
	// 	if err := decoder.Decode(&msg); err != nil {
	// 		return c.JSON(500, map[string]string{"error": err.Error()})
	// 	}
	// 	if stream, ok := msg["stream"].(string); ok {
	// 		fmt.Println(c.Response().Writer, "<div>%s</div>", stream)
	// 		c.Response().Flush()
	// 	}
	// 	// time.Sleep(200 * time.Millisecond)
	// }

	containerResp, err := apiClient.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: "test-docker-project",
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				nat.Port("3000/tcp"): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: "8888",
					},
				},
			},
		},
		&network.NetworkingConfig{
		},
		nil, "",
	)

	if err != nil {
		fmt.Println("Error creating container:", err)
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	err = apiClient.ContainerStart(context.Background(), containerResp.ID, container.StartOptions{})
	
	if err != nil {
		fmt.Println("Error starting container:", err)
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	return c.JSON(200, map[string]string{"message": "Webhook successful"})
}