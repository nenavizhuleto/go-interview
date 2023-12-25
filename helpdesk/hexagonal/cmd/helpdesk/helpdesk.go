package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	apiv3 "helpdesk/internals/api/v3"
	"helpdesk/internals/data"
	"helpdesk/internals/megaplan"
	"helpdesk/internals/util"
	"helpdesk/telegram"
)

func initTelegram() {
	if err := telegram.InitDefault(util.MustGetEnvVar("TELEGRAM_TOKEN")); err != nil {
		panic(err)

	}

	go telegram.Bot.Run()
}

func main() {
	initEnv()
	initTelegram()
	initDb()
	initMegaplan()

	// initializing fiber application
	app := fiber.New(fiber.Config{
		ErrorHandler: apiv3.HandleError,
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	apiRouter := app.Group("/api")
	apiRouter.Post("/megaplan/event", apiv3.HandleMegaplanEvent)

	apiV3 := apiRouter.Group("/v3")

	apiV3.Get("/identity", apiv3.GetIdentity)
	apiV3.Post("/register", apiv3.Register)

	apiv3.SetUsersRoutes("/users", apiV3)
	apiv3.SetDevicesRoutes("/devices", apiV3)
	apiv3.SetCompaniesRoutes("/company", apiV3)
	apiv3.SetBranchesRoutes("/branch", apiV3)
	apiv3.SetNetworksRoutes("/network", apiV3)
	apiv3.SetTasksRoutes("/tasks", apiV3)

	log.Fatal(app.Listen(":3000"))
}

func initEnv() {
	// initializing environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func initDb() {
	// initializing database connection
	if db, err := data.NewDB(); err != nil {
		log.Fatalf("Error initializing SQlite database: %v", err)
	} else {
		data.DB = db
	}

	data.MustConnectMongo("helpdesk")
}

func initMegaplan() {
	// initializing megaplan connection
	opts := megaplan.NewAuthOpt(util.MustGetEnvVar("MEGAPLAN_USER"), util.MustGetEnvVar("MEGAPLAN_PASSWORD"), util.MustGetEnvVar("MEGAPLAN_RESPONSIBLE"))
	megaplan.MP = megaplan.New(util.MustGetEnvVar("MEGAPLAN_URL"), opts).MustAuthenticateWithPassword(opts)
}
