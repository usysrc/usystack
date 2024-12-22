package model

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Item struct {
	ID   int           `json:"id"`
	Name template.HTML `json:"name"`
}

func GetAllItems(c *fiber.Ctx) ([]Item, error) {
	rows, err := db.Query("SELECT id, name FROM items")
	if err != nil {
		slog.Error(err.Error())
		c.Status(http.StatusInternalServerError)
		return nil, err
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
	return items, nil
}

func NewItem(c *fiber.Ctx, newItem Item) error {
	_, err := db.Exec("INSERT into items (name) VALUES ($1)", newItem.Name)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return err
	}
	return nil
}
