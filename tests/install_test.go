package goxfree

import (
	"testing"

	goxfree "github.com/niubirbang/go-xfree"
)

func TestInstall(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("test failed: %v", r)
		}
	}()

	installer := goxfree.NewInstaller(goxfree.NewOption(
		"./tmp",
		goxfree.WithLogLevel(goxfree.LevelDebug),
	))
	if err := installer.Run(); err != nil {
		panic(err)
	}
	defer installer.Quit()
}
