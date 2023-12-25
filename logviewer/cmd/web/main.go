package main

import (
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

var LogDirectory = "logs/"

var files map[string]string

func traverse(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			err := traverse(dir + entry.Name() + "/")
			if err != nil {
				return err
			}
			continue
		}
		filepath := dir + entry.Name()
		content, err := os.ReadFile(filepath)
		if err != nil {
			return err
		}

		files[filepath] = string(content)
	}
	return nil
}

func read() error {
	files = make(map[string]string)

	err := traverse(LogDirectory)
	if err != nil {
		return err
	}

	return nil

}

func main() {
	engine := html.New("views", ".html")
	engine.Reload(true)

	if err := read(); err != nil {
		panic(err)
	}

	config := fiber.Config{
		Views: engine,
	}
	app := fiber.New(config)
	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		var dirs []string
		for key := range files {
			dirs = append(dirs, key)
		}

		return c.Render("index", fiber.Map{
			"Dirs": dirs,
		})
	})

	app.Post("/search", func(c *fiber.Ctx) error {
		term := c.FormValue("search")

		counts := make(fiber.Map)
		for key, value := range files {
			counts[key] = strings.Count(value, term)
		}

		return c.Render("search", counts)

	})

	app.Post("/content/*", func(c *fiber.Ctx) error {
		filename := c.Params("*")
		log.Printf(filename)
		file, ok := files[filename]
		if !ok {
			return c.SendString("not found")
		}

		return c.SendString(string(file))
	})

	app.Listen(":3000")
}
