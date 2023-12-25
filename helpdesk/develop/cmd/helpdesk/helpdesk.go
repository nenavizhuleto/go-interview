package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"helpdesk/internals/api"
	"helpdesk/internals/api/admin"
	"helpdesk/internals/data"
	"helpdesk/internals/megaplan"
	"helpdesk/internals/util"
	"helpdesk/telegram"
)

func main() {
	// --- Environment ---
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	// --- Telegram Bot ---
	if err := telegram.InitDefault(util.MustGetEnvVar("TELEGRAM_TOKEN")); err != nil {
		panic(err)
	}

	go telegram.Bot.Run()

	// --- Database ---
	if db, err := data.NewDB(); err != nil {
		log.Fatalf("Error initializing SQlite database: %v", err)
	} else {
		data.DB = db
	}
	data.MustConnectMongo("helpdesk")

	// --- Megaplan ---
	opts := megaplan.NewAuthOpt(util.MustGetEnvVar("MEGAPLAN_USER"), util.MustGetEnvVar("MEGAPLAN_PASSWORD"), util.MustGetEnvVar("MEGAPLAN_RESPONSIBLE"))
	megaplan.MP = megaplan.New(util.MustGetEnvVar("MEGAPLAN_URL"), opts).MustAuthenticateWithPassword(opts)

	// --- Fiber ---
	app := fiber.New(fiber.Config{
		ErrorHandler: api.HandleError,
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	apiRouter := app.Group("/api")
	apiRouter.Post("/megaplan/event", api.HandleMegaplanEvent)

	// --- Authentication ---
	apiRouter.Get("/auth/token", api.GetToken)
	apiRouter.Post("/auth/register", api.Register)

	// --- User ---
	hd := apiRouter.Group("/helpdesk")
	hd.Use(api.UserMiddleware)
	hd.Get("/profile", api.GetUserProfile)
	hd.Get("/telegram", api.GetUserTelegram)
	hd.Post("/telegram", api.NewUserTelegram)
	hd.Get("/tasks", api.GetUserTasks)
	hd.Post("/tasks", api.NewUserTask)
	hd.Get("/tasks/:id", api.GetUserTask)
	hd.Get("/tasks/:id/comments", api.GetUserTaskComments)
	hd.Post("/tasks/:id/comments", api.NewUserTaskComment)

	// --- Admin ---
	adm := apiRouter.Group("/admin")
	adm.Use(admin.AdminMiddleware)
	// --- Company
	adm.Get("/company", admin.GetCompanies)
	adm.Get("/company/:id", admin.GetCompany)
	adm.Post("/company", admin.NewCompany)
	adm.Delete("/company/:id", admin.DeleteCompany)
	// --- Branch
	adm.Get("/branch", admin.GetBranches)
	adm.Get("/branch/:id", admin.GetBranch)
	adm.Post("/branch/:company_id", admin.NewBranch)
	adm.Delete("/branch/:id", admin.DeleteBranch)
	// --- User
	adm.Get("/user", admin.GetUsers)
	adm.Get("/user/:id", admin.GetUser)
	adm.Post("/user", admin.NewUser)
	adm.Delete("/user/:id", admin.DeleteUser)

	log.Fatal(app.Listen(":3000"))
}
