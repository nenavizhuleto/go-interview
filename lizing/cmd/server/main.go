package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

const PASSPHRASE string = "passphrase"
const LIZING int = 1
const RAVIS int = 3

func ConnectDB() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "2fa",
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

}

type Row struct {
	Date time.Time `json:"date"`
	Name string    `json:"name"`
}

func GetData(company_id int, start, end time.Time) ([]Row, error) {
	query := `SELECT l.datetime, e.name FROM log as l
			JOIN employee as e ON (l.employee_id = e.id)
			WHERE e.company_id = ? AND
			l.datetime >= DATE(?) AND l.datetime <= DATE(?)
			ORDER BY l.datetime`

	rows, err := db.Query(query, company_id, start, end)
	if err != nil {
		return nil, err
	}

	var res []Row

	for rows.Next() {
		var datestr string
		var name string
		rows.Scan(&datestr, &name)
		date, err := time.Parse(time.DateTime, datestr)
		if err != nil {
			return nil, err
		}
		res = append(res, Row{
			Date: date,
			Name: name,
		})
	}
	defer rows.Close()

	return res, nil
}

type GetDataParams struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func GetError(message string) []byte {
	body, err := json.Marshal(map[string]interface{}{
		"status":  http.StatusNotAcceptable,
		"message": message,
	})
	if err != nil {
		body = []byte("Error")
	}

	return body
}

func main() {
	ConnectDB()

	http.HandleFunc("/lizing", func(w http.ResponseWriter, r *http.Request) {
		pass := r.Header.Get("passphrase")
		if pass != PASSPHRASE {
			w.WriteHeader(http.StatusForbidden)
			log.Printf("Failed request from: %s ", r.RemoteAddr)
			return
		}

		log.Printf("Authorized access from: %s", r.RemoteAddr)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		var params GetDataParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write(GetError(err.Error()))
			return
		}

		start, _ := time.Parse(time.DateOnly, params.Start)
		end, _ := time.Parse(time.DateOnly, params.End)

		log.Printf("Params: start = %s, end = %s", params.Start, params.End)

		rows, err := GetData(LIZING, start, end)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(GetError(err.Error()))
			return
		}

		log.Printf("rows: %#v", rows)

		body, err:= json.Marshal(rows)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(GetError(err.Error()))
			return
		}

		log.Printf("Successfully fullfiled request")
		
		w.Write(body)
		return
	})

	http.HandleFunc("/ravis", func(w http.ResponseWriter, r *http.Request) {
		pass := r.Header.Get("passphrase")
		if pass != PASSPHRASE {
			w.WriteHeader(http.StatusForbidden)
			log.Printf("Failed request from: %s ", r.RemoteAddr)
			return
		}

		log.Printf("Authorized access from: %s", r.RemoteAddr)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		var params GetDataParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write(GetError(err.Error()))
			return
		}

		start, _ := time.Parse(time.DateOnly, params.Start)
		end, _ := time.Parse(time.DateOnly, params.End)

		log.Printf("Params: start = %s, end = %s", params.Start, params.End)

		rows, err := GetData(RAVIS, start, end)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(GetError(err.Error()))
			return
		}

		log.Printf("rows: %#v", rows)

		body, err:= json.Marshal(rows)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(GetError(err.Error()))
			return
		}

		log.Printf("Successfully fullfiled request")
		
		w.Write(body)
		return
	})

	port := os.Getenv("PORT")
	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))

}
