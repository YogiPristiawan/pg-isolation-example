package read_phenomena

import (
	"database/sql"
	"fmt"
	"log"
)

type PhantomRead struct {
	Migrate func() error
	DB      *sql.DB
}

func (p *PhantomRead) Run() (err error) {
	if p.Migrate == nil {
		return fmt.Errorf("migrate function is required")
	}
	if p.DB == nil {
		return fmt.Errorf("database connection is required")
	}

	// migrate
	log.Println("[o] migrating database...")
	err = p.Migrate()
	if err != nil {
		return
	}
	log.Print("[v] migrate success...\n\n")

	// create first transaction
	log.Print("1. we are going to create our FIRST database transaction")
	tx1, err := p.DB.Begin()
	if err != nil {
		return
	}
	log.Print("-> FIRST TRANSACTION CREATED\n\n")

	// create second transaction
	log.Print("2. next, we are going to create our SECOND database transaction")
	tx2, err := p.DB.Begin()
	if err != nil {
		if e := tx1.Rollback(); e != nil {
			return e
		}

		return
	}
	log.Print("-> SECOND TRANSACTION CREATED\n\n")

	log.Println("3. read the product with quantity less than 50 from the FIRST database transaction")
	rows, err := tx1.Query("SELECT name, quantity FROM products WHERE quantity < 50")
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
	log.Println("-> we get this list of products:")
	initialRowsLength := 0
	for rows.Next() {
		initialRowsLength++
		var product Product
		err = rows.Scan(&product.Name, &product.Quantity)
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

		fmt.Printf("(%d)----------\nName: %s\nQuantity: %d\n-------------\n\n", initialRowsLength, product.Name, product.Quantity)
	}

	log.Println("4. then, insert a product with quantity of 20 from the SECOND database transaction, and then commit immediately")
	_, err = tx2.Exec("INSERT INTO products (name, quantity) VALUES ($1, $2)", "Carpet", 20)
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
	log.Print("-> ROW INSERTED\n\n")

	log.Println("5. read the product again with quantity less than 50 from the FIRST database transaction")
	rows, err = tx1.Query("SELECT name, quantity FROM products WHERE quantity < 50")
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
	log.Println("-> now, we get this list of products:")
	finalRowsLength := 0
	for rows.Next() {
		finalRowsLength++
		var product Product
		err = rows.Scan(&product.Name, &product.Quantity)
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

		fmt.Printf("(%d)----------\nName: %s\nQuantity: %d\n-------------\n\n", finalRowsLength, product.Name, product.Quantity)
	}
	err = tx1.Commit()
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

	log.Printf(`As you can see, there are %d rows of products with quantities less than 50. This is called a 'Phantom Read', because the uncommited FIRST database transaction can see matched rows inserted from the SECOND database transaction`,
		finalRowsLength)

	return nil
}
