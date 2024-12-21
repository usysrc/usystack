package main

import (
	"database/sql"
	"log"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/template/html/v2"
	_ "github.com/mattn/go-sqlite3"
	slogfiber "github.com/samber/slog-fiber"

	"github.com/usysrc/usystack/controller"
	"github.com/usysrc/usystack/filter"
)

func main() {
	// Connect to PostgreSQL
	conn, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Fatal(err)
	}
	db := conn
	defer db.Close()

	// load 'init.sql' and execute it
	file, err := os.ReadFile("init.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(string(file))
	if err != nil {
		log.Fatal(err)
	}

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

	// Define routes
	itemHandler := controller.NewItemHandler(db)
	app.Get("/", itemHandler.IndexHandler)
	app.Post("/add-item", itemHandler.AddItem)

	// Start server
	if err := app.Listen(":3000"); err != nil {
		slog.Error(err.Error())
	}
}
