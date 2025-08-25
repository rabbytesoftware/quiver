package process

import (
	"context"
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/rabbytesoftware/quiver/internal/logger"
)

// ProcessInfo holds information about a running process
type ProcessInfo struct {
	PID       int       `json:"pid"`
	ArrowName string    `json:"arrow_name"`
	Command   string    `json:"command"`
	StartTime time.Time `json:"start_time"`
	Process   *os.Process `json:"-"` // Don't serialize the actual process
}

// ProcessTracker manages running processes for arrows
type ProcessTracker struct {
	logger    *logger.Logger
	processes map[string][]*ProcessInfo // arrow_name -> list of processes
	mutex     sync.RWMutex
}

// NewProcessTracker creates a new process tracker
func NewProcessTracker(logger *logger.Logger) *ProcessTracker {
	return &ProcessTracker{
		logger:    logger.WithService("process-tracker"),
		processes: make(map[string][]*ProcessInfo),
	}
}

// TrackProcess adds a process to be tracked for the given arrow
func (pt *ProcessTracker) TrackProcess(arrowName string, process *os.Process, command string) {
	pt.mutex.Lock()
	defer pt.mutex.Unlock()

	processInfo := &ProcessInfo{
		PID:       process.Pid,
		ArrowName: arrowName,
		Command:   command,
		StartTime: time.Now(),
		Process:   process,
	}

	if pt.processes[arrowName] == nil {
		pt.processes[arrowName] = make([]*ProcessInfo, 0)
	}

	pt.processes[arrowName] = append(pt.processes[arrowName], processInfo)
	pt.logger.Info("Tracking process PID %d for arrow %s: %s", process.Pid, arrowName, command)
}

// GetProcesses returns all processes for a given arrow
func (pt *ProcessTracker) GetProcesses(arrowName string) []*ProcessInfo {
	pt.mutex.RLock()
	defer pt.mutex.RUnlock()

	processes := pt.processes[arrowName]
	if processes == nil {
		return []*ProcessInfo{}
	}

	// Return a copy to avoid race conditions
	result := make([]*ProcessInfo, len(processes))
	copy(result, processes)
	return result
}

// GetAllProcesses returns all tracked processes
func (pt *ProcessTracker) GetAllProcesses() map[string][]*ProcessInfo {
	pt.mutex.RLock()
	defer pt.mutex.RUnlock()

	result := make(map[string][]*ProcessInfo)
	for arrowName, processes := range pt.processes {
		result[arrowName] = make([]*ProcessInfo, len(processes))
		copy(result[arrowName], processes)
	}
	return result
}

// StopProcesses stops all processes for a given arrow
func (pt *ProcessTracker) StopProcesses(arrowName string, graceful bool, timeout time.Duration) error {
	pt.mutex.Lock()
	defer pt.mutex.Unlock()

	processes := pt.processes[arrowName]
	if len(processes) == 0 {
		pt.logger.Info("No processes to stop for arrow %s", arrowName)
		return nil
	}

	pt.logger.Info("Stopping %d processes for arrow %s (graceful: %t)", len(processes), arrowName, graceful)

	var errors []string
	stoppedProcesses := make([]*ProcessInfo, 0)

	for _, processInfo := range processes {
		if err := pt.stopSingleProcess(processInfo, graceful, timeout); err != nil {
			pt.logger.Error("Failed to stop process PID %d: %v", processInfo.PID, err)
			errors = append(errors, fmt.Sprintf("PID %d: %v", processInfo.PID, err))
		} else {
			stoppedProcesses = append(stoppedProcesses, processInfo)
		}
	}

	// Remove stopped processes from tracking
	if len(stoppedProcesses) > 0 {
		remaining := make([]*ProcessInfo, 0)
		for _, processInfo := range processes {
			stopped := false
			for _, stoppedProcess := range stoppedProcesses {
				if processInfo.PID == stoppedProcess.PID {
					stopped = true
					break
				}
			}
			if !stopped {
				remaining = append(remaining, processInfo)
			}
		}
		
		if len(remaining) == 0 {
			delete(pt.processes, arrowName)
		} else {
			pt.processes[arrowName] = remaining
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to stop some processes: %v", errors)
	}

	pt.logger.Info("Successfully stopped all processes for arrow %s", arrowName)
	return nil
}

// stopSingleProcess stops a single process with optional graceful shutdown
func (pt *ProcessTracker) stopSingleProcess(processInfo *ProcessInfo, graceful bool, timeout time.Duration) error {
	pid := processInfo.PID
	process := processInfo.Process

	// Check if process is still running
	if !pt.isProcessRunning(pid) {
		pt.logger.Info("Process PID %d is already stopped", pid)
		return nil
	}

	if graceful {
		pt.logger.Info("Attempting graceful shutdown of process PID %d", pid)
		
		// Send SIGTERM for graceful shutdown
		if err := process.Signal(syscall.SIGTERM); err != nil {
			pt.logger.Warn("Failed to send SIGTERM to PID %d: %v", pid, err)
		} else {
			// Wait for graceful shutdown with timeout
			done := make(chan error, 1)
			go func() {
				_, err := process.Wait()
				done <- err
			}()

			select {
			case err := <-done:
				if err != nil {
					pt.logger.Info("Process PID %d exited with error: %v", pid, err)
				} else {
					pt.logger.Info("Process PID %d shut down gracefully", pid)
				}
				return nil
			case <-time.After(timeout):
				pt.logger.Warn("Graceful shutdown timeout for PID %d, forcing termination", pid)
				// Fall through to forced termination
			}
		}
	}

	// Forced termination
	pt.logger.Info("Force terminating process PID %d", pid)
	if err := process.Kill(); err != nil {
		return fmt.Errorf("failed to kill process: %v", err)
	}

	// Wait for process to actually die
	go func() {
		process.Wait()
	}()

	pt.logger.Info("Process PID %d terminated", pid)
	return nil
}

// isProcessRunning checks if a process is still running
func (pt *ProcessTracker) isProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// On Unix systems, sending signal 0 tests if process exists
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// CleanupDeadProcesses removes dead processes from tracking
func (pt *ProcessTracker) CleanupDeadProcesses() {
	pt.mutex.Lock()
	defer pt.mutex.Unlock()

	for arrowName, processes := range pt.processes {
		aliveProcesses := make([]*ProcessInfo, 0)
		
		for _, processInfo := range processes {
			if pt.isProcessRunning(processInfo.PID) {
				aliveProcesses = append(aliveProcesses, processInfo)
			} else {
				pt.logger.Debug("Removing dead process PID %d from tracking", processInfo.PID)
			}
		}

		if len(aliveProcesses) == 0 {
			delete(pt.processes, arrowName)
		} else {
			pt.processes[arrowName] = aliveProcesses
		}
	}
}

// StartCleanupRoutine starts a background routine to periodically clean up dead processes
func (pt *ProcessTracker) StartCleanupRoutine(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	pt.logger.Info("Starting process cleanup routine (interval: %v)", interval)

	for {
		select {
		case <-ctx.Done():
			pt.logger.Info("Process cleanup routine stopped")
			return
		case <-ticker.C:
			pt.CleanupDeadProcesses()
		}
	}
}

// HasRunningProcesses returns true if the arrow has any running processes
func (pt *ProcessTracker) HasRunningProcesses(arrowName string) bool {
	processes := pt.GetProcesses(arrowName)
	for _, processInfo := range processes {
		if pt.isProcessRunning(processInfo.PID) {
			return true
		}
	}
	return false
} 