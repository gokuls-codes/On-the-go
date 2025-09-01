package git

import (
	"context"
	"encoding/json"
	"log"
	"os/exec"

	"github.com/docker/go-connections/nat"
	"github.com/gokuls-codes/on-the-go/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/moby/moby/api/types/build"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/image"
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

	log.Println("Repository: ", payload.Repository.Name)

	// jsonBytes, err := json.MarshalIndent(payload, "", "  ")
	// if err != nil {
	// 	return c.JSON(500, map[string]string{"error": "Failed to marshal payload"})
	// }

	// jsonString := string(jsonBytes)

	// log.Println("Git Push Payload:", jsonBytes)

	headers := c.Request().Header

	log.Println("Headers:", headers.Get("X-Hub-Signature"))

	go func() {

	cmd := exec.Command("git", "pull")
    cmd.Dir = "../test-docker-project"

    output, err := cmd.CombinedOutput()
    if err != nil {
		log.Println("Error executing git pull:", string(output))
		return
    }
	log.Println("Git pull successful\nOutput:", string(output))

	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Println("Error creating Docker client:", err)
		return
	}

	defer apiClient.Close()

	containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		log.Println("Error listing containers:", err)
		return
	}

	var containerId string

	for _, container := range containers {
		if container.Image == "test-docker-project" {
			containerId = container.ID
			break
		}
	}

	log.Println("Found existing container ID:", containerId)

	if containerId != "" {
		err = apiClient.ContainerStop(context.Background(), containerId, container.StopOptions{})
		if err != nil {
			log.Println("Error stopping container:", err)
		return
		}
		log.Println("Container stopped successfully")

		err = apiClient.ContainerRemove(context.Background(), containerId, container.RemoveOptions{})
		if err != nil {
			log.Println("Error removing container:", err)
		return
		}
		log.Println("Container removed successfully")
	}

	images, err := apiClient.ImageList(context.Background(), image.ListOptions{All: true})
	if err != nil {
		log.Println("Error listing images:", err)
		return
	}

	var imageId string
	
	for _, image := range images {
		log.Println("Image:", image.RepoTags, image.ID)
		if len(image.RepoTags) > 0 && image.RepoTags[0] == "test-docker-project:latest" {
			imageId = image.ID
		}
	}

	log.Println("Found existing image ID:", imageId)

	if imageId != "" {
		_, err = apiClient.ImageRemove(context.Background(), imageId, image.RemoveOptions{Force: true, PruneChildren: true})
		if err != nil {
			log.Println("Error removing image:", err)
		return
		}
		log.Println("Image removed successfully")
	}

	buildContext, err := utils.TarDirectory("../test-docker-project")
	if err != nil {
		log.Println("Error creating build context:", err)
		return
	}

	buildOptions := build.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags: []string{"test-docker-project"},
		Remove: true,
	}

	imageResponse, err := apiClient.ImageBuild(context.Background(), buildContext, buildOptions)

	if err != nil {
		log.Println("Error building image:", err)
		return
	}
	
	log.Println("Image built successfully")



	decoder := json.NewDecoder(imageResponse.Body)
	for decoder.More() {
		var msg map[string]interface{}
		if err := decoder.Decode(&msg); err != nil {
		return
		}
		if stream, ok := msg["stream"].(string); ok {
			log.Println(stream)
		}
		// time.Sleep(200 * time.Millisecond)
	}

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
		log.Println("Error creating container:", err)
		return
	}
	log.Println("Container created successfully")

	err = apiClient.ContainerStart(context.Background(), containerResp.ID, container.StartOptions{})
	
	if err != nil {
		log.Println("Error starting container:", err)
		return
	}
	
	log.Println("Container started successfully")

}()

	return c.JSON(200, map[string]string{"message": "Webhook successful"})
}