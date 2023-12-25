package main

import (
	"log"
	"os/exec"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/google/uuid"
)

func main() {

	engine := html.New("./views", ".html")
	engine.AddFunc("string", func(b []byte) string {
		return string(b)
	})
	engine.Reload(true)
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		cmd := exec.Command("ls", "-la")
		return c.Render("index", cmd)
	})

	app.Post("/exec", execCmd)
	app.Get("/listen/:id", listen)

	app.Listen(":3000")
}

var Processes = make(map[string]*exec.Cmd)

func execCmd(c *fiber.Ctx) error {
	var body struct {
		Input string
	}

	log.Printf("%#v", c.BodyRaw())

	if err := c.BodyParser(&body); err != nil {
		return err
	}

	args := strings.Split(body.Input, " ")
	cmd := exec.Command(args[0], args[1:]...)

	id := uuid.NewString()
	Processes[id] = cmd
	log.Printf("%s", id)

	return c.Render("output", id)

}

func listen(c *fiber.Ctx) error {
	id := c.Params("id")

	cmd, ok := Processes[id]
	if !ok {
		return fiber.ErrNotFound
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	bw := c.Response().BodyWriter()
	for {
		out := make([]byte, 1024)
		n, err := stdout.Read(out)
		if err != nil {
			break
		}
		bw.Write(out[:n])
	}

	return c.SendStatus(286)
}
