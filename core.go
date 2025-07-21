package goxfree

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"
)

type Core struct {
	client *client
	api    *api
	ws     *ws
}

func NewCore(option Option) *Core {
	client := newClientCore(option)
	return &Core{
		client: client,
		api:    newHttpUnixClient(client.option.GetServerUnixAddress()),
		ws:     newWsUnixDialer(client.option.GetServerUnixAddress()),
	}
}

func (c *Core) Run() error {
	if err := c.client.Run(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var err error
	for {
		select {
		case <-ctx.Done():
			return errors.New("timeout")
		case <-ticker.C:
			if err = c.TestClient(); err == nil {
				return nil
			}
		}
	}
}

func (c *Core) quit() error {
	_, err := c.api.put("/quit", nil, nil)
	if err != nil {
		return err
	}
	return err
}
func (c *Core) Quit() error {
	if err := c.quit(); err != nil {
		log.Println("use api quit failed:", err)
	}
	if err := c.client.Quit(); err != nil {
		return err
	}
	return nil
}

func (c *Core) TestClient() error {
	_, err := c.api.get("/test", nil)
	if err != nil {
		return err
	}
	return err
}
func (c *Core) GetStatus() (Status, error) {
	var data Status
	body, err := c.api.get("/status", nil)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(body, &data)
	return data, err
}
func (c *Core) GetNetMode() (NetMode, error) {
	var data NetMode
	body, err := c.api.get("/net-mode", nil)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(body, &data)
	return data, err
}
func (c *Core) GetProxyMode() (ProxyMode, error) {
	var data ProxyMode
	body, err := c.api.get("/proxy-mode", nil)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(body, &data)
	return data, err
}
func (c *Core) GetDelay(name string) (int, error) {
	var data int
	query := make(url.Values)
	query.Set("name", name)
	body, err := c.api.get("/delay", query)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(body, &data)
	return data, err
}
func (c *Core) GetAllDelay() (map[string]int, error) {
	var data map[string]int
	body, err := c.api.get("/all-delay", nil)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(body, &data)
	return data, err
}
func (c *Core) GetStore() (CoreStore, error) {
	var data CoreStore
	body, err := c.api.get("/store", nil)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(body, &data)
	return data, err
}

func (c *Core) Open() error {
	_, err := c.api.put("/open", nil, nil)
	return err
}
func (c *Core) Close() error {
	_, err := c.api.put("/close", nil, nil)
	return err
}
func (c *Core) ChangeNetMode(mode NetMode) error {
	_, err := c.api.put("/change-net-mode", nil, mode)
	return err
}
func (c *Core) ChangeProxyMode(mode ProxyMode) error {
	_, err := c.api.put("/change-proxy-mode", nil, mode)
	return err
}
func (c *Core) ChangeNodes(nodes Nodes) error {
	param := map[string]interface{}{
		"model": nodes.Model,
	}
	switch nodes.Model {
	case NODE_MODEL_YAML:
		param["yaml"] = nodes.Value
	case NODE_MODEL_URI:
		param["uri"] = nodes.Value
	case NODE_MODEL_BASE64:
		param["base64"] = nodes.Value
	default:
		return fmt.Errorf("unknow model: %s", nodes.Model)
	}
	_, err := c.api.put("/change-nodes", nil, param)
	return err
}
func (c *Core) ChangeNodeAuto() error {
	_, err := c.api.put("/change-node-auto", nil, nil)
	return err
}
func (c *Core) ChangeNodeFixed(name string) error {
	_, err := c.api.put("/change-node-fixed", nil, name)
	return err
}
func (c *Core) TestDelay(name string) (int, error) {
	var data int
	body, err := c.api.put("/test-delay", nil, name)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(body, &data)
	return data, err
}
func (c *Core) TestAllDelay(name string) (map[string]int, error) {
	var data map[string]int
	body, err := c.api.put("/test-all-delay", nil, nil)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(body, &data)
	return data, err
}

func (c *Core) ListenMemery(fn func(Memery)) {
	if fn == nil {
		return
	}
	go func() {
		for {
			conn, err := c.ws.conn("/listen-memery")
			if err != nil {
				time.Sleep(time.Second * 2)
				continue
			}
			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					conn.Close()
					break
				}
				var data Memery
				if err := json.Unmarshal(msg, &data); err == nil {
					fn(data)
				}
			}
		}
	}()
}
func (c *Core) ListenTraffic(fn func(Traffic)) {
	if fn == nil {
		return
	}
	go func() {
		for {
			conn, err := c.ws.conn("/listen-traffic")
			if err != nil {
				time.Sleep(time.Second * 2)
				continue
			}
			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					conn.Close()
					break
				}
				var data Traffic
				if err := json.Unmarshal(msg, &data); err == nil {
					fn(data)
				}
			}
		}
	}()
}
func (c *Core) ListenConnections(fn func(Connections)) {
	if fn == nil {
		return
	}
	go func() {
		for {
			conn, err := c.ws.conn("/listen-connections")
			if err != nil {
				time.Sleep(time.Second * 2)
				continue
			}
			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					conn.Close()
					break
				}
				var data Connections
				if err := json.Unmarshal(msg, &data); err == nil {
					fn(data)
				}
			}
		}
	}()
}
func (c *Core) ListenDelay(fn func(map[string]int)) {
	if fn == nil {
		return
	}
	go func() {
		for {
			conn, err := c.ws.conn("/listen-delay")
			if err != nil {
				time.Sleep(time.Second * 2)
				continue
			}
			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					conn.Close()
					break
				}
				var data map[string]int
				if err := json.Unmarshal(msg, &data); err == nil {
					fn(data)
				}
			}
		}
	}()
}
func (c *Core) ListenStore(fn func(CoreStore)) {
	if fn == nil {
		return
	}
	go func() {
		for {
			conn, err := c.ws.conn("/listen-store")
			if err != nil {
				time.Sleep(time.Second * 2)
				continue
			}
			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					conn.Close()
					break
				}
				var data CoreStore
				if err := json.Unmarshal(msg, &data); err == nil {
					fn(data)
				}
			}
		}
	}()
}
