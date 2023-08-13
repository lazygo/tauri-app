//go:build darwin
// +build darwin

package appdata

import (
	"os"
	"path/filepath"
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
	FirstAppDataPath = filepath.Join(UserHomeDir, "/Library/Preferences/")
	FirstGlobalAppDataPath = "/Library/Preferences/"
	SecondAppDataPath = filepath.Join(UserHomeDir, "/Library/Preferences/")
	SecondGlobalAppDataPath = "/Library/Preferences/"

	return nil
}
