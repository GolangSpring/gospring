package controller

import (
	"fmt"
	"github.com/GolangSpring/gospring/application"
	"github.com/go-fuego/fuego"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"net/http"
)

var _ application.IController = (*SystemController)(nil)

func formatPercentage(value float64) string {
	return fmt.Sprintf("%.2f%%", value)
}

func GetSystemMetrics() map[string]string {
	metrics := make(map[string]string)

	cpuPercentages, err := cpu.Percent(0, false)
	if err == nil && len(cpuPercentages) > 0 {
		metrics["cpu"] = formatPercentage(cpuPercentages[0])
	} else {
		metrics["cpu"] = "error"
	}

	// Get memory usage
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		metrics["memory"] = formatPercentage(memInfo.UsedPercent)
	} else {
		metrics["memory"] = "error"
	}

	// Get disk usage
	diskInfo, err := disk.Usage("/")
	if err == nil {
		metrics["disk"] = formatPercentage(diskInfo.UsedPercent)
	} else {
		metrics["disk"] = "error"
	}

	return metrics
}

type SystemController struct{}

func NewSystemController() *SystemController {
	return &SystemController{}
}

func (controller *SystemController) Routes(server *fuego.Server) {
	fuego.Get(server, "/api-public/health", controller.Health)

}

func (controller *SystemController) Middlewares() []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{}
}

func (controller *SystemController) Health(c fuego.ContextNoBody) (map[string]string, error) {
	return GetSystemMetrics(), nil
}
