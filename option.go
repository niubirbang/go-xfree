package goxfree

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

var (
	defaultNetMode                = MODE_SYSPROXY
	defaultProxyMode              = MODE_ABROAD
	defaultTestDelayURL           = "https://www.gstatic.com/generate_204"
	defaultTestDelayTimeout       = 5 * time.Second
	defaultMixedPort              = 12401
	defaultExternalControllerPort = 12402
	defaultDoCloseSysproxy        = true
	defaultNeedAuto               = true
	defaultNeedMinDelay           = true
	defaultServerUnixAddress      string
)

func init() {

}

type setter func(*Option)
type Option struct {
	platform *string
	arch     *string
	dir      string
	logLevel *LogLevel

	mixedPort              *int
	externalControllerPort *int
	netMode                *NetMode
	proxyMode              *ProxyMode
	serverUnixAddress      *string
	serverTcpAddress       *string
	doCloseSysproxy        *bool
	testDelayURL           *string
	testDelayTimeout       *time.Duration
	needAuto               *bool
	needMinDelay           *bool
}

func NewOption(dir string, options ...setter) Option {
	var o Option
	if !filepath.IsAbs(dir) {
		currentDir, _ := os.Getwd()
		dir = path.Join(currentDir, dir)
	}
	o.dir = dir
	for _, option := range options {
		option(&o)
	}
	return o
}

// core: ok
// manager: ok
func WithPlatform(platform string) setter {
	return func(o *Option) {
		o.platform = &platform
	}
}

// core: ok
// manager: ok
func WithArch(arch string) setter {
	return func(o *Option) {
		o.arch = &arch
	}
}

// core: ok
// manager: ok
func WithLogLevel(l LogLevel) setter {
	return func(o *Option) {
		o.logLevel = &l
	}
}

// core: ok
// manager: ok
func WithNetMode(mode NetMode) setter {
	return func(o *Option) {
		o.netMode = &mode
	}
}

// core: ok
// manager: ok
func WithProxyMode(mode ProxyMode) setter {
	return func(o *Option) {
		o.proxyMode = &mode
	}
}

// core: ok
// manager: ok
func WithMixedPort(port int) setter {
	return func(o *Option) {
		o.mixedPort = &port
	}
}

// core: ok
// manager: ok
func WithExternalControllerPort(port int) setter {
	return func(o *Option) {
		o.externalControllerPort = &port
	}
}

// core: ok
// manager: ok
func WithDoCloseSysproxy(close bool) setter {
	return func(o *Option) {
		o.doCloseSysproxy = &close
	}
}

// core: ok
// manager: ok
func WithTestDelayURL(url string) setter {
	return func(o *Option) {
		o.testDelayURL = &url
	}
}

// core: ok
// manager: ok
func WithTestDelayTimeout(d time.Duration) setter {
	return func(o *Option) {
		o.testDelayTimeout = &d
	}
}

// core: ok
// manager: ok
func WithServerUnixAddress(address string) setter {
	return func(o *Option) {
		o.serverUnixAddress = &address
	}
}

// core: ok
// manager: ok
func WithServerTcpAddress(address string) setter {
	return func(o *Option) {
		o.serverTcpAddress = &address
	}
}

// core: invalid
// manager: ok
func WithNeedAuto(need bool) setter {
	return func(o *Option) {
		o.needAuto = &need
	}
}

// core: invalid
// manager: ok
func WithNeedMinDelay(need bool) setter {
	return func(o *Option) {
		o.needMinDelay = &need
	}
}

func (o Option) GetPlatform() string {
	if o.platform != nil {
		return *o.platform
	}
	return runtime.GOOS
}
func (o Option) GetArch() string {
	if o.arch != nil {
		return *o.arch
	}
	return runtime.GOARCH
}
func (o Option) GetDir() string {
	return o.dir
}
func (o Option) GetLogLevel() LogLevel {
	if o.logLevel != nil {
		return *o.logLevel
	}
	return DefaultLogLevel
}
func (o Option) GetNetMode() NetMode {
	if o.netMode != nil {
		return *o.netMode
	}
	return defaultNetMode
}
func (o Option) GetProxyMode() ProxyMode {
	if o.proxyMode != nil {
		return *o.proxyMode
	}
	return defaultProxyMode
}
func (o Option) GetMixedPort() int {
	if o.mixedPort != nil {
		return *o.mixedPort
	}
	return defaultMixedPort
}
func (o Option) GetExternalControllerPort() int {
	if o.externalControllerPort != nil {
		return *o.externalControllerPort
	}
	return defaultExternalControllerPort
}
func (o Option) GetDoCloseSysproxy() bool {
	if o.doCloseSysproxy != nil {
		return *o.doCloseSysproxy
	}
	return defaultDoCloseSysproxy
}
func (o Option) GetTestDelayURL() string {
	if o.testDelayURL != nil {
		return *o.testDelayURL
	}
	return defaultTestDelayURL
}
func (o Option) GetTestDelayTimeout() time.Duration {
	if o.testDelayTimeout != nil {
		return *o.testDelayTimeout
	}
	return defaultTestDelayTimeout
}
func (o Option) GetNeedAuto() bool {
	if o.needAuto != nil {
		return *o.needAuto
	}
	return defaultNeedAuto
}
func (o Option) GetNeedMinDelay() bool {
	if o.needMinDelay != nil {
		return *o.needMinDelay
	}
	return defaultNeedMinDelay
}
func (o Option) GetServerUnixAddress() string {
	if o.serverUnixAddress != nil {
		return *o.serverUnixAddress
	}
	switch runtime.GOOS {
	case "windows":
		return `\\.\pipe\xfree`
	case "darwin":
		return `/tmp/xfree.sock`
	case "linux":
		return `/tmp/xfree.sock`
	default:
		return ""
	}
}
func (o Option) GetServerTcpAddress() string {
	if o.serverTcpAddress != nil {
		return *o.serverTcpAddress
	}
	return ""
}
