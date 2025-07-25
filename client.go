package goxfree

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

type (
	client struct {
		mu sync.Mutex

		option Option

		mode    Mode
		cmdPath string
		cmd     *exec.Cmd
		cmdPgid int
	}
	checker struct {
		mu        sync.Mutex
		buffer    bytes.Buffer
		checked   bool
		errChan   chan error
		timeoutDo func()
	}
)

func newClientCore(option Option) *client {
	c := &client{
		option: option,
		mode:   MODE_CORE,
	}
	c.init()
	return c
}

func newClientManager(option Option) *client {
	c := &client{
		option: option,
		mode:   MODE_MANAGER,
	}
	c.init()
	return c
}

func (c *client) checkPermission() error {
	permission, err := c.getPermission()
	if err != nil {
		return err
	}
	if !permission {
		if err := c.setPermission(); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) init() {
	c.cmdPath = path.Join(
		c.option.GetDir(),
		cmdNames[fmt.Sprintf(
			"%s-%s",
			c.option.GetPlatform(), c.option.GetArch(),
		)],
	)
}

func (c *client) Run() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cmd != nil {
		return errors.New("runed")
	}

	if err := c.checkPermission(); err != nil {
		return err
	}

	// build args
	var args []string
	switch c.mode {
	case MODE_CORE:
		args = []string{
			"crun",
			"--dir", c.option.GetDir(),
			"--level", string(c.option.GetLogLevel()),
			"--mixed", strconv.Itoa(c.option.GetMixedPort()),
			"--ext", strconv.Itoa(c.option.GetExternalControllerPort()),
			"--net", string(c.option.GetNetMode()),
			"--proxy", string(c.option.GetProxyMode()),
			"--unix", c.option.GetServerUnixAddress(),
			"--tcp", c.option.GetServerTcpAddress(),
			"--close-sysproxy", strconv.FormatBool(c.option.GetDoCloseSysproxy()),
			"--delay-url", c.option.GetTestDelayURL(),
			"--delay-timeout", c.option.GetTestDelayTimeout().String(),
		}
	case MODE_MANAGER:
		args = []string{
			"mrun",
			"--dir", c.option.GetDir(),
			"--level", string(c.option.GetLogLevel()),
			"--mixed", strconv.Itoa(c.option.GetMixedPort()),
			"--ext", strconv.Itoa(c.option.GetExternalControllerPort()),
			"--net", string(c.option.GetNetMode()),
			"--proxy", string(c.option.GetProxyMode()),
			"--unix", c.option.GetServerUnixAddress(),
			"--tcp", c.option.GetServerTcpAddress(),
			"--close-sysproxy", strconv.FormatBool(c.option.GetDoCloseSysproxy()),
			"--delay-url", c.option.GetTestDelayURL(),
			"--delay-timeout", c.option.GetTestDelayTimeout().String(),
			"--auto", strconv.FormatBool(c.option.GetNeedAuto()),
			"--delay", strconv.FormatBool(c.option.GetNeedMinDelay()),
		}
	}

	// new checker
	checker := newChecker(func() {
		c.quit()
	})
	checkerListener := checker.getListener()

	// cmd
	// log.Println("client command", c.cmdPath, args)
	c.cmd = exec.Command(c.cmdPath, args...)
	c.cmd.Stdout = checker
	c.cmd.Stderr = checker

	if err := c.cmd.Start(); err != nil {
		return err
	}
	c.cmdPgid = c.cmd.Process.Pid

	return checkerListener()
}

func (c *client) Quit() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.quit()
	return nil
}

func newChecker(timeoutDo func()) *checker {
	return &checker{
		errChan:   make(chan error, 1),
		timeoutDo: timeoutDo,
	}
}

func (c *checker) Write(p []byte) (n int, err error) {
	// log.Println(string(p))

	c.mu.Lock()
	defer c.mu.Unlock()
	n, err = c.buffer.Write(p)
	if err != nil {
		return
	}
	for {
		line, err := c.buffer.ReadString('\n')
		if err != nil {
			c.buffer.WriteString(line)
			break
		}
		trimmed := strings.TrimSpace(line)
		if !c.checked {
			switch {
			case strings.HasPrefix(trimmed, "CMD:SUCCESS"):
				c.errChan <- nil
				c.checked = true
			case strings.HasPrefix(trimmed, "CMD:ERROR:"):
				c.errChan <- errors.New(strings.TrimPrefix(trimmed, "CMD:ERROR:"))
				c.checked = true
			}
		}
	}
	return
}

func (c *checker) getListener() func() error {
	return func() error {
		timeout := time.After(20 * time.Second)
		for {
			select {
			case err := <-c.errChan:
				return err
			case <-timeout:
				if c.timeoutDo != nil {
					c.timeoutDo()
				}
				return errors.New("timeout")
			}
		}
	}
}
