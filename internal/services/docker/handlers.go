package docker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/docker/go-connections/nat"
	"github.com/gokuls-codes/on-the-go/internal/types"
	"github.com/gokuls-codes/on-the-go/internal/utils"
	"github.com/gokuls-codes/on-the-go/internal/web/templates/components"
	"github.com/gokuls-codes/on-the-go/internal/web/templates/pages"
	"github.com/labstack/echo/v4"
	"github.com/moby/moby/api/types/build"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/image"
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

func (h *Handler) listImages(c echo.Context) error {

	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to create Docker client"})
	}

	defer apiClient.Close()

	images, err := apiClient.ImageList(c.Request().Context(), image.ListOptions{})
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	return utils.Render(c, pages.Images(images))
}

func (h *Handler) createContainerPage(c echo.Context) error {
	return utils.Render(c, pages.CreateContainerPage())
}

func (h *Handler) newProjectPage(c echo.Context) error {
	return utils.Render(c, pages.NewProjectPage())
}

func (h *Handler) createProject(c echo.Context) error {

	p := new(types.CreateProjectPayload)

	if err := c.Bind(p); err != nil {
		log.Println(err)
		return c.String(http.StatusBadRequest, "bad request")
	}

	targetDir := "../otg-projects/"
	repoName := utils.GetRepoName(p.GitHubURL)

	cmd := exec.Command("git", "clone", p.GitHubURL, targetDir+repoName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return c.JSON(500, map[string]string{"error": fmt.Sprintf("Git clone failed: %s", string(output))})
	}

	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to create Docker client"})
	}

	defer apiClient.Close()

	buildContext, err := utils.TarDirectory(targetDir + repoName)
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	buildOptions := build.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{repoName},
		Remove:     true,
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
			log.Println(stream)
		}
		// time.Sleep(200 * time.Millisecond)
	}

	containerResp, err := apiClient.ContainerCreate(
		c.Request().Context(),
		&container.Config{
			Image: repoName,
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				nat.Port(fmt.Sprintf("%d/tcp", p.ContainerPort)): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: fmt.Sprintf("%d", p.HostPort),
					},
				},
			},
		},
		&network.NetworkingConfig{},
		nil, "",
	)

	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	err = apiClient.ContainerStart(c.Request().Context(), containerResp.ID, container.StartOptions{})

	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	return nil
}

func (h *Handler) getEnvVarRow(c echo.Context) error {
	return utils.Render(c, components.EnvVariablesRow())
}
