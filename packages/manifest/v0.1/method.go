package v0_1

type Methods map[string]map[string][]string

func (m Methods) GetInstall() map[string][]string {
	if install, exists := m["install"]; exists {
		return install
	}
	return nil
}

func (m Methods) GetExecute() map[string][]string {
	if execute, exists := m["execute"]; exists {
		return execute
	}
	return nil
}

func (m Methods) GetUninstall() map[string][]string {
	if uninstall, exists := m["uninstall"]; exists {
		return uninstall
	}
	return nil
}

func (m Methods) GetUpdate() map[string][]string {
	if update, exists := m["update"]; exists {
		return update
	}
	return nil
}

func (m Methods) GetValidate() map[string][]string {
	if validate, exists := m["validate"]; exists {
		return validate
	}
	return nil
}

func (m Methods) GetMethod(methodName string) map[string][]string {
	if method, exists := m[methodName]; exists {
		return method
	}
	return nil
}
