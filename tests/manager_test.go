package goxfree

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"testing"

	goxfree "github.com/niubirbang/go-xfree"
)

func buildSubs() (goxfree.Subs, error) {
	var nodes goxfree.Subs
	body, err := os.ReadFile("./tmp/test_subs.json")
	if err != nil {
		return nodes, err
	}
	err = json.Unmarshal(body, &nodes)
	return nodes, err
}

func TestManager(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("test failed: %v", r)
		}
	}()

	manager := goxfree.NewManager(goxfree.NewOption(
		"./tmp",
		goxfree.WithLogLevel(goxfree.LevelDebug),
	))
	if err := manager.Run(); err != nil {
		panic(err)
	}
	defer manager.Quit()

	store, err := manager.GetStore()
	if err != nil {
		t.Error("Get store failed:", err)
	} else {
		t.Log("Get store success:", store)
	}

	manager.ListenStore(func(store goxfree.ManagerStore) {
		t.Log("Listen store:", store)
	})

	subs, err := buildSubs()
	if err != nil {
		t.Error("Build subs failed:", err)
	}
	if err := manager.ChangeSubs(subs); err != nil {
		t.Error("Change subs failed:", err)
	} else {
		t.Log("Change subs success")
	}
	// if err := manager.ChangeNetMode(goxfree.MODE_TUN); err != nil {
	// 	t.Error("Change net mode failed:", err)
	// } else {
	// 	t.Log("Change net mode success")
	// }
	if err := manager.Open(); err != nil {
		t.Error("Open failed:", err)
	} else {
		t.Log("Open success")
	}

	termSign := make(chan os.Signal, 1)
	signal.Notify(termSign, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-termSign
}
