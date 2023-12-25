package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const database_filename string = "finance.db"

const sql_accounts_create string = `
CREATE TABLE IF NOT EXISTS accounts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE,
	balance FLOAT NOT NULL DEFAULT(0)
);
`

const sql_transactions_create string = `
CREATE TABLE IF NOT EXISTS transactions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	account_id INTEGER NOT NULL,
	timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	description TEXT DEFAULT "",
	value FLOAT NOT NULL
);
`

var Db *sql.DB

func InitDatabase() {
	db, err := sql.Open("sqlite3", database_filename)
	if err != nil {
		log.Fatalf("error while opening %s database: %s", database_filename, err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("database ping failed: %s", err)
	}

	if _, err := db.Exec(sql_accounts_create); err != nil {
		log.Fatalf("error executing `create accounts table`: %s", err)
	} else {
		log.Println("`accounts` table created")
	}

	if _, err := db.Exec(sql_transactions_create); err != nil {
		log.Fatalf("error executing `create transactions table`: %s", err)
	} else {
		log.Println("`transactions` table created")
	}

	Db = db
	if _, err := DbSelectAccountByName("default"); err != nil {
		if _, err := DbInsertAccount("default"); err != nil {
			log.Fatalf("error creating default `account`: %s", err)
		}
	}
	log.Println("database initialized successfully")
}

func DbSelectAccount(id int64) (Account, error) {
	var account Account
	res := Db.QueryRow("SELECT * FROM accounts WHERE id = ?", id)
	if err := res.Scan(&account.ID, &account.Name, &account.Balance); err != nil {
		return account, err
	}

	return account, nil
}

func DbSelectAccountByName(name string) (Account, error) {
	var account Account
	res := Db.QueryRow("SELECT * FROM accounts WHERE name LIKE ?", name)
	if err := res.Scan(&account.ID, &account.Name, &account.Balance); err != nil {
		return account, err
	}

	return account, nil
}

func DbInsertAccount(name string) (int64, error) {
	res, err := Db.Exec("INSERT INTO accounts (name) VALUES (?)", name)
	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

func DbUpdateAccountBalance(id int64, new_balance float64) (int64, error) {
	res, err := Db.Exec("UPDATE accounts SET balance = ? WHERE id = ?", new_balance, id)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func DbSelectTransaction(id int64) (Transaction, error) {
	var t Transaction
	res := Db.QueryRow("SELECT * FROM transactions WHERE id = ?", id)
	if err := res.Scan(&t.ID, &t.AccountID, &t.Timestamp, &t.Description, &t.Value); err != nil {
		return t, err
	}

	return t, nil

}

func DbSelectAccountTransactions(account_id int64) ([]Transaction, error) {
	var ts []Transaction
	res, err := Db.Query("SELECT * FROM transactions WHERE account_id = ?", account_id)
	if err != nil {
		return ts, err
	}

	for res.Next() {
		var t Transaction
		if err := res.Scan(&t.ID, &t.AccountID, &t.Timestamp, &t.Description, &t.Value); err != nil {
			return ts, err
		}
	}

	return ts, nil
}

func DbInsertTransaction(account_id int64, description string, value float64) (int64, error) {
	log.Printf("db_insert_transaction: account: %d | description: %s | value: %f", account_id, description, value)
	res, err := Db.Exec("INSERT INTO transactions (account_id, description, value) VALUES (?, ?, ?)", account_id, description, value)
	if err != nil {
		return -1, err
	}

	return res.LastInsertId()

}
