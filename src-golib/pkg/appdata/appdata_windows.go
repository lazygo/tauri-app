//go:build windows
// +build windows

package appdata

import (
	"io/fs"
	"os"
	"syscall"
)

var (
	UserHomeDir             string
	FirstAppDataPath        string
	FirstGlobalAppDataPath  string
	SecondAppDataPath       string
	SecondGlobalAppDataPath string
)

func initialize() error {
	var err error
	UserHomeDir, err = os.UserHomeDir()
	if err != nil {
		return err
	}
	FirstAppDataPath = os.Getenv("APPDATA")
	FirstGlobalAppDataPath = os.Getenv("PROGRAMDATA")
	SecondAppDataPath = os.Getenv("PROGRAMDATA")
	SecondGlobalAppDataPath = os.Getenv("PROGRAMDATA")

	return nil
}
