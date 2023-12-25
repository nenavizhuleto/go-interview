package main

import (
	"database/sql"
	"html/template"
	"io"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

type Metric struct {
	Site      string
	Address   string
	Timestamp string
}

var metrics = make(map[string][]Metric)

type Server struct {
	DB         *sql.DB
	ListenAddr string
	App        *echo.Echo
}

var server Server

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	db, err := sql.Open("sqlite3", "metrics.db")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	db.Exec(`CREATE TABLE IF NOT EXISTS metrics (
        site TEXT NOT NULL,
        ip TEXT NOT NULL,
        timestamp TEXT NOT NULL
    )`)
	t := &Template{
		templates: template.Must(template.ParseGlob("*.html")),
	}
	app := echo.New()
	app.Renderer = t
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())
	server.App = app
	server.DB = db
	server.ListenAddr = ":3000"

	app.GET("/monitor", handleMonitor)
	app.GET("/monitor/statistics", handleMonitorStats)
	app.GET("/monitor/metrics", handleMonitorMetrics)
	app.GET("/metrics", handleCollectMetrics)
	app.GET("/cdn", handleCDNMetricScript)

	app.Start(server.ListenAddr)
}

func handleMonitorMetrics(c echo.Context) error {
	var metrics []Metric
	rows, err := server.DB.Query("SELECT * FROM metrics")
	if err != nil {
		return err
	}
	for rows.Next() {
		var metric Metric
		rows.Scan(&metric.Site, &metric.Address, &metric.Timestamp)
		metrics = append(metrics, metric)
	}
	defer rows.Close()
	return c.Render(200, "metrics.html", metrics)
}

type Stat struct {
	Site   string
	Visits int
}

func handleMonitorStats(c echo.Context) error {
	rows, err := server.DB.Query("SELECT site, count(site) FROM metrics GROUP BY site")
	if err != nil {
		return err
	}
	var stats []Stat
	for rows.Next() {
		var site string
		var visits int
		rows.Scan(&site, &visits)
		stats = append(stats, Stat{Site: site, Visits: visits})
	}
	defer rows.Close()
	return c.Render(200, "statistics.html", stats)
}

func handleMonitor(c echo.Context) error {
	return c.Render(200, "monitor.html", nil)
}

func handleCDNMetricScript(c echo.Context) error {
	return c.File("./metrics.js")
}

func handleCollectMetrics(c echo.Context) error {
	site := c.QueryParam("site")

	db := server.DB
	if _, err := db.Exec("INSERT INTO metrics VALUES (?, ?, ?)", site, c.RealIP(), time.Now().String()); err != nil {
		return err
	}

	return c.NoContent(200)
}
