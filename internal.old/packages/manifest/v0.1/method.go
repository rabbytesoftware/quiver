package v0_1

import (
	"fmt"
)

// Methods represents the OS -> ARCH -> METHOD -> [commands] structure
// e.g., methods: { windows: { amd64: { install: [...], execute: [...] } } }
type Methods map[string]map[string]map[string][]string

func (m Methods) GetInstall() map[string]map[string][]string {
	return m["install"]
}

func (m Methods) GetExecute() map[string]map[string][]string {
	return m["execute"]
}

func (m Methods) GetUninstall() map[string]map[string][]string {
	return m["uninstall"]
}

func (m Methods) GetUpdate() map[string]map[string][]string {
	return m["update"]
}

func (m Methods) GetValidate() map[string]map[string][]string {
	return m["validate"]
}

func (m Methods) GetMethod(methodName string) map[string]map[string][]string {
	return m[methodName]
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

	m.ensureMethodExists(methodName)
	m.ensureOSExists(methodName, osName)
	(*m)[methodName][osName][archName] = commands

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

// ensureMethodExists initializes the method map if it doesn't exist
func (m *Methods) ensureMethodExists(methodName string) {
	if (*m)[methodName] == nil {
		(*m)[methodName] = make(map[string]map[string][]string)
	}
}

// ensureOSExists initializes the OS map for a method if it doesn't exist
func (m *Methods) ensureOSExists(methodName, osName string) {
	if (*m)[methodName][osName] == nil {
		(*m)[methodName][osName] = make(map[string][]string)
	}
}
