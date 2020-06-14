package database

import (
	"database/sql"
	"fmt"

	"github.com/squadcast_assignment/internal/config"
)

type DBClient interface {
	Execute(statement string, kind string) (int64, error)
	Query(statement string) (Row, error)
	Close() error
}

type dbClient struct {
	instance *sql.DB
}

func InitDatabaseConnection(dbConfig config.Database) (*dbClient, error) {
	var err error

	db := dbClient{}
	dbInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)
	db.instance, err = sql.Open(dbConfig.Dialect, dbInfo)

	if err != nil {
		return nil, fmt.Errorf("connection to MySQL failed: %s", err.Error())
	}

	if err = db.instance.Ping(); err != nil {
		return nil, fmt.Errorf("connection to MySQL failed: %s", err.Error())
	}

	return &db, nil
}

func (db *dbClient) Execute(statement string, kind string) (int64, error) {
	result, err := db.instance.Exec(statement)
	if err != nil {
		return -1, err
	}
	if kind == "UPDATE" {
		return result.RowsAffected()
	} else if kind == "CREATE" {
		return result.LastInsertId()
	} else if kind == "DELETE" {
		return result.RowsAffected()
	}
	return -1, fmt.Errorf("unknown operation: %s", kind)
}

func (db *dbClient) Query(statement string) (Row, error) {
	rows, err := db.instance.Query(statement)
	if err != nil {
		return new(MySQLRow), err
	}
	row := new(MySQLRow)
	row.Rows = rows
	return row, nil
}

func (db *dbClient) Close() error {
	if err := db.instance.Close(); err != nil {
		return err
	}
	return nil
}

type Row interface {
	Scan(dest ...interface{})
	Next() bool
}

type MySQLRow struct {
	Rows *sql.Rows
}

func (r MySQLRow) Scan(dest ...interface{}) {
	r.Rows.Scan(dest...)
}

func (r MySQLRow) Next() bool {
	return r.Rows.Next()
}
