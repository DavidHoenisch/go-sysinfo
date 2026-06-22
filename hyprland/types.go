package hyprland

type Version struct {
	Branch              string   `json:"branch"`
	Commit              string   `json:"commit"`
	Version             string   `json:"version"`
	Dirty               bool     `json:"dirty"`
	CommitMessage       string   `json:"commit_message"`
	CommitDate          string   `json:"commit_date"`
	Tag                 string   `json:"tag"`
	Commits             string   `json:"commits"`
	BuildAquamarine     string   `json:"buildAquamarine"`
	BuildHyprlang       string   `json:"buildHyprlang"`
	BuildHyprutils      string   `json:"buildHyprutils"`
	BuildHyprcursor     string   `json:"buildHyprcursor"`
	BuildHyprgraphics   string   `json:"buildHyprgraphics"`
	SystemAquamarine    string   `json:"systemAquamarine"`
	SystemHyprlang      string   `json:"systemHyprlang"`
	SystemHyprutils     string   `json:"systemHyprutils"`
	SystemHyprcursor    string   `json:"systemHyprcursor"`
	SystemHyprgraphics  string   `json:"systemHyprgraphics"`
	ABIHash             string   `json:"abiHash"`
	Flags               []string `json:"flags"`
}

type WorkspaceRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Monitor struct {
	ID              int          `json:"id"`
	Name            string       `json:"name"`
	Description     string       `json:"description"`
	Make            string       `json:"make"`
	Model           string       `json:"model"`
	Serial          string       `json:"serial"`
	Width           int          `json:"width"`
	Height          int          `json:"height"`
	PhysicalWidth   int          `json:"physicalWidth"`
	PhysicalHeight  int          `json:"physicalHeight"`
	RefreshRate     float64      `json:"refreshRate"`
	X               int          `json:"x"`
	Y               int          `json:"y"`
	ActiveWorkspace WorkspaceRef `json:"activeWorkspace"`
	Scale           float64      `json:"scale"`
	Transform       int          `json:"transform"`
	Focused         bool         `json:"focused"`
	DPMSStatus      bool         `json:"dpmsStatus"`
	VRR             bool         `json:"vrr"`
	Disabled        bool         `json:"disabled"`
	AvailableModes  []string     `json:"availableModes"`
}

type Workspace struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Monitor         string `json:"monitor"`
	MonitorID       int    `json:"monitorID"`
	Windows         int    `json:"windows"`
	HasFullscreen   bool   `json:"hasfullscreen"`
	LastWindow      string `json:"lastwindow"`
	LastWindowTitle string `json:"lastwindowtitle"`
	IsPersistent    bool   `json:"ispersistent"`
	TiledLayout     string `json:"tiledLayout"`
}

type Client struct {
	Address        string       `json:"address"`
	Mapped         bool         `json:"mapped"`
	Hidden         bool         `json:"hidden"`
	Visible        bool         `json:"visible"`
	AcceptsInput   bool         `json:"acceptsInput"`
	At             []int        `json:"at"`
	Size           []int        `json:"size"`
	Workspace      WorkspaceRef `json:"workspace"`
	Floating       bool         `json:"floating"`
	Monitor        int          `json:"monitor"`
	Class          string       `json:"class"`
	Title          string       `json:"title"`
	InitialClass   string       `json:"initialClass"`
	InitialTitle   string       `json:"initialTitle"`
	PID            int          `json:"pid"`
	XWayland       bool         `json:"xwayland"`
	Pinned         bool         `json:"pinned"`
	Fullscreen     int          `json:"fullscreen"`
	Tags           []string     `json:"tags"`
	FocusHistoryID int          `json:"focusHistoryID"`
	StableID       string       `json:"stableId"`
}

type Bind struct {
	Locked          bool   `json:"locked"`
	Mouse           bool   `json:"mouse"`
	Release         bool   `json:"release"`
	Repeat          bool   `json:"repeat"`
	LongPress       bool   `json:"longPress"`
	ModMask         int    `json:"modmask"`
	Submap          string `json:"submap"`
	Key             string `json:"key"`
	Keycode         int    `json:"keycode"`
	Description     string `json:"description"`
	Dispatcher      string `json:"dispatcher"`
	Arg             string `json:"arg"`
	HasDescription  bool   `json:"has_description"`
}

type Option struct {
	Option string  `json:"option"`
	Custom string  `json:"custom"`
	Int    int     `json:"int"`
	Float  float64 `json:"float"`
	Str    string  `json:"str"`
	Set    bool    `json:"set"`
}

type SessionInfo struct {
	Version              string
	ActiveWorkspaceID    int
	ActiveWorkspaceName  string
	FocusedMonitor       string
	ConfigErrorCount     int
}
