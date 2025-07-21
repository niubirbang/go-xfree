package goxfree

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"time"
)

type Manager struct {
	client *client
	api    *api
	ws     *ws
}

func NewManager(option Option) *Manager {
	client := newClientManager(option)
	return &Manager{
		client: client,
		api:    newHttpUnixClient(client.option.GetServerUnixAddress()),
		ws:     newWsUnixDialer(client.option.GetServerUnixAddress()),
	}
}

func (m *Manager) Run() error {
	if err := m.client.Run(); err != nil {
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
			if err = m.TestClient(); err == nil {
				return nil
			}
		}
	}
}

func (m *Manager) quit() error {
	_, err := m.api.put("/quit", nil, nil)
	if err != nil {
		return err
	}
	return err
}
func (m *Manager) Quit() error {
	if err := m.quit(); err != nil {
		log.Println("use api quit failed:", err)
	}
	if err := m.client.Quit(); err != nil {
		return err
	}
	return nil
}

func (m *Manager) TestClient() error {
	_, err := m.api.get("/test", nil)
	if err != nil {
		return err
	}
	return err
}
func (m *Manager) GetStatus() (Status, error) {
	var data Status
	body, err := m.api.get("/status", nil)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(body, &data)
	return data, err
}
func (m *Manager) GetNetMode() (NetMode, error) {
	var data NetMode
	body, err := m.api.get("/net-mode", nil)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(body, &data)
	return data, err
}
func (m *Manager) GetProxyMode() (ProxyMode, error) {
	var data ProxyMode
	body, err := m.api.get("/proxy-mode", nil)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(body, &data)
	return data, err
}
func (m *Manager) GetDelay(name string) (int, error) {
	var data int
	query := make(url.Values)
	query.Set("name", name)
	body, err := m.api.get("/delay", query)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(body, &data)
	return data, err
}
func (m *Manager) GetAllDelay() (map[string]int, error) {
	var data map[string]int
	body, err := m.api.get("/all-delay", nil)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(body, &data)
	return data, err
}
func (m *Manager) GetStore() (ManagerStore, error) {
	var data ManagerStore
	body, err := m.api.get("/store", nil)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(body, &data)
	return data, err
}

func (m *Manager) Open() error {
	_, err := m.api.put("/open", nil, nil)
	return err
}
func (m *Manager) Close() error {
	_, err := m.api.put("/close", nil, nil)
	return err
}
func (m *Manager) ChangeNetMode(mode NetMode) error {
	_, err := m.api.put("/change-net-mode", nil, mode)
	return err
}
func (m *Manager) ChangeProxyMode(mode ProxyMode) error {
	_, err := m.api.put("/change-proxy-mode", nil, mode)
	return err
}
func (m *Manager) ChangeSubs(subs Subs) error {
	_, err := m.api.put("/change-subs", nil, subs)
	return err
}
func (m *Manager) ChangeNodeAuto() error {
	_, err := m.api.put("/change-node-auto", nil, nil)
	return err
}
func (m *Manager) ChangeNodeFixed(chain Chain) error {
	_, err := m.api.put("/change-node-fixed", nil, chain)
	return err
}
func (m *Manager) TestDelay(name string) (int, error) {
	var data int
	body, err := m.api.put("/test-delay", nil, name)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(body, &data)
	return data, err
}
func (m *Manager) TestAllDelay(name string) (map[string]int, error) {
	var data map[string]int
	body, err := m.api.put("/test-all-delay", nil, nil)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(body, &data)
	return data, err
}

func (m *Manager) ListenMemery(fn func(Memery)) {
	if fn == nil {
		return
	}
	go func() {
		for {
			conn, err := m.ws.conn("/listen-memery")
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
func (m *Manager) ListenTraffic(fn func(Traffic)) {
	if fn == nil {
		return
	}
	go func() {
		for {
			conn, err := m.ws.conn("/listen-traffic")
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
func (m *Manager) ListenConnections(fn func(Connections)) {
	if fn == nil {
		return
	}
	go func() {
		for {
			conn, err := m.ws.conn("/listen-connections")
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
func (m *Manager) ListenDelay(fn func(map[string]int)) {
	if fn == nil {
		return
	}
	go func() {
		for {
			conn, err := m.ws.conn("/listen-delay")
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
func (m *Manager) ListenStore(fn func(ManagerStore)) {
	if fn == nil {
		return
	}
	go func() {
		for {
			conn, err := m.ws.conn("/listen-store")
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
				var data ManagerStore
				if err := json.Unmarshal(msg, &data); err == nil {
					fn(data)
				}
			}
		}
	}()
}
