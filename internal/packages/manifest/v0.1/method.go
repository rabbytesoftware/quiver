package v0_1

import (
	"fmt"
)

// Methods now represents platform -> method -> commands structure
// e.g., methods: { windows: { install: [...], execute: [...] }, linux: { install: [...] } }
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
	var temp map[string]interface{}
	if err := unmarshal(&temp); err != nil {
		return err
	}
	*m = make(Methods)
	for method, data := range temp {
		switch val := data.(type) {
		case map[interface{}]interface{}:
			osMap := make(map[string]map[string][]string)
			for osKey, inner := range val {
				osStr := fmt.Sprintf("%v", osKey)
				switch innerVal := inner.(type) {
				case []interface{}:
					cmds := make([]string, len(innerVal))
					for i, c := range innerVal {
						cmds[i] = fmt.Sprintf("%v", c)
					}
					osMap[osStr] = map[string][]string{"any": cmds}
				case map[interface{}]interface{}:
					archMap := make(map[string][]string)
					for archKey, cmdsVal := range innerVal {
						archStr := fmt.Sprintf("%v", archKey)
						switch cmds := cmdsVal.(type) {
						case []interface{}:
							cmdList := make([]string, len(cmds))
							for i, c := range cmds {
								cmdList[i] = fmt.Sprintf("%v", c)
							}
							archMap[archStr] = cmdList
						default:
							return fmt.Errorf("unexpected commands type for %s/%s/%s", method, osStr, archStr)
						}
					}
					osMap[osStr] = archMap
				default:
					return fmt.Errorf("unexpected inner type for %s/%s", method, osStr)
				}
			}
			(*m)[method] = osMap
		default:
			return fmt.Errorf("unexpected data type for method %s", method)
		}
	}
	return nil
}
