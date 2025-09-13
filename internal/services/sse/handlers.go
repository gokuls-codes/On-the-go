package sse

import (
	"bytes"
	"log"

	"github.com/gokuls-codes/on-the-go/internal/utils"
	"github.com/gokuls-codes/on-the-go/internal/web/templates/components"
	"github.com/labstack/echo/v4"
)

func (h *Handler) getProjectLogsSSE(c echo.Context) error {
	projectChan, err := h.MessageQ.GetChannel(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(404, "Project not found")
	}

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	for {
		select {
		case <-c.Request().Context().Done():
			log.Printf("SSE client disconnected, ip: %v", c.RealIP())
			return nil

		case msg, ok := <-projectChan:
			if !ok {
				log.Printf("Project %s channel closed", c.Param("id"))
				return c.String(200, "Channel closed")
			}

			var buf bytes.Buffer

			err := components.LogLine(msg).Render(c.Request().Context(), &buf)

			if err != nil {
				log.Println("Error rendering log line template: ", err)
				return echo.NewHTTPError(500, "Internal server error")
			}

			event := utils.Event{
				Event: []byte("project-log"),
				Data:  buf.Bytes(),
			}

			if err := event.MarshalTo(w); err != nil {
				log.Printf("Error writing to client: %v", err)
				return err
			}

			w.Flush()
		}
	}
}
