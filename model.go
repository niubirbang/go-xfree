package goxfree

import "time"

var (
	cmdNames = map[string]string{
		"windows-amd64": "xfree-windows-amd64.exe",
		"windows-arm64": "xfree-windows-arm64.exe",
		"darwin-amd64":  "xfree-darwin-amd64",
		"darwin-arm64":  "xfree-darwin-arm64",
		"linux-amd64":   "xfree-linux-amd64",
		"linux-arm64":   "xfree-linux-arm64",
	}
)

const (
	LevelFatal LogLevel = "FATAL"
	LevelError LogLevel = "ERROR"
	LevelWarn  LogLevel = "WARN"
	LevelInfo  LogLevel = "INFO"
	LevelDebug LogLevel = "DEBUG"
	LevelTrace LogLevel = "TRACE"

	DefaultLogLevel = LevelWarn

	MODE_CORE    Mode = "CORE"
	MODE_MANAGER Mode = "MANAGER"

	MODE_SYSPROXY NetMode = "SYSPROXY"
	MODE_TUN      NetMode = "TUN"

	MODE_ABROAD    ProxyMode = "ABROAD"
	MODE_RETURNING ProxyMode = "RETURNING"
	MODE_GLOBAL    ProxyMode = "GLOBAL"

	STATUS_OPENED Status = "OPENED"
	STATUS_CLOSED Status = "CLOSED"

	AUTO_NODE_NAME                 = "****AUTO****"
	CURRENT_MODE_NONE  CurrentMode = "NONE"
	CURRENT_MODE_AUTO  CurrentMode = "AUTO"
	CURRENT_MODE_FIXED CurrentMode = "FIXED"

	NODE_MODEL_YAML   NodeModel = "YAML"
	NODE_MODEL_URI    NodeModel = "URI"
	NODE_MODEL_BASE64 NodeModel = "BASE64"

	MODEL_NODE  SubModel = "NODE"
	MODEL_GROUP SubModel = "GROUP"
	MODEL_AUTO  SubModel = "AUTO"
)

type (
	Mode        string
	NetMode     string
	ProxyMode   string
	Status      string
	CurrentMode string
	NodeModel   string
	LogLevel    string

	SubModel string
	Chain    []string
	Subs     []Sub
	Sub      struct {
		Name      string     `json:"name"`
		Chain     Chain      `json:"chain"`
		Model     SubModel   `json:"model"`
		ExpiredAt *time.Time `json:"expiredAt"`
		Usable    *bool      `json:"usable"`
		URI       string     `json:"uri"`
		Delay     int        `json:"delay"`
		Children  Subs       `json:"children"`
		NodeName  string     `json:"nodeName"`
	}

	CoreStore struct {
		Running   bool      `json:"running"`
		NetMode   NetMode   `json:"netMode"`
		ProxyMode ProxyMode `json:"proxyMode"`
		Current   string    `json:"current"`
		Status    Status    `json:"status"`
	}
	ManagerStore struct {
		CoreStore
		Subs         Subs        `json:"subs"`
		ConnectTime  *time.Time  `json:"connectTime"`
		CurrentSub   *Sub        `json:"currentSub"`
		CurrentChain Chain       `json:"currentChain"`
		CurrentMode  CurrentMode `json:"currentMode"`
	}

	Nodes struct {
		Model NodeModel   `json:"model"`
		Value interface{} `json:"value"`
	}

	Memery struct {
		Inuse   int `json:"inuse"`
		Oslimit int `json:"oslimit"`
	}
	Traffic struct {
		Up   int `json:"up"`
		Down int `json:"down"`
	}
	Connection struct {
		ID          string      `json:"id"`
		Upload      int         `json:"upload"`
		Download    int         `json:"download"`
		Start       time.Time   `json:"start"`
		Rule        string      `json:"rule"`
		RulePayload string      `json:"rulePayload"`
		Metadata    interface{} `json:"metadata"`
	}
	Connections struct {
		DownloadTotal int          `json:"downloadTotal"`
		UploadTotal   int          `json:"uploadTotal"`
		Memory        int          `json:"memory"`
		Connections   []Connection `json:"connections"`
	}
)

func NewNodesYaml(value []interface{}) Nodes {
	return Nodes{
		Model: NODE_MODEL_YAML,
		Value: value,
	}
}
func NewNodesUri(value []string) Nodes {
	return Nodes{
		Model: NODE_MODEL_URI,
		Value: value,
	}
}
func NewNodesBase64(value string) Nodes {
	return Nodes{
		Model: NODE_MODEL_BASE64,
		Value: value,
	}
}

func (s Sub) CanUse() bool {
	if s.Usable != nil && !*s.Usable {
		return false
	}
	if s.ExpiredAt != nil && s.ExpiredAt.Before(time.Now()) {
		return false
	}
	return true
}
