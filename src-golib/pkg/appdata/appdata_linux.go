//go:build linux
// +build linux

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
	FirstAppDataPath = UserHomeDir
	FirstGlobalAppDataPath = "/etc/"
	SecondAppDataPath = UserHomeDir
	SecondGlobalAppDataPath = "/etc/"

	return nil
}
