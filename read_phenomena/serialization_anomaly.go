package read_phenomena

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type SerializationAnomaly struct {
	Migrate func() error
	DB      *sql.DB
}

func (s *SerializationAnomaly) Run() (err error) {
	if s.Migrate == nil {
		return fmt.Errorf("migrate function is required")
	}
	if s.DB == nil {
		return fmt.Errorf("database connection is required")
	}

	// migrate
	log.Println("[o] migrating database...")
	err = s.Migrate()
	if err != nil {
		return
	}
	log.Print("[v] migrate success...\n\n")

	log.Println("1. as a starting point, let's check the data with the ID 1 or 3 in the database")
	rows, err := s.DB.Query("SELECT id, name, quantity FROM products WHERE id IN (1, 3)")
	if err != nil {
		return
	}
	rowNumber := 0
	totalQuantity := 0
	for rows.Next() {
		rowNumber++
		var product Product
		err = rows.Scan(&product.ID, &product.Name, &product.Quantity)
		if err != nil {
			return
		}

		totalQuantity += int(product.Quantity)
		fmt.Printf("(%d)----------\nID: %d\nName: %s\nQuantity: %d\n-------------\n\n", rowNumber, product.ID, product.Name, product.Quantity)
	}
	log.Printf("-> remember that we get this lists of product above, and the sum of their quantity is %d\n\n", totalQuantity)

	// create first transaction
	log.Print("2. we are going to create our FIRST database transaction")
	tx1, err := s.DB.Begin()
	if err != nil {
		return
	}
	log.Print("-> FIRST TRANSACTION CREATED\n\n")

	// create second transaction
	log.Print("3. then, we are going to create our SECOND database transaction")
	tx2, err := s.DB.Begin()
	if err != nil {
		if e := tx1.Rollback(); e != nil {
			return e
		}

		return
	}
	log.Print("-> SECOND TRANSACTION CREATED\n\n")

	log.Println("4. let's update the quantity of product with ID 1 by 10 from the FIRST database transaction")
	_, err = tx1.Exec("UPDATE products SET quantity = 10 WHERE id = 1")
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

	log.Println("5. let's see the sum of product quantities with ID 1 or 3 from the FIRST database transaction")
	var sumFromTx1 int64
	err = tx1.QueryRow("SELECT SUM(quantity) FROM products WHERE id IN (1, 3)").Scan(&sumFromTx1)
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
	log.Printf("-> SUM: %d\n", sumFromTx1)
	log.Printf("-> remember that now we get the sum of product quantities with ID 1 or 3 from the FIRST database transaction is %d\n\n", sumFromTx1)

	log.Println("6. let's update update the quantity of product with ID 3 by 5 from the SECOND database transaction")
	_, err = tx2.Exec("UPDATE products SET quantity = 5 WHERE id = 3")
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

	log.Println("7. let's see the sum of product quantities with ID 1 or 3 from the SECOND database transaction")
	var sumFromTx2 int64
	err = tx2.QueryRow("SELECT SUM(quantity) FROM products WHERE id IN (1, 3)").Scan(&sumFromTx2)
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
	log.Printf("-> SUM: %d\n", sumFromTx2)
	log.Printf("-> remember that now we get the sum of product quantities with ID 1 or 3 from the SECOND database transaction is %d\n\n", sumFromTx2)

	log.Println("8. then, we are going to commit the FIRST and SECOND database transaction")
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
	err = tx2.Commit()
	if err != nil {
		var e2 error
		e2 = tx2.Rollback()

		if e2 != nil {
			return e2
		}

		return
	}
	log.Print("-> TRANSACTION COMMITTED\n\n")

	str := strings.Builder{}
	str.WriteString("As you can see, after committing both database transaction no error occured.")
	str.WriteString(" This phenomenon called 'Serialization Anomaly'. These two transactions work with the same data and the order of commits will affect results, which could lead to unpredictable outcomes.")

	fmt.Print(str.String())

	return nil
}
