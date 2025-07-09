package v0_1

// Methods now represents platform -> method -> commands structure
// e.g., methods: { windows: { install: [...], execute: [...] }, linux: { install: [...] } }
type Methods map[string]map[string][]string

func (m Methods) GetInstall() map[string][]string {
	result := make(map[string][]string)
	for platform, methods := range m {
		if install, exists := methods["install"]; exists {
			result[platform] = install
		}
	}
	return result
}

func (m Methods) GetExecute() map[string][]string {
	result := make(map[string][]string)
	for platform, methods := range m {
		if execute, exists := methods["execute"]; exists {
			result[platform] = execute
		}
	}
	return result
}

func (m Methods) GetUninstall() map[string][]string {
	result := make(map[string][]string)
	for platform, methods := range m {
		if uninstall, exists := methods["uninstall"]; exists {
			result[platform] = uninstall
		}
	}
	return result
}

func (m Methods) GetUpdate() map[string][]string {
	result := make(map[string][]string)
	for platform, methods := range m {
		if update, exists := methods["update"]; exists {
			result[platform] = update
		}
	}
	return result
}

func (m Methods) GetValidate() map[string][]string {
	result := make(map[string][]string)
	for platform, methods := range m {
		if validate, exists := methods["validate"]; exists {
			result[platform] = validate
		}
	}
	return result
}

func (m Methods) GetMethod(methodName string) map[string][]string {
	result := make(map[string][]string)
	for platform, methods := range m {
		if method, exists := methods[methodName]; exists {
			result[platform] = method
		}
	}
	return result
}
