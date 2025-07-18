//go:build linux
// +build linux

package goxfree

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

func (c *client) quit() {
	pgid := c.cmdPgid
	if pgid == 0 {
		return
	}
	_ = syscall.Kill(-pgid, syscall.SIGTERM)
	time.AfterFunc(5*time.Second, func() {
		_ = syscall.Kill(-pgid, syscall.SIGKILL)
	})
}

func newHttpUnixClient(address string) *api {
	dialer := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return net.Dial("unix", address)
	}
	return &api{
		client: &http.Client{
			Transport: &http.Transport{
				DialContext:        dialer,
				DisableCompression: false,
				ForceAttemptHTTP2:  false,
			},
			Timeout: 5 * time.Second,
		},
	}
}

func newWsUnixDialer(address string) *ws {
	return &ws{
		dialer: &websocket.Dialer{
			NetDialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", address)
			},
		},
	}
}

func (c *client) getPermission() (bool, error) {
	return c.getChownChmod([]string{c.cmdPath})
}

func (c *client) setPermission() error {
	return c.setChownChmod([]string{c.cmdPath})
}

func (c *client) getChownChmod(filePaths []string) (bool, error) {
	for _, filePath := range filePaths {
		fi, err := os.Stat(filePath)
		if err != nil {
			return false, err
		}
		stat, ok := fi.Sys().(*syscall.Stat_t)
		if !ok {
			return false, errors.New("unable to obtain file owner information")
		}
		if stat.Uid != 0 {
			return false, nil
		}
		mode := fi.Mode()
		if mode&os.ModeSetuid == 0 ||
			mode&0100 == 0 ||
			mode&0010 == 0 ||
			mode&0001 == 0 {
			return false, nil
		}
	}
	return true, nil
}

func (c *client) setChownChmod(filePaths []string) error {
	var innerCmds []string
	for _, filePath := range filePaths {
		innerCmds = append(innerCmds,
			fmt.Sprintf(`chown root:root "%s"`, filePath),
			fmt.Sprintf(`chmod 4755 "%s"`, filePath),
		)
	}
	shell := strings.Join(innerCmds, " && ")

	cmd := exec.Command("pkexec", "sh", "-c", shell)
	cmd.Env = []string{
		"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		"LANG=C.UTF-8",
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set chmod and chown: %v\noutput: %s", err, string(output))
	}
	return nil
}
