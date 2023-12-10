package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"pg-isolation/isolation"
	"pg-isolation/read_phenomena"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var conn *sql.DB
var clear map[string]func()

func init() {
	clear = make(map[string]func())

	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["darwin"] = func() {
		cmd := exec.Command("clear") //Mac
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	if _, ok := clear[runtime.GOOS]; !ok {
		log.Fatal("Your platform is unsupported!")
	}
}

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	var topic Topic
	var readPhenomena read_phenomena.ReadPhenomena
	var isolationLevel isolation.Isolation

	config := ParseConfig()
	db := &DB{
		Config: config,
	}
	conn, err = db.CreateConn()
	if err != nil {
		log.Fatal("[x] failed to connect to database", err)
	}

MainLoop:
	for {
		Clear()

		// choose topic
		for {
			topic, err = chooseTopic()
			if err == nil && topic != TOPIC_UNSPECIFIED {
				break
			}

			fmt.Println("\nChoose the right Topic option!")
		}

		if topic == TOPIC_READ_PHENOMENA {
			for {
				readPhenomena, err = chooseReadPhenomena()
				if err == nil && readPhenomena != read_phenomena.READ_PHENOMENA_UNSPECIFIED {
					break
				}

				fmt.Println("\nChoose the right Read Phenomena option!")
			}

			// repeatable read
			if readPhenomena == read_phenomena.READ_PHENOMENA_NON_REPEATABLE_READ {
				Clear()

				nonRepeatableRead := read_phenomena.NonRepeatableRead{
					DB:      conn,
					Migrate: db.Migrate,
				}
				err := nonRepeatableRead.Run()
				if err != nil {
					log.Fatal(err)
				}
			} else if readPhenomena == read_phenomena.READ_PHENOMENA_PHANTOM_READ {
				Clear()

				phantomRead := read_phenomena.PhantomRead{
					DB:      conn,
					Migrate: db.Migrate,
				}
				err := phantomRead.Run()
				if err != nil {
					log.Fatal(err)
				}
			} else if readPhenomena == read_phenomena.READ_PHENOMENA_SERIALIZATION_ANOMALY {
				Clear()

				serializationAnomaly := read_phenomena.SerializationAnomaly{
					DB:      conn,
					Migrate: db.Migrate,
				}
				err := serializationAnomaly.Run()
				if err != nil {
					log.Fatal(err)
				}
			}

		} else if topic == TOPIC_PG_ISOLATION_LEVEL {
			for {
				isolationLevel, err = chooseIsolationLevel()
				if err == nil && isolationLevel != isolation.ISOLATION_UNSPECIFIED {
					break
				}

				fmt.Println("\nChoose the right Isolation Level option!")
			}

			// read committed
			if isolationLevel == isolation.ISOLATION_READ_COMMITTED {
				Clear()

				readCommitted := isolation.RadCommitted{
					DB:      conn,
					Migrate: db.Migrate,
				}
				err = readCommitted.Run()
				if err != nil {
					log.Fatal(err)
				}
			} else if isolationLevel == isolation.ISOLATION_REPEATABLE_READ {
				// TODO: Implement this
			} else if isolationLevel == isolation.ISOLATION_SERIALIZABLE {
				// TODO: Implement this
			}
		}

	RunAgainLoop:
		for {
			var confirmMessage = "\nDo you want to see another example?"
			fmt.Println(confirmMessage)
			fmt.Println(strings.Repeat("-", len(confirmMessage)))

			fmt.Println("[0] No")
			fmt.Println("[1] Yes")
			fmt.Print("\n-> ")

			var runAgain int
			_, err = fmt.Scan(&runAgain)
			if err != nil {
				log.Fatal(err)
			}

			switch runAgain {
			case 0:
				break MainLoop
			case 1:
				break RunAgainLoop
			default:

				fmt.Println("\nChoose the right option!")
			}
		}

	}
}

func chooseTopic() (Topic, error) {
	var msg1 = "\nChoose one of the topics below:"
	fmt.Println(msg1)
	fmt.Println(strings.Repeat("-", len(msg1)))

	fmt.Println("[1] Read Phenomena")
	fmt.Println("[2] PG Isolation Level")
	fmt.Print("\n-> ")

	var topic Topic
	_, err := fmt.Scan(&topic)
	if err != nil || !topic.Valid() {
		return TOPIC_UNSPECIFIED, ErrInvalidInput
	}

	return topic, nil
}

func chooseReadPhenomena() (read_phenomena.ReadPhenomena, error) {
	var msg1 = "\nChoose one of the Read Phenomena below:"
	fmt.Println(msg1)
	fmt.Println(strings.Repeat("-", len(msg1)))

	fmt.Println("[1] Non Repeatable Read")
	fmt.Println("[2] Phantom Read")
	fmt.Println("[3] Serialization Anomaly")
	fmt.Print("\n-> ")

	var readPhenomena read_phenomena.ReadPhenomena
	_, err := fmt.Scan(&readPhenomena)
	if err != nil || !readPhenomena.Valid() {
		return read_phenomena.READ_PHENOMENA_UNSPECIFIED, ErrInvalidInput
	}

	return readPhenomena, nil
}

func chooseIsolationLevel() (isolation.Isolation, error) {
	var msg1 = "\nChoose one of the Isolation Level below:"
	fmt.Println(msg1)
	fmt.Println(strings.Repeat("-", len(msg1)))

	fmt.Println("[1] Read Committed")
	fmt.Println("[2] Repeatable Read")
	fmt.Println("[3] Serializable")
	fmt.Print("\n-> ")

	var isolationLevel isolation.Isolation
	_, err := fmt.Scan(&isolationLevel)
	if err != nil || !isolationLevel.Valid() {
		return isolation.ISOLATION_UNSPECIFIED, ErrInvalidInput
	}

	return isolationLevel, nil
}

func Clear() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		log.Fatal("Your platform is unsupported!")
	}
}
