package hyprland

func GetWorkspaces(r Reader) []Workspace {
	out, ok := hyprctlJSON(r, "workspaces")
	if !ok {
		return nil
	}
	var workspaces []Workspace
	if !decodeJSON(out, &workspaces) {
		return nil
	}
	return workspaces
}

func (r Reader) GetWorkspaces() []Workspace {
	return GetWorkspaces(r)
}

func GetActiveWorkspace(r Reader) *Workspace {
	out, ok := hyprctlJSON(r, "activeworkspace")
	if !ok {
		return nil
	}
	var workspace Workspace
	if !decodeJSON(out, &workspace) {
		return nil
	}
	return &workspace
}

func (r Reader) GetActiveWorkspace() *Workspace {
	return GetActiveWorkspace(r)
}
