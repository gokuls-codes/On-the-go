package system

import (
	"bytes"
	"errors"
	"log"
	"time"

	"github.com/gokuls-codes/on-the-go/internal/utils"
	"github.com/gokuls-codes/on-the-go/internal/web/templates/components"
	"github.com/labstack/echo/v4"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"

	"fmt"
)

func (h *Handler) systemInfo(c echo.Context) error {

	v, _ := mem.VirtualMemory()
	fmt.Printf("Total: %v, Available: %v, UsedPercent:%f%%\n", v.Total, v.Available, v.UsedPercent)
	fmt.Println(v)

	physicalCnt, _ := cpu.Counts(false)
	logicalCnt, _ := cpu.Counts(true)
	fmt.Printf("physical count:%d logical count:%d\n", physicalCnt, logicalCnt)

	totalPercent, _ := cpu.Percent(1*time.Second, false)
	perPercents, _ := cpu.Percent(1*time.Second, true)
	fmt.Printf("total percent:%v per percents:%v", totalPercent, perPercents)

	return c.JSON(200, map[string]string{"status": "ok", "memory_used_percent": fmt.Sprintf("%f", v.UsedPercent), "cpu_count": fmt.Sprintf("%d", logicalCnt), "cpu_used_percent": fmt.Sprintf("%f", totalPercent[0])})
}

func (h *Handler) memoryInfoSSe(c echo.Context) error {

	log.Printf("SSE client connected, ip: %v", c.RealIP())

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	var buf bytes.Buffer

	memory, _ := mem.VirtualMemory()
	totalCpuPercent, _ := cpu.Percent(1*time.Second, false)
	diskInfo, _ := disk.Usage("/")

	err := components.SystemStatsComponent(totalCpuPercent[0], memory.UsedPercent, diskInfo.UsedPercent).Render(c.Request().Context(), &buf)

	if err != nil {
		log.Printf("Error rendering template: %v\n", err)
		return errors.New("error rendering template")
	}

	event := utils.Event{
		Event: []byte("system-stats"),
		Data:  buf.Bytes(),
	}
	if err := event.MarshalTo(w); err != nil {
		return err
	}

	w.Flush()

	for {
		select {
		case <-c.Request().Context().Done():
			log.Printf("SSE client disconnected, ip: %v", c.RealIP())
			return nil

		case <-ticker.C:
			var buf bytes.Buffer

			memory, _ := mem.VirtualMemory()
			totalCpuPercent, _ := cpu.Percent(1*time.Second, false)
			diskInfo, _ := disk.Usage("/")

			err := components.SystemStatsComponent(totalCpuPercent[0], memory.UsedPercent, diskInfo.UsedPercent).Render(c.Request().Context(), &buf)

			if err != nil {
				log.Printf("Error rendering template: %v\n", err)
				return errors.New("error rendering template")
			}

			event := utils.Event{
				Event: []byte("system-stats"),
				Data:  buf.Bytes(),
			}
			if err := event.MarshalTo(w); err != nil {
				return err
			}

			w.Flush()
		}

	}

}
