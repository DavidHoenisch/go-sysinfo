package hyprland

func GetBinds(r Reader) []Bind {
	out, ok := hyprctlJSON(r, "binds")
	if !ok {
		return nil
	}
	var binds []Bind
	if !decodeJSON(out, &binds) {
		return nil
	}
	return binds
}

func (r Reader) GetBinds() []Bind {
	return GetBinds(r)
}
