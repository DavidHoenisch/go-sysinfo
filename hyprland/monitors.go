package hyprland

func GetMonitors(r Reader) []Monitor {
	out, ok := hyprctlJSON(r, "monitors")
	if !ok {
		return nil
	}
	var monitors []Monitor
	if !decodeJSON(out, &monitors) {
		return nil
	}
	return monitors
}

func (r Reader) GetMonitors() []Monitor {
	return GetMonitors(r)
}
