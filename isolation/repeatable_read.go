package isolation

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type RepeatableRead struct {
	Migrate func() error
	DB      *sql.DB
}

func (r *RepeatableRead) Run() (err error) {
	if r.Migrate == nil {
		return fmt.Errorf("migrate function is required")
	}
	if r.DB == nil {
		return fmt.Errorf("database connection is required")
	}

	// migrate
	log.Println("[o] migrating database...")
	err = r.Migrate()
	if err != nil {
		return
	}
	log.Print("[v] migrate success...\n\n")

	// create first transaction
	log.Print("1. we are going to create our FIRST database transaction, and set the isolation level to REPEATABLE READ")
	tx1, err := r.DB.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})
	if err != nil {
		return
	}
	log.Print("-> FIRST TRANSACTION CREATED\n\n")

	log.Print("2. then, let's see the product with ID 1 from the FIRST database transaction")
	var product Product
	err = tx1.QueryRow("SELECT id, name, quantity FROM products WHERE id = 1").Scan(&product.ID, &product.Name, &product.Quantity)
	if err != nil {
		var e error
		e = tx1.Rollback()

		if e != nil {
			return e
		}

		return
	}
	log.Println("->")
	fmt.Printf("(1)----------\nID: %d\nName: %s\nQuantity: %d\n-------------\n\n", product.ID, product.Name, product.Quantity)

	log.Println("3. update the quantity of product with ID 1 by 5 from the SECOND connection")
	_, err = r.DB.Exec("UPDATE products SET quantity = 5 WHERE id = 1")
	if err != nil {
		var e error
		e = tx1.Rollback()

		if e != nil {
			return e
		}

		return
	}
	log.Print("-> ROW UPDATED\n\n")

	log.Println("4. then, update the quantity of product with ID 1 by 10 from FIRST database transaction")
	_, err = tx1.Exec("UPDATE products SET quantity = 10 WHERE id = 1")
	if err != nil {
		var e error
		e = tx1.Rollback()

		if e != nil {
			return e
		}

		log.Println("->")
		fmt.Print("error during committing FIRST database transaction\n", err)
		fmt.Print("\n\n")
	}
	// log.Print("-> ROW UPDATED\n\n")

	// log.Println("5. commit the FIRST database transaction")
	// err = tx1.Commit()
	// if err != nil {
	// 	log.Print("-> error during committing FIRST database transaction", err)

	// 	var e error
	// 	e = tx1.Rollback()
	// 	if e != nil {
	// 		return e
	// 	}
	// }

	str := strings.Builder{}
	str.WriteString("-> As you can see, we are attempting to update the same row from the FIRST database transaction.")
	str.WriteString(" However, since that row has already been updated by the SECOND connection, an error occured")
	log.Print(str.String())

	return nil
}
