//go:build windows
// +build windows

package goxfree

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/Microsoft/go-winio"
	"github.com/gorilla/websocket"
)

func (c *client) quit() {
	pgid := c.cmdPgid
	if pgid == 0 {
		return
	}
	proc, err := os.FindProcess(pgid)
	if err != nil {
		return
	}
	if err := proc.Kill(); err != nil {
		return
	}
}

func newHttpUnixClient(address string) *api {
	dialer := func(network, addr string) (net.Conn, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return winio.DialPipeContext(ctx, address)
	}
	return &api{
		client: &http.Client{
			Transport: &http.Transport{
				Dial: dialer,
			},
		},
	}
}

func newWsUnixDialer(address string) *ws {
	return &ws{
		dialer: &websocket.Dialer{
			NetDialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return winio.DialPipeContext(ctx, address)
			},
			HandshakeTimeout: 45 * time.Second,
		},
	}
}

func (c *client) getPermission() (bool, error) {
	return c.getNetFirewallRule("xfree")
}

func (c *client) setPermission() error {
	return c.setNetFirewallRule("xfree", c.cmdPath)
}

func (c *client) getNetFirewallRule(name string) (bool, error) {
	cmd := exec.Command(
		"powershell",
		"-Command",
		fmt.Sprintf(
			`Get-NetFirewallRule -DisplayName "%s" | Select-Object -Property DisplayName, Direction, Action, Enabled, Profile | ConvertTo-Json`,
			name,
		),
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	var rule struct {
		DisplayName string `json:"DisplayName"`
		Direction   int    `json:"Direction"`
		Action      int    `json:"Action"`
		Enabled     bool   `json:"Enabled"`
		Profile     int    `json:"Profile"`
	}
	if err := json.Unmarshal(output, &rule); err != nil {
		return false, err
	}
	directionOK := rule.Direction == 1
	actionOK := rule.Action == 2
	enabledOK := rule.Enabled
	profileOK := rule.Profile == 0
	return directionOK && actionOK && enabledOK && profileOK, nil
}

func (c *client) setNetFirewallRule(name, filePath string) error {
	removeFirewallCMD := exec.Command(
		"powershell",
		"-Command",
		fmt.Sprintf(
			`Remove-NetFirewallRule -DisplayName "%s" -ErrorAction SilentlyContinue`,
			name,
		),
	)
	removeFirewallOutput, err := removeFirewallCMD.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove firewall rule: %v\noutput: %s", err, string(removeFirewallOutput))
	}

	setFirewallCMD := exec.Command(
		"powershell",
		"-Command",
		fmt.Sprintf(
			`New-NetFirewallRule -DisplayName "%s" -Direction Inbound -Action Allow -Program "%s" -Enabled True -Profile Any`,
			name,
			filePath,
		),
	)
	setFirewallOutput, err := setFirewallCMD.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set firewall rule: %v\noutput: %s", err, string(setFirewallOutput))
	}
	return nil
}
