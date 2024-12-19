package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	_ "github.com/mattn/go-sqlite3"
	slogfiber "github.com/samber/slog-fiber"
)

var db *sql.DB

func main() {
	// Connect to PostgreSQL
	conn, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Fatal(err)
	}
	db = conn
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
	engine.Reload(true)

	// Start fiber
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Add structured logging middleware
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	app.Use(slogfiber.New(logger))

	// Define routes
	app.Get("/", indexHandler)
	app.Post("/add-item", itemHandler)

	// Start server
	app.Listen(":3000")
}

// add items to the db
func itemHandler(c *fiber.Ctx) error {
	fmt.Print(string(c.Body()))
	var newItem Item
	if err := c.BodyParser(&newItem); err != nil {
		c.Status(500)
		return err
	}

	_, err := db.Exec("INSERT into items (name) VALUES ($1)", newItem.Name)
	if err != nil {
		c.Status(500)
		return err
	}
	return indexHandler(c)
}

func indexHandler(c *fiber.Ctx) error {
	rows, err := db.Query("SELECT id, name FROM items")
	if err != nil {
		return err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name)
		if err != nil {
			slog.Error(err.Error())
			return fmt.Errorf("failed to scan row: %w", err)
		}
		items = append(items, item)
	}

	return c.Render("index", fiber.Map{
		"Items": items,
	})
}

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
