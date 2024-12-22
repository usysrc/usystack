package controller

import (
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/usysrc/usystack/model"
)

// add an item to the db
func AddItem(c *fiber.Ctx) error {
	slog.Debug(string(c.Body()))
	var newItem model.Item
	if err := c.BodyParser(&newItem); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		slog.Error(err.Error())
		return err
	}
	err := model.NewItem(c, newItem)
	if err != nil {
		return err
	}
	err = ListItems(c)
	return err
}

// list the items
func ListItems(c *fiber.Ctx) error {
	items, err := model.GetAllItems(c)
	if err != nil {
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
func IndexHandler(c *fiber.Ctx) error {
	items, err := model.GetAllItems(c)
	if err != nil {
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

// write single item
func SingleHandler(c *fiber.Ctx) error {
	type Param struct {
		ID int `json:"id"`
	}
	param := Param{}
	if err := c.ParamsParser(&param); err != nil {
		slog.Error(err.Error())
		return err
	}
	item, err := model.GetItem(c, param.ID)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	err = c.Render("single", fiber.Map{
		"Item": item,
	})
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}