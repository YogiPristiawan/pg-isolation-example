package isolation

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type RadCommitted struct {
	Migrate func() error
	DB      *sql.DB
}

func (r *RadCommitted) Run() (err error) {
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
	log.Print("1. we are going to create our FIRST database transaction, and set the isolation level to READ COMMITTED")
	tx1, err := r.DB.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return
	}
	log.Print("-> FIRST TRANSACTION CREATED\n\n")

	// create second transaction
	log.Print("2. then, we are going to create our SECOND database transaction")
	tx2, err := r.DB.Begin()
	if err != nil {
		if e := tx1.Rollback(); e != nil {
			return e
		}

		return
	}
	log.Print("-> SECOND TRANSACTION CREATED\n\n")

	log.Println("3. read the product with ID 1 from the FIRST database transaction")
	var product1 Product
	err = tx1.QueryRow("SELECT id, name, quantity FROM products WHERE id = 1").Scan(&product1.ID, &product1.Name, &product1.Quantity)
	if err != nil {
		var e error
		var e2 error
		e = tx1.Rollback()
		e2 = tx2.Rollback()

		if e != nil {
			return e
		}
		if e2 != nil {
			return e2
		}

		return
	}
	log.Println("->")
	fmt.Printf("(1)----------\nID: %d\nName: %s\nQuantity: %d\n-------------\n\n", product1.ID, product1.Name, product1.Quantity)

	log.Println("4. then, update the product quantity with ID 1 by 10 from the SECOND database transaction")
	_, err = tx2.Exec("UPDATE products SET quantity = 10 WHERE id = 1")
	if err != nil {
		var e error
		var e2 error
		e = tx1.Rollback()
		e2 = tx2.Rollback()

		if e != nil {
			return e
		}
		if e2 != nil {
			return e2
		}

		return
	}
	log.Print("-> ROW UPDATED\n\n")

	log.Println("5. read the product with ID 1 again from the FIRST database transaction")
	var product2 Product
	err = tx1.QueryRow("SELECT id, name, quantity FROM products WHERE id = 1").Scan(&product2.ID, &product2.Name, &product2.Quantity)
	if err != nil {
		var e error
		var e2 error
		e = tx1.Rollback()
		e2 = tx2.Rollback()

		if e != nil {
			return e
		}
		if e2 != nil {
			return e2
		}

		return
	}
	fmt.Printf("(1)----------\nID: %d\nName: %s\nQuantity: %d\n-------------\n", product2.ID, product2.Name, product2.Quantity)
	log.Print("-> as you can see, the product with ID 1 has the same value as in the step 3. this is because the SECOND database transaction has not been committed yet.\n\n")

	log.Println("6. let's commit the SECOND database transaction")
	err = tx2.Commit()
	if err != nil {
		var e error
		var e2 error
		e = tx1.Rollback()
		e2 = tx2.Rollback()

		if e != nil {
			return e
		}
		if e2 != nil {
			return e2
		}

		return
	}
	log.Print("-> SECOND TRANSACTION COMMITTED\n\n")

	log.Println("7. let's read the product with ID 1 again from the FIRST database transaction")
	var product3 Product
	err = tx1.QueryRow("SELECT id, name, quantity FROM products WHERE id = 1").Scan(&product3.ID, &product3.Name, &product3.Quantity)
	if err != nil {
		var e error
		e = tx1.Rollback()

		if e != nil {
			return e
		}

		return
	}
	err = tx1.Commit()
	if err != nil {
		var e error
		e = tx1.Rollback()

		if e != nil {
			return e
		}

		return
	}
	fmt.Printf("(1)----------\nID: %d\nName: %s\nQuantity: %d\n-------------\n", product3.ID, product3.Name, product3.Quantity)
	log.Print("-> now, the product with ID 1, read from the FIRST database transaction, has the updated value from the SECOND database transaction after it has been committed")

	return nil
}
