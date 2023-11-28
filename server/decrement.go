package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
)

func Decrement(w http.ResponseWriter, r *http.Request) {
	var selectSql = `SELECT quantity FROM product WHERE id = 1`

	var qty int64
	tx, err := pool.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		log.Fatal(err)
	}
	err = tx.QueryRow(selectSql).Scan(&qty)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("got quantty %d", qty)

	var sql = `UPDATE product SET quantity = $1 WHERE id = 1`

	// tx, err := pool.Begin()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	res, err := tx.Exec(sql, qty-1)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			log.Fatal(err)
		}
		log.Fatal(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%d row updated", affected)

	err = tx.Commit()
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			log.Fatal(err2)
		}
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("oke"))
}
