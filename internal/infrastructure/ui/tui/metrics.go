package tui

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemMetricsCollector implements MetricsCollector interface
type SystemMetricsCollector struct {
	metrics    ResourceMetrics
	stopChan   chan bool
	updateChan chan ResourceMetrics
}

// NewSystemMetricsCollector creates a new metrics collector
func NewSystemMetricsCollector() *SystemMetricsCollector {
	return &SystemMetricsCollector{
		stopChan:   make(chan bool),
		updateChan: make(chan ResourceMetrics, 1),
	}
}

// GetMetrics returns the current metrics
func (c *SystemMetricsCollector) GetMetrics() ResourceMetrics {
	return c.metrics
}

// Start begins collecting metrics
func (c *SystemMetricsCollector) Start() error {
	go c.collectMetrics()
	return nil
}

// Stop stops collecting metrics
func (c *SystemMetricsCollector) Stop() error {
	close(c.stopChan)
	return nil
}

// GetUpdateChannel returns the channel for metric updates
func (c *SystemMetricsCollector) GetUpdateChannel() <-chan ResourceMetrics {
	return c.updateChan
}

// collectMetrics runs the metrics collection loop
func (c *SystemMetricsCollector) collectMetrics() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			metrics := c.gatherMetrics()
			c.metrics = metrics
			select {
			case c.updateChan <- metrics:
			default:
				// Channel is full, skip this update
			}
		case <-c.stopChan:
			return
		}
	}
}

// gatherMetrics collects current system metrics
func (c *SystemMetricsCollector) gatherMetrics() ResourceMetrics {
	cpuPercents, err := cpu.Percent(0, false)
	var cpuUsage float64
	if err == nil && len(cpuPercents) > 0 {
		cpuUsage = cpuPercents[0]
	}

	memInfo, err := mem.VirtualMemory()
	var memUsage, memUsed, memTotal float64
	if err == nil {
		memUsage = memInfo.UsedPercent
		memUsed = float64(memInfo.Used)
		memTotal = float64(memInfo.Total)
	}

	return ResourceMetrics{
		CPUPercent:    cpuUsage,
		MemoryPercent: memUsage,
		MemoryUsed:    uint64(memUsed),
		MemoryTotal:   uint64(memTotal),
		LastUpdate:    time.Now(),
	}
}
