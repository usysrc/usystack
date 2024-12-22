package main

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/template/html/v2"
	_ "github.com/mattn/go-sqlite3"
	slogfiber "github.com/samber/slog-fiber"

	"github.com/usysrc/usystack/controller"
	"github.com/usysrc/usystack/filter"
	"github.com/usysrc/usystack/model"
)

func main() {
	// Initialize standard Go html template engine
	engine := html.New("./views", ".html")
	engine.Funcmap = map[string]any{
		"markdown": filter.MarkdownFilter,
	}

	engine.Reload(true)

	// Start fiber
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Ignore favicon requests
	app.Use(favicon.New())

	// Add structured logging middleware
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	app.Use(slogfiber.New(logger))

	// Add session middleware
	controller.CreateSessionStore()

	// Define routes
	model.Connect()
	defer model.Close()
	app.Get("/", controller.IndexHandler)
	app.Get("/:id", controller.SingleHandler)
	app.Post("/add-item", controller.AddItem)
	app.Post("/login", controller.Login)
	app.Post("/logout", controller.Logout)

	// Start server
	if err := app.Listen(":3000"); err != nil {
		slog.Error(err.Error())
	}
}
