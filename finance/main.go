package main

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	InitDatabase()

	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e := echo.New()
	e.Renderer = t
	e.Debug = true
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		account, err := DbSelectAccountByName("default")
		if err != nil {
			return err
		}

		log.Printf("Account: %v\n", account)
		return c.Render(http.StatusOK, "index.html", account)
	})

	e.POST("/transaction/:name", func(c echo.Context) error {
		name := c.Param("name")
		account, err := DbSelectAccountByName(name)
		if err != nil {
			return err
		}

		value_str := c.FormValue("value")
		log.Printf("VALUE: %#v", value_str)

		value, err := strconv.ParseFloat(value_str, 64)
		if err != nil {
			return err
		}

		account.Balance += value
		DbUpdateAccountBalance(account.ID, account.Balance)
		buf := bytes.NewBufferString("")
		t.Render(buf, "balance", account)

		return c.Render(http.StatusOK, "balance", account)
	})

	e.Start(":3000")
}
