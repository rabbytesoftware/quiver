// cmd/quiver_demo/main.go
package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"gopkg.in/yaml.v3"
)

type Severity string

const (
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

type ValidationResult struct {
	ID          string    `json:"id" yaml:"id"`
	Title       string    `json:"title" yaml:"title"`
	Severity    Severity  `json:"severity" yaml:"severity"`
	Blocking    bool      `json:"blocking" yaml:"blocking"`
	Timestamp   time.Time `json:"timestamp" yaml:"timestamp"`
	Remediation string    `json:"remediation,omitempty" yaml:"remediation,omitempty"`
	Details     any       `json:"details,omitempty" yaml:"details,omitempty"`
}

type Report struct {
	Results []ValidationResult `json:"results" yaml:"results"`
}

// Manager reads rules and runs validators.
type Manager struct {
	Rules map[string]interface{}
}

func NewManager() *Manager {
	return &Manager{Rules: make(map[string]interface{})}
}

func (m *Manager) LoadRules(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &m.Rules)
}

func (m *Manager) ValidateAll() Report {
	var report Report

	report.Results = append(report.Results, m.validateOS()...)
	report.Results = append(report.Results, m.validateMemory()...)
	report.Results = append(report.Results, m.validateDisk()...)
	report.Results = append(report.Results, m.validateDependencies()...)
	report.Results = append(report.Results, m.validateNetwork()...)
	report.Results = append(report.Results, m.validateSecurity()...)

	return report
}

// ---------------- Validators ----------------

func (m *Manager) validateOS() []ValidationResult {
	var out []ValidationResult
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	allowedOS := m.getStringList("os.allowed")
	allowedArch := m.getStringList("os.arch")

	if len(allowedOS) > 0 && !contains(allowedOS, goos) {
		out = append(out, ValidationResult{
			ID:        "os-compat",
			Title:     fmt.Sprintf("Unsupported OS: %s", goos),
			Severity:  SeverityError,
			Blocking:  true,
			Timestamp: time.Now(),
			Remediation: fmt.Sprintf("Supported OS: %v. Use a supported OS or update configs/requirements.yaml",
				allowedOS),
		})
	} else {
		out = append(out, ValidationResult{
			ID:        "os-compat",
			Title:     fmt.Sprintf("OS supported: %s", goos),
			Severity:  SeverityInfo,
			Blocking:  false,
			Timestamp: time.Now(),
		})
	}

	if len(allowedArch) > 0 && !contains(allowedArch, goarch) {
		out = append(out, ValidationResult{
			ID:        "arch-compat",
			Title:     fmt.Sprintf("Unsupported architecture: %s", goarch),
			Severity:  SeverityError,
			Blocking:  true,
			Timestamp: time.Now(),
			Remediation: fmt.Sprintf("Supported architectures: %v. Use a supported CPU or update configs/requirements.yaml",
				allowedArch),
		})
	} else {
		out = append(out, ValidationResult{
			ID:        "arch-compat",
			Title:     fmt.Sprintf("Architecture supported: %s", goarch),
			Severity:  SeverityInfo,
			Blocking:  false,
			Timestamp: time.Now(),
		})
	}

	return out
}

func (m *Manager) validateMemory() []ValidationResult {
	var out []ValidationResult
	res, ok := m.Rules["resources"].(map[string]interface{})
	if !ok {
		out = append(out, ValidationResult{
			ID:        "memory-no-config",
			Title:     "No memory requirement configured; skipping memory check",
			Severity:  SeverityWarning,
			Blocking:  false,
			Timestamp: time.Now(),
		})
		return out
	}
	memReq := toInt(res["memory_mb"])
	if memReq == 0 {
		out = append(out, ValidationResult{
			ID:        "memory-no-config",
			Title:     "No memory requirement set; skipping memory check",
			Severity:  SeverityWarning,
			Blocking:  false,
			Timestamp: time.Now(),
		})
		return out
	}
	vm, err := mem.VirtualMemory()
	if err != nil {
		out = append(out, ValidationResult{
			ID:        "memory-read-failed",
			Title:     "Failed to read memory stats",
			Severity:  SeverityWarning,
			Blocking:  false,
			Timestamp: time.Now(),
			Remediation: "gopsutil failed to read memory stats. Try running with appropriate permissions.",
			Details:    err.Error(),
		})
		return out
	}
	totalMB := int(vm.Total / 1024 / 1024)
	if totalMB < memReq {
		out = append(out, ValidationResult{
			ID:        "memory-check",
			Title:     fmt.Sprintf("Not enough memory: %d MB (required %d MB)", totalMB, memReq),
			Severity:  SeverityError,
			Blocking:  true,
			Timestamp: time.Now(),
			Remediation: "Increase RAM or run the project on a machine with larger memory.",
			Details: map[string]int{"total_mb": totalMB, "required_mb": memReq},
		})
	} else {
		out = append(out, ValidationResult{
			ID:        "memory-check",
			Title:     fmt.Sprintf("Memory OK: %d MB", totalMB),
			Severity:  SeverityInfo,
			Blocking:  false,
			Timestamp: time.Now(),
			Details:   map[string]int{"total_mb": totalMB, "required_mb": memReq},
		})
	}
	return out
}

func (m *Manager) validateDisk() []ValidationResult {
	var out []ValidationResult
	res, ok := m.Rules["resources"].(map[string]interface{})
	if !ok {
		out = append(out, ValidationResult{
			ID:        "disk-no-config",
			Title:     "No disk requirement configured; skipping disk check",
			Severity:  SeverityWarning,
			Blocking:  false,
			Timestamp: time.Now(),
		})
		return out
	}
	diskReqGB := toInt(res["disk_gb"])
	if diskReqGB == 0 {
		out = append(out, ValidationResult{
			ID:        "disk-no-config",
			Title:     "No disk requirement set; skipping disk check",
			Severity:  SeverityWarning,
			Blocking:  false,
			Timestamp: time.Now(),
		})
		return out
	}

	path := rootPathForOS()
	usage, err := disk.Usage(path)
	if err != nil {
		out = append(out, ValidationResult{
			ID:        "disk-read-failed",
			Title:     "Failed to read disk usage",
			Severity:  SeverityWarning,
			Blocking:  false,
			Timestamp: time.Now(),
			Remediation: "gopsutil failed to read disk usage. Try running with appropriate permissions or check mounting.",
			Details:    err.Error(),
		})
		return out
	}

	freeGB := int(usage.Free / 1024 / 1024 / 1024)
	if freeGB < diskReqGB {
		out = append(out, ValidationResult{
			ID:        "disk-check",
			Title:     fmt.Sprintf("Not enough free disk: %d GB available (required %d GB)", freeGB, diskReqGB),
			Severity:  SeverityError,
			Blocking:  true,
			Timestamp: time.Now(),
			Remediation: fmt.Sprintf("Free up %d GB or increase disk size on %s", (diskReqGB - freeGB), path),
			Details: map[string]int{"free_gb": freeGB, "required_gb": diskReqGB},
		})
	} else {
		out = append(out, ValidationResult{
			ID:        "disk-check",
			Title:     fmt.Sprintf("Disk OK: %d GB free on %s", freeGB, path),
			Severity:  SeverityInfo,
			Blocking:  false,
			Timestamp: time.Now(),
			Details:   map[string]int{"free_gb": freeGB, "required_gb": diskReqGB},
		})
	}
	return out
}

func (m *Manager) validateDependencies() []ValidationResult {
	var out []ValidationResult
	deps := m.getStringList("dependencies.required_bins")
	if len(deps) == 0 {
		out = append(out, ValidationResult{
			ID:        "deps-no-config",
			Title:     "No required dependencies configured; skipping dependency checks",
			Severity:  SeverityWarning,
			Blocking:  false,
			Timestamp: time.Now(),
		})
		return out
	}

	for _, bin := range deps {
		_, err := exec.LookPath(bin)
		if err != nil {
			out = append(out, ValidationResult{
				ID:        "dep-missing-" + bin,
				Title:     fmt.Sprintf("Dependency not found: %s", bin),
				Severity:  SeverityError,
				Blocking:  true,
				Timestamp: time.Now(),
				Remediation: fmt.Sprintf("Install %s and ensure it is on PATH (or update configs/requirements.yaml).", bin),
			})
		} else {
			out = append(out, ValidationResult{
				ID:        "dep-ok-" + bin,
				Title:     fmt.Sprintf("Dependency available: %s", bin),
				Severity:  SeverityInfo,
				Blocking:  false,
				Timestamp: time.Now(),
			})
		}
	}
	return out
}

// ---------------- New Validators ----------------

func (m *Manager) validateNetwork() []ValidationResult {
	var out []ValidationResult
	hostsMap, ok := m.Rules["network"].(map[string]interface{})
	if !ok {
		return out
	}
	requiredHosts := []string{}
	if arr, ok := hostsMap["required_hosts"].([]interface{}); ok {
		for _, h := range arr {
			if s, ok := h.(string); ok {
				requiredHosts = append(requiredHosts, s)
			}
		}
	}
	for _, host := range requiredHosts {
		_, err := net.LookupHost(host)
		if err != nil {
			out = append(out, ValidationResult{
				ID:        "network-" + host,
				Title:     fmt.Sprintf("Cannot resolve host: %s", host),
				Severity:  SeverityError,
				Blocking:  true,
				Timestamp: time.Now(),
				Remediation: fmt.Sprintf("Check network connection or DNS for %s", host),
			})
		} else {
			out = append(out, ValidationResult{
				ID:        "network-" + host,
				Title:     fmt.Sprintf("Host reachable: %s", host),
				Severity:  SeverityInfo,
				Blocking:  false,
				Timestamp: time.Now(),
			})
		}
	}
	return out
}

func (m *Manager) validateSecurity() []ValidationResult {
	var out []ValidationResult
	secMap, ok := m.Rules["security"].(map[string]interface{})
	if !ok {
		return out
	}

	requireSudo := false
	if val, ok := secMap["require_sudo"].(bool); ok {
		requireSudo = val
	}

	if requireSudo {
		u, err := user.Current()
		if err != nil {
			out = append(out, ValidationResult{
				ID:        "security-sudo",
				Title:     "Cannot determine current user",
				Severity:  SeverityWarning,
				Blocking:  false,
				Timestamp: time.Now(),
				Details:   err.Error(),
			})
		} else if u.Uid != "0" && runtime.GOOS != "windows" {
			out = append(out, ValidationResult{
				ID:        "security-sudo",
				Title:     "Not running as root/admin",
				Severity:  SeverityWarning,
				Blocking:  false,
				Timestamp: time.Now(),
				Remediation: "Run as root/admin for full access",
			})
		} else {
			out = append(out, ValidationResult{
				ID:        "security-sudo",
				Title:     "Running with sufficient privileges",
				Severity:  SeverityInfo,
				Blocking:  false,
				Timestamp: time.Now(),
			})
		}
	}
	return out
}

// ------------- Helpers ----------------

func (m *Manager) getStringList(path string) []string {
	switch path {
	case "os.allowed":
		if osv, ok := m.Rules["os"]; ok {
			if osMap, ok := osv.(map[string]interface{}); ok {
				if arr, ok := osMap["allowed"].([]interface{}); ok {
					return toStringSlice(arr)
				}
			}
		}
	case "os.arch":
		if osv, ok := m.Rules["os"]; ok {
			if osMap, ok := osv.(map[string]interface{}); ok {
				if arr, ok := osMap["arch"].([]interface{}); ok {
					return toStringSlice(arr)
				}
			}
		}
	case "dependencies.required_bins":
		if dv, ok := m.Rules["dependencies"]; ok {
			if depMap, ok := dv.(map[string]interface{}); ok {
				if arr, ok := depMap["required_bins"].([]interface{}); ok {
					return toStringSlice(arr)
				}
			}
		}
	}
	return nil
}

func toStringSlice(arr []interface{}) []string {
	out := make([]string, 0, len(arr))
	for _, v := range arr {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func toInt(v interface{}) int {
	switch t := v.(type) {
	case int:
		return t
	case int64:
		return int(t)
	case float64:
		return int(t)
	case uint64:
		return int(t)
	default:
		return 0
	}
}

func contains(arr []string, s string) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}
	return false
}

func rootPathForOS() string {
	if runtime.GOOS == "windows" {
		return "C:\\"
	}
	return "/"
}

func main() {
	cfgPath := "configs/requirements.yaml"
	m := NewManager()
	if err := m.LoadRules(cfgPath); err != nil {
		fmt.Println("Error loading rules:", err)
		return
	}
	report := m.ValidateAll()

	// -------- Terminal Table --------
	fmt.Println("\n=== System Requirements Validation ===")
	fmt.Printf("%-20s %-50s %-10s\n", "Check ID", "Title", "Severity")
	fmt.Println(strings.Repeat("-", 85))
	for _, res := range report.Results {
		fmt.Printf("%-20s %-50s %-10s\n", res.ID, res.Title, res.Severity)
		if res.Remediation != "" {
			fmt.Printf("  -> Fix: %s\n", res.Remediation)
		}
	}
	fmt.Println(strings.Repeat("=", 85))

	// -------- JSON Report --------
	jsonFile := "validation_report.json"
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Println("Failed to generate JSON report:", err)
		return
	}
	err = os.WriteFile(jsonFile, data, 0644)
	if err != nil {
		fmt.Println("Failed to save JSON report:", err)
		return
	}
	fmt.Printf("\n[info] JSON report saved to %s\n", jsonFile)
}
