package hyprland

func GetClients(r Reader) []Client {
	out, ok := hyprctlJSON(r, "clients")
	if !ok {
		return nil
	}
	var clients []Client
	if !decodeJSON(out, &clients) {
		return nil
	}
	return clients
}

func (r Reader) GetClients() []Client {
	return GetClients(r)
}

func GetActiveWindow(r Reader) *Client {
	out, ok := hyprctlJSON(r, "activewindow")
	if !ok {
		return nil
	}
	var client Client
	if !decodeJSON(out, &client) {
		return nil
	}
	return &client
}

func (r Reader) GetActiveWindow() *Client {
	return GetActiveWindow(r)
}
