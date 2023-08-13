package config

import (
	_ "embed"
	"fmt"
	"path/filepath"

	"github.com/lazygo/client/pkg/appdata"
	"github.com/lazygo/lazygo/config"
	"github.com/lazygo/lazygo/logger"
)

type Logger struct {
	DefaultName string          `json:"default" toml:"default"`
	Adapter     []logger.Config `json:"adapter" toml:"adapter"`
}

type App struct {
	Name  string `json:"name" toml:"name"`
	Debug bool   `json:"debug" toml:"debug"`
}

var AppConfig App

//go:embed config.toml
var configData []byte

func Init() error {

	var base *config.Config
	var err error
	for _, loader := range []config.Loader{config.Json, config.Toml} {
		base, err = loader(configData)
		if err == nil {
			break
		}
	}
	if err != nil || base == nil {
		return fmt.Errorf("config file format fail")
	}

	// load server config
	err = base.Load("app", func(conf App) error {
		AppConfig = conf
		appdata.Init(conf.Name)
		return nil
	})
	if err != nil {
		return fmt.Errorf("load app config fail: %w", err)
	}
	userProfilePath, err := appdata.Path()
	if err != nil {
		return fmt.Errorf("get user profile path fail: %w", err)
	}

	// load logger config
	err = base.Load("logger", func(conf Logger) error {
		for index := range conf.Adapter {
			if filename, ok := conf.Adapter[index].Option["filename"]; ok {
				conf.Adapter[index].Option["filename"] = filepath.Join(userProfilePath, filename)
			}
		}
		return logger.Init(conf.Adapter, conf.DefaultName)
	})
	if err != nil {
		return fmt.Errorf("load logger config fail: %w", err)
	}

	return nil
}
