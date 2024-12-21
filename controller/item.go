package controller

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/usysrc/usystack/model"
)

type ItemHandler struct {
	Items model.ItemStore
}

func NewItemHandler(db *sql.DB) ItemHandler {
	return ItemHandler{
		Items: model.NewItemStore(db),
	}
}

// add an item to the db
func (ih *ItemHandler) AddItem(c *fiber.Ctx) error {
	slog.Debug(string(c.Body()))
	var newItem model.Item
	if err := c.BodyParser(&newItem); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return err
	}
	err := ih.Items.NewItem(c, newItem)
	if err != nil {
		return err
	}
	err = ih.ListItems(c)
	return err
}

// list the items
func (ih *ItemHandler) ListItems(c *fiber.Ctx) error {
	items, err := ih.Items.GetAllItems(c)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	err = c.Render("list", fiber.Map{
		"Items": items,
	})
	if err != nil {
		slog.Error(err.Error())
	}
	return err
}

// write the index
func (ih *ItemHandler) IndexHandler(c *fiber.Ctx) error {
	items, err := ih.Items.GetAllItems(c)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	err = c.Render("index", fiber.Map{
		"Items": items,
	})
	if err != nil {
		slog.Error(err.Error())
	}
	return err
}
