package read_phenomena

import (
	"database/sql"
	"fmt"
	"log"
)

type NonRepeatableRead struct {
	Migrate func() error
	DB      *sql.DB
}

func (n *NonRepeatableRead) Run() (err error) {
	if n.Migrate == nil {
		return fmt.Errorf("migrate function is required")
	}
	if n.DB == nil {
		return fmt.Errorf("database connection is required")
	}

	// migrate
	log.Println("[o] migrating database...")
	err = n.Migrate()
	if err != nil {
		return
	}
	log.Print("[v] migrate success...\n\n")

	// create first transaction
	log.Print("1. we are going to create our FIRST database transaction")
	tx1, err := n.DB.Begin()
	if err != nil {
		return
	}
	log.Print("-> FIRST TRANSACTION CREATED\n\n")

	// create second transaction
	log.Print("2. then, we are going to create our SECOND database transaction")
	tx2, err := n.DB.Begin()
	if err != nil {
		if e := tx1.Rollback(); e != nil {
			return e
		}

		return
	}
	log.Print("-> SECOND TRANSACTION CREATED\n\n")

	log.Println("3. read the product with ID 1 from the FIRST database transaction")
	var product1 Product
	err = tx1.QueryRow("SELECT name, quantity FROM products WHERE id = 1").Scan(&product1.Name, &product1.Quantity)
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
	log.Printf("-> we get '%s' with a quantity of %d\n\n", product1.Name, product1.Quantity)

	log.Println("4. update the product quantity with ID 1 from the SECOND database transaction to 5, and then commit immediately")
	_, err = tx2.Exec("UPDATE products SET quantity = 5")
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
	log.Print("-> ROW UPDATED\n\n")

	log.Println("5. update the product with ID 1 from the FIRST database transaction to 3, and then commit immediately")
	_, err = tx1.Exec("UPDATE products SET quantity = 3 WHERE id = 1")
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
	log.Print("-> ROW UPDATED\n\n")

	log.Println("6. now read the final product data with ID 1")
	var finalProduct Product
	err = n.DB.QueryRow("SELECT name, quantity FROM products WHERE id = 1").Scan(&finalProduct.Name, &finalProduct.Quantity)
	if err != nil {
		return
	}

	fmt.Printf("Name: %s\nQuantity: %d\n\n", finalProduct.Name, finalProduct.Quantity)

	log.Printf(`As you can see, the product now has a quantity of %d successfully updated from the FIRST database transaction without checking if that row has been updated from the SECOND database transaction`, finalProduct.Quantity)

	return nil
}
