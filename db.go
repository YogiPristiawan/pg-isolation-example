package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type DB struct {
	Config Config
	conn   *sql.DB
}

func (d *DB) CreateConn() (*sql.DB, error) {
	var err error
	conn, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		d.Config.DB.Username,
		d.Config.DB.Password,
		d.Config.DB.Host,
		d.Config.DB.Port,
		d.Config.DB.Database,
	))
	if err != nil {
		return nil, err
	}
	log.Println("[v] database connected...")

	d.conn = conn

	return conn, nil
}

func (d *DB) Migrate() (err error) {
	if d.conn == nil {
		return errors.New("connection is empty")
	}

	tx, err := d.conn.Begin()
	if err != nil {
		return
	}

	// drop table
	var dropTableSql = `DROP TABLE IF EXISTS products`
	_, err = tx.Exec(dropTableSql)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			return e
		}

		return
	}

	// create table
	var createTableSql = `CREATE TABLE IF NOT EXISTS products (
		id BIGSERIAL NOT NULL PRIMARY KEY,
		name VARCHAR(63) NOT NULL,
		quantity INTEGER NOT NULL
	)`
	_, err = tx.Exec(createTableSql)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			return e
		}

		return
	}

	// create dummy data
	stmt, err := tx.Prepare(`INSERT INTO products (name, quantity) VALUES ($1, $2)`)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			return e
		}

		return
	}

	rows := []struct {
		name     string
		quantity int
	}{
		{
			name:     "Hat",
			quantity: 90,
		},
		{
			name:     "T-shirt",
			quantity: 100,
		},
		{
			name:     "Hoodie",
			quantity: 50,
		},
		{
			name:     "Pants",
			quantity: 30,
		},
		{
			name:     "Dress",
			quantity: 50,
		},
	}
	for _, row := range rows {
		_, err = tx.Stmt(stmt).Exec(row.name, row.quantity)
		if err != nil {
			if e := tx.Rollback(); e != nil {
				return e
			}

			return
		}
	}

	err = tx.Commit()
	if err != nil {
		if e := tx.Rollback(); e != nil {
			return e
		}

		return
	}

	return
}
