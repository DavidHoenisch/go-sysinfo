package hyprland

import "strings"

func GetConfigErrors(r Reader) []string {
	out, ok := hyprctlJSON(r, "configerrors")
	if !ok {
		return nil
	}
	var errors []string
	if !decodeJSON(out, &errors) {
		return nil
	}
	filtered := make([]string, 0, len(errors))
	for _, errMsg := range errors {
		errMsg = strings.TrimSpace(errMsg)
		if errMsg != "" {
			filtered = append(filtered, errMsg)
		}
	}
	return filtered
}

func (r Reader) GetConfigErrors() []string {
	return GetConfigErrors(r)
}
