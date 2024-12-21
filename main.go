package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/template/html/v2"
	_ "github.com/mattn/go-sqlite3"
	slogfiber "github.com/samber/slog-fiber"
	"github.com/usysrc/usystack/filter"
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
	app.Get("/", indexHandler)
	app.Post("/add-item", addItem)

	// Start server
	app.Listen(":3000")
}

// add items to the db
func addItem(c *fiber.Ctx) error {
	slog.Info(string(c.Body()))
	var newItem Item
	if err := c.BodyParser(&newItem); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return err
	}

	_, err := db.Exec("INSERT into items (name) VALUES ($1)", newItem.Name)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return err
	}
	err = listItems(c)
	return err
}

func getItems(c *fiber.Ctx) []Item {
	rows, err := db.Query("SELECT id, name FROM items")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return nil
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name)
		if err != nil {
			slog.Error(err.Error())
		}

		items = append(items, item)
	}
	return items
}

// list the items
func listItems(c *fiber.Ctx) error {
	err := c.Render("list", fiber.Map{
		"Items": getItems(c),
	})
	if err != nil {
		slog.Error(err.Error())
	}
	return err
}

func indexHandler(c *fiber.Ctx) error {
	err := c.Render("index", fiber.Map{
		"Items": getItems(c),
	})
	if err != nil {
		slog.Error(err.Error())
	}
	return err
}

type Item struct {
	ID   int           `json:"id"`
	Name template.HTML `json:"name"`
}
