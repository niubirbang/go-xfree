package goxfree

import (
	"os"
	"os/signal"
	"syscall"
	"testing"

	goxfree "github.com/niubirbang/go-xfree"
)

func buildNodes() (goxfree.Nodes, error) {
	var nodes goxfree.Nodes
	body, err := os.ReadFile("./tmp/test_nodes.txt")
	if err != nil {
		return nodes, err
	}
	return goxfree.NewNodesBase64(string(body)), nil
}

func TestCore(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("test failed: %v", r)
		}
	}()

	core := goxfree.NewCore(goxfree.NewOption(
		"./tmp",
		goxfree.WithLogLevel(goxfree.LevelDebug),
	))
	if err := core.Run(); err != nil {
		panic(err)
	}
	defer core.Quit()

	store, err := core.GetStore()
	if err != nil {
		t.Error("Get store failed:", err)
	} else {
		t.Log("Get store success:", store)
	}

	core.ListenStore(func(store goxfree.CoreStore) {
		t.Log("Listen store:", store)
	})

	nodes, err := buildNodes()
	if err != nil {
		t.Error("Build nodes failed:", err)
	}
	if err := core.ChangeNodes(nodes); err != nil {
		t.Error("Change nodes failed:", err)
	} else {
		t.Log("Change nodes success")
	}
	// if err := core.ChangeNetMode(goxfree.MODE_TUN); err != nil {
	// 	t.Error("Change net mode failed:", err)
	// } else {
	// 	t.Log("Change net mode success")
	// }
	if err := core.Open(); err != nil {
		t.Error("Open failed:", err)
	} else {
		t.Log("Open success")
	}

	termSign := make(chan os.Signal, 1)
	signal.Notify(termSign, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-termSign
}
