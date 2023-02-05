package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	User         string
	Password     string
	Host         string
	Name         string
	MaxIdleConns int
	MaxOpenConns int
	DbType       string
	//Not included: TLS disable/enable. out of scope
}

func Open(cfg Config) (*sqlx.DB, error) {
	//Mysql or the relevant driver is pretty odd. Normally we'd like to use a url.URL here, but that was straight up
	//Not working with some weird errors. So We'll use a good old conn string instead.
	connString := fmt.Sprintf("%s:%s@(%s)/%s?parseTime=true", cfg.User, cfg.Password, cfg.Host, cfg.Name)

	db, err := sqlx.Open("mysql", connString)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	return db, nil
}
