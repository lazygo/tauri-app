package framework

import (
	"github.com/lazygo/client/pkg/sqlite"
)

type Model struct {
	manager  *sqlite.Manager
	table    string
	dbname   string
	noprefix bool
}

func (m *Model) SetManager(manager *sqlite.Manager) {
	m.manager = manager
}

func (m *Model) SetTable(table string) {
	m.table = table
}

func (m *Model) SetDb(dbname string) {
	m.dbname = dbname
}

func (m *Model) SetNoPrefix(n bool) {
	m.noprefix = n
}

func (m *Model) GetDb() *sqlite.DB {
	return m.db(m.dbname)
}

func (m *Model) db(dbName string) *sqlite.DB {
	if m.manager == nil {
		panic("no db manager")
	}
	database, err := m.manager.Database(dbName)
	if err != nil {
		panic(err)
	}
	return database
}

func (m *Model) QueryBuilder() sqlite.Builder {
	if m.manager == nil {
		panic("no db manager")
	}
	table := m.table
	if table == "" {
		// 没有指定表名
		panic("no table name")
	}
	database, err := m.manager.Database(m.dbname)
	if err != nil {
		panic("no database")
	}
	if m.noprefix {
		return database.TableRaw(table)
	}
	return database.Table(table)
}
