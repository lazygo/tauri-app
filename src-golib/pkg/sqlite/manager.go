package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"path"
	"sync"
	"time"
)

type Config struct {
	Name            string `json:"name"`
	DbName          string `json:"dbname"`
	Prefix          string `json:"prefix"`
	MaxOpenConns    int    `json:"max_open_conns"`
	MaxIdleConns    int    `json:"max_idle_conns"`
	ConnMaxLifetime int    `json:"conn_max_lifetime"`
}

type Manager struct {
	sync.Map
}

// init 初始化数据库连接
func (m *Manager) init(prefix string, conf []Config) error {
	for _, item := range conf {
		if _, ok := m.Load(item.Name); ok {
			// 已连接的就不再次连接了
			continue
		}
		db, err := m.open(prefix, item)
		if err != nil {
			return err
		}
		err = db.Ping()
		if err != nil {
			return err
		}
		m.Store(item.Name, newDb(item.Name, db, item.Prefix))
	}
	return nil
}

// CloseAll 关闭数据库连接
func (m *Manager) CloseAll() error {
	var err error
	m.Range(func(name, db interface{}) bool {
		_ = db.(*DB).Close()
		m.Delete(name)
		return true
	})
	return err
}

// open 连接数据库
func (m *Manager) open(prefix string, item Config) (*sql.DB, error) {
	dsn := path.Join(prefix, item.DbName)
	database, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	database.SetMaxIdleConns(item.MaxIdleConns)
	database.SetMaxOpenConns(item.MaxOpenConns)
	if item.ConnMaxLifetime > 0 {
		database.SetConnMaxLifetime(time.Duration(item.ConnMaxLifetime) * time.Second)
	}
	return database, nil
}

// Init 初始化数据库
func Init(prefix string, conf []Config) (*Manager, error) {
	var manager = &Manager{}
	err := manager.init(prefix, conf)
	return manager, err
}

// Database 通过名称获取数据库
func (m *Manager) Database(name string) (*DB, error) {
	if databases, ok := m.Load(name); ok {
		return databases.(*DB), nil
	}
	return nil, ErrDatabaseNotExists
}
