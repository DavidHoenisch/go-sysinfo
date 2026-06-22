package hyprland

func GetVersion(r Reader) *Version {
	out, ok := hyprctlJSON(r, "version")
	if !ok {
		return nil
	}
	var version Version
	if !decodeJSON(out, &version) {
		return nil
	}
	return &version
}

func (r Reader) GetVersion() *Version {
	return GetVersion(r)
}
