package git

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/docker/go-connections/nat"
	"github.com/gokuls-codes/on-the-go/internal/db/sqlc"
	"github.com/gokuls-codes/on-the-go/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/moby/moby/api/types/build"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/image"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

func (h *Handler) gitPush(c echo.Context) error {
	var payload struct {
		Repository struct {
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
	repoName := payload.Repository.Name

	log.Println("Headers:", headers.Get("X-Hub-Signature"))

	project, err := h.Store.GetProjectByRepoName(c.Request().Context(), repoName)

	if err != nil {
		log.Println("Error fetching project:", err)
		return c.JSON(500, map[string]string{"error": "Project not found"})
	}

	log.Println("Project found:", project)

	go func() {

		cmd := exec.Command("git", "pull")
		cmd.Dir = "../otg-projects/" + repoName

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

		containerId := project.ContainerID.String
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

		imageId := project.ImageID.String

		log.Println("Found existing image ID:", imageId)

		if imageId != "" {
			_, err = apiClient.ImageRemove(context.Background(), imageId, image.RemoveOptions{Force: true, PruneChildren: true})
			if err != nil {
				log.Println("Error removing image:", err)
				return
			}
			log.Println("Image removed successfully")
		}

		buildContext, err := utils.TarDirectory("../" + repoName)
		if err != nil {
			log.Println("Error creating build context:", err)
			return
		}

		buildOptions := build.ImageBuildOptions{
			Dockerfile: "Dockerfile",
			Tags:       []string{repoName},
			Remove:     true,
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

		err = h.Store.Queries.UpdateImageId(context.Background(), sqlc.UpdateImageIdParams{
			ID:      project.ID,
			ImageID: sql.NullString{String: repoName, Valid: true},
		})

		if err != nil {
			log.Println("Updating image ID failed:", err)
			return
		}

		containerResp, err := apiClient.ContainerCreate(
			context.Background(),
			&container.Config{
				Image: repoName,
			},
			&container.HostConfig{
				PortBindings: nat.PortMap{
					nat.Port(fmt.Sprintf("%d/tcp", project.ContainerPort.Int64)): []nat.PortBinding{
						{
							HostIP:   "0.0.0.0",
							HostPort: fmt.Sprintf("%d", project.HostPort.Int64),
						},
					},
				},
			},
			&network.NetworkingConfig{},
			nil, "",
		)

		if err != nil {
			log.Println("Error creating container:", err)
			return
		}
		log.Println("Container created successfully")

		err = h.Store.Queries.UpdateContainerId(context.Background(), sqlc.UpdateContainerIdParams{
			ID:          project.ID,
			ContainerID: sql.NullString{String: containerResp.ID, Valid: true},
		})

		if err != nil {
			log.Println("Updating container ID failed:", err)
			return
		}

		err = apiClient.ContainerStart(context.Background(), containerResp.ID, container.StartOptions{})

		if err != nil {
			log.Println("Error starting container:", err)
			return
		}

		log.Println("Container started successfully")

	}()

	return c.JSON(200, map[string]string{"message": "Webhook successful"})
}
