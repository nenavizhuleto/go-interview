package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func main() {
	if len(os.Args) < 2 {
		log.Fatal("pass start and end date in 'YYYY-MM-DD' format")
	}
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

	log.Println("Connected to mysql database")

	start := os.Args[1]
	end := os.Args[2]
	res, err := getReport(start, end)
	if err != nil {
		log.Fatal(err)
	}

	printReport(res)
}

func printReport(res []Row) {
	report := make(Report)
	for _, row := range res {
		hours := 24 - row.Date.Hour()
		date := row.Date.Format(time.DateOnly)
		if _, ok := report[date]; !ok {
			report[date] = make(ReportEntry)
			report[date][row.Name] = hours
		} else {
			// if val[row.Name] < hours {
			// 	continue
			// }
			if report[date][row.Name] > hours {
				continue
			} else {
				report[date][row.Name] = hours
			}
		}
	}

	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	defer w.Flush()

	fmt.Fprintf(w, "\n %s\t%s\t%s\t", "Date", "Name", "Hours")

	keys := make([]string, 0, len(report))
	for k := range report {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		a, err := time.Parse(time.DateOnly, keys[i])
		if err != nil {
			log.Fatal("Error sorting dates")
		}
		b, err := time.Parse(time.DateOnly, keys[j])
		if err != nil {
			log.Fatal("Error sorting dates")
		}
		return a.Compare(b) < 0
	})

	names := []string{
		"Храбров",
	}

	for _, date := range keys {
		for name, hours := range report[date] {
			pass := true
			for _, n := range names {
				if strings.Contains(strings.ToLower(name), strings.ToLower(n)) {
					pass = true
				}
			}

			if !pass {
				continue
			}
			fmt.Fprintf(w, "\n %s\t%s\t%d\t", date, name, hours)
		}
	}
	fmt.Fprintf(w, "\n")
}

type (
	Report      map[string]ReportEntry
	ReportEntry map[string]int
)

type Row struct {
	Date time.Time
	Name string
}

func getReport(start, end string) ([]Row, error) {
	query := `SELECT l.datetime as date, e.name as employee FROM log as l 
        JOIN employee as e ON (l.employee_id = e.id)
        WHERE e.company_id = 1 AND
        l.datetime >= DATE(?) AND l.datetime <= DATE(?) 
        ORDER BY l.datetime`

	rows, err := db.Query(query, start, end)
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
