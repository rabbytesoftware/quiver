package v0_2

import (
	"fmt"
)

// Methods represents the OS -> ARCH -> METHOD -> [commands] structure for v0.2
// e.g., methods: { windows: { amd64: { install: [...], execute: [...] } } }
type Methods map[string]map[string]map[string][]string

func (m Methods) GetInstall() map[string]map[string][]string {
	result := make(map[string]map[string][]string)
	for osName, osData := range m {
		for archName, archData := range osData {
			if installCommands, exists := archData["install"]; exists {
				if result[osName] == nil {
					result[osName] = make(map[string][]string)
				}
				result[osName][archName] = installCommands
			}
		}
	}
	return result
}

func (m Methods) GetExecute() map[string]map[string][]string {
	result := make(map[string]map[string][]string)
	for osName, osData := range m {
		for archName, archData := range osData {
			if executeCommands, exists := archData["execute"]; exists {
				if result[osName] == nil {
					result[osName] = make(map[string][]string)
				}
				result[osName][archName] = executeCommands
			}
		}
	}
	return result
}

func (m Methods) GetUninstall() map[string]map[string][]string {
	result := make(map[string]map[string][]string)
	for osName, osData := range m {
		for archName, archData := range osData {
			if uninstallCommands, exists := archData["uninstall"]; exists {
				if result[osName] == nil {
					result[osName] = make(map[string][]string)
				}
				result[osName][archName] = uninstallCommands
			}
		}
	}
	return result
}

func (m Methods) GetUpdate() map[string]map[string][]string {
	result := make(map[string]map[string][]string)
	for osName, osData := range m {
		for archName, archData := range osData {
			if updateCommands, exists := archData["update"]; exists {
				if result[osName] == nil {
					result[osName] = make(map[string][]string)
				}
				result[osName][archName] = updateCommands
			}
		}
	}
	return result
}

func (m Methods) GetValidate() map[string]map[string][]string {
	result := make(map[string]map[string][]string)
	for osName, osData := range m {
		for archName, archData := range osData {
			if validateCommands, exists := archData["validate"]; exists {
				if result[osName] == nil {
					result[osName] = make(map[string][]string)
				}
				result[osName][archName] = validateCommands
			}
		}
	}
	return result
}

func (m Methods) GetMethod(methodName string) map[string]map[string][]string {
	result := make(map[string]map[string][]string)
	for osName, osData := range m {
		for archName, archData := range osData {
			if methodCommands, exists := archData[methodName]; exists {
				if result[osName] == nil {
					result[osName] = make(map[string][]string)
				}
				result[osName][archName] = methodCommands
			}
		}
	}
	return result
}

func (m *Methods) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var rawMethods map[string]interface{}
	if err := unmarshal(&rawMethods); err != nil {
		return fmt.Errorf("failed to unmarshal methods: %w", err)
	}

	*m = make(Methods)

	for osName, osData := range rawMethods {
		if err := m.processOS(osName, osData); err != nil {
			return err
		}
	}

	return nil
}

// processOS handles the OS level of the methods structure
func (m *Methods) processOS(osName string, osData interface{}) error {
	osMap, ok := osData.(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf("invalid OS data for '%s': expected map, got %T", osName, osData)
	}

	(*m)[osName] = make(map[string]map[string][]string)

	for archKey, archData := range osMap {
		archName := fmt.Sprintf("%v", archKey)
		if err := m.processArch(osName, archName, archData); err != nil {
			return err
		}
	}

	return nil
}

// processArch handles the architecture level of the methods structure
func (m *Methods) processArch(osName, archName string, archData interface{}) error {
	archMap, ok := archData.(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf("invalid arch data for '%s/%s': expected map, got %T", osName, archName, archData)
	}

	(*m)[osName][archName] = make(map[string][]string)

	for methodKey, methodData := range archMap {
		methodName := fmt.Sprintf("%v", methodKey)
		if err := m.processMethod(osName, archName, methodName, methodData); err != nil {
			return err
		}
	}

	return nil
}

// processMethod handles the method level of the methods structure
func (m *Methods) processMethod(osName, archName, methodName string, methodData interface{}) error {
	commands, err := parseCommands(methodData)
	if err != nil {
		return fmt.Errorf("invalid commands for '%s/%s/%s': %w", osName, archName, methodName, err)
	}

	(*m)[osName][archName][methodName] = commands

	return nil
}

// parseCommands converts interface{} to []string for commands
func parseCommands(data interface{}) ([]string, error) {
	commandList, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected command list, got %T", data)
	}

	commands := make([]string, len(commandList))
	for i, cmd := range commandList {
		commands[i] = fmt.Sprintf("%v", cmd)
	}

	return commands, nil
} 