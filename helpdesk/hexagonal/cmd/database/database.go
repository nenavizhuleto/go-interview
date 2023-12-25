package main

import (
	"database/sql"
	"log"
	"os"

	"helpdesk/internals/data"
)

func initDatabase() {
	log.Printf("Initializing database with sql/init.sql")
	buf, err := os.ReadFile("sql/init.sql")
	if err != nil {
		log.Fatalf("Couldn't read init.sql file: %v", err)
	}
	query := string(buf)

	log.Println(query)
	db, err := sql.Open("sqlite3", data.DBNAME)
	if err != nil {
		log.Fatalf("Couldn't initialize database: %v", err)
	}

	if _, err = db.Exec(query); err != nil {
		log.Fatalf("Error executing sql query: %v", err)
	}
}

func dropDatabase() {
	log.Printf("Dropping database with sql/drop.sql")
	buf, err := os.ReadFile("sql/drop.sql")
	if err != nil {
		log.Fatalf("Couldn't read drop.sql file: %v", err)
	}
	query := string(buf)

	log.Println(query)
	db, err := sql.Open("sqlite3", data.DBNAME)
	if err != nil {
		log.Fatalf("Couldn't initialize database: %v", err)
	}

	if _, err = db.Exec(query); err != nil {
		log.Fatalf("Error executing sql query: %v", err)
	}
}

func main() {
	args := os.Args[1:]
	log.Printf("Args: %v", args)
	if len(args) < 1 {
		log.Fatalf("Not enough arguments: pass init or drop")
	}

	action := args[0]

	if action == "init" {
		initDatabase()
	}

	if action == "drop" {
		dropDatabase()
	}
}
