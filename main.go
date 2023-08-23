package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/jackc/pgx/v4"
)

var db *pgx.Conn

func main() {
	// Connect to PostgreSQL
	conn, err := pgx.Connect(context.Background(), "postgresql://username:password@usystack-db/dbname")
	if err != nil {
		log.Fatal(err)
	}
	db = conn
	defer db.Close(context.Background())

	// Initialize standard Go html template engine
	engine := html.New("./views", ".html")

	// Start fiber
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Define routes
	app.Get("/", indexHandler)

	// Start server
	log.Fatal(app.Listen(":3000"))
}

func indexHandler(c *fiber.Ctx) error {
	rows, err := db.Query(context.Background(), "SELECT id, name FROM items")
	if err != nil {
		return err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name)
		if err != nil {
			return err
		}
		items = append(items, item)
	}

	return c.Render("index", fiber.Map{
		"Items": items,
	})
}

type Item struct {
	ID   int
	Name string
}
