// +build !android,!ios,!mobile,!nacl

package app

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

func TestDefaultTheme(t *testing.T) {
	assert.Equal(t, theme.DarkTheme(), defaultTheme())
}

func TestEnsureDir(t *testing.T) {
	tmpDir := filepath.Join(rootConfigDir(), "fynetest")

	ensureDirExists(tmpDir)
	if st, err := os.Stat(tmpDir); err != nil || !st.IsDir() {
		t.Error("Could not ensure directory exists")
	}

	os.Remove(tmpDir)
}

func TestWatchSettings(t *testing.T) {
	settings := &settings{}
	listener := make(chan fyne.Settings)
	settings.AddChangeListener(listener)

	settings.fileChanged() // simulate the settings file changing

	select {
	case _ = <-listener:
	case <-time.After(100 * time.Millisecond):
		t.Error("Settings listener was not called")
	}
}

func TestWatchFile(t *testing.T) {
	path := filepath.Join(rootConfigDir(), "fyne-temp-watch.txt")
	os.Create(path)
	defer os.Remove(path)

	called := make(chan interface{})
	watchFile(path, func() {
		called <- true
	})
	file, _ := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	file.WriteString(" ")
	file.Close()

	select {
	case _ = <-called:
	case <-time.After(100 * time.Millisecond):
		t.Error("File watcher callback was not called")
	}
}

func TestFileWatcher_FileDeleted(t *testing.T) {
	path := filepath.Join(rootConfigDir(), "fyne-temp-watch.txt")
	os.Create(path)
	defer os.Remove(path)

	called := make(chan interface{})
	watcher := watchFile(path, func() {
		called <- true
	})
	if watcher == nil {
		assert.Fail(t, "Could not start watcher")
		return
	}

	defer watcher.Close()
	os.Remove(path)
	os.Create(path)

	select {
	case _ = <-called:
	case <-time.After(100 * time.Millisecond):
		t.Error("File watcher callback was not called")
	}
}
