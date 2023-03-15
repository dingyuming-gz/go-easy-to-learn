package database

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

// Database 包装了数据库连接操作。
type Database struct {
	driver   string
	connStr  string
	db       *sql.DB
	connOpen bool
}

// NewDatabase 创建一个新的数据库实例。
func NewDatabase(driver string, connStr string) *Database {
	return &Database{
		driver:  driver,
		connStr: connStr,
	}
}

// Connect 连接到数据库。
func (d *Database) Connect() error {
	if d.connOpen {
		return nil
	}

	db, err := sql.Open(d.driver, d.connStr)
	if err != nil {
		println("sql.Open", err)
		return err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		println("sql.Open", err)
		return err
	}

	d.db = db
	d.connOpen = true

	return nil
}

// Close 关闭与数据库的连接。
func (d *Database) Close() error {
	if !d.connOpen {
		return nil
	}

	err := d.db.Close()
	if err != nil {
		return err
	}

	d.db = nil
	d.connOpen = false

	return nil
}

// Save 将数据保存到数据库中。
func (d *Database) Save(data []byte) error {
	if !d.connOpen {
		return ErrDatabaseNotConnected
	}

	_, err := d.db.Exec("INSERT INTO data (value) VALUES ($1)", data)
	if err != nil {
		return err
	}

	return nil
}

var (
	ErrDatabaseNotConnected = errors.New("database not connected")
)
