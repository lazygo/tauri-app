package appdata

import (
	"os"
	"path/filepath"

	"github.com/lazygo/client/utils"
)

var config Config

type Config struct {
	Appname string `json:"appname" toml:"appname"`
}

func Init(appname string) error {
	config.Appname = appname
	return initialize()
}

func Path() (string, error) {

	firstPath := filepath.Join(FirstAppDataPath, config.Appname)
	firstPath, err := filepath.Abs(firstPath)
	if err != nil {
		return "", err
	}
	if utils.PathExists(firstPath) && utils.IsDir(firstPath) {
		return firstPath, nil
	}

	secondPath := filepath.Join(SecondAppDataPath, config.Appname)
	secondPath, err = filepath.Abs(secondPath)
	if err != nil {
		return "", err
	}
	if utils.PathExists(secondPath) && utils.IsDir(secondPath) {
		return secondPath, nil
	}

	err = os.MkdirAll(firstPath, os.ModePerm)
	if err != nil {
		return "", err
	}
	return firstPath, nil
}
