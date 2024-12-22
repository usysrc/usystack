package model

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ID       string `json:"id"`
	LoggedIn bool
}

func GetUserByName(c *fiber.Ctx, loginData LoginData) (*User, error) {
	rows, err := db.Query("SELECT id,username, password FROM users where username = ($1)", loginData.Username)
	if err != nil {
		slog.Error(err.Error())
		c.Status(http.StatusInternalServerError)
		return nil, err
	}
	defer rows.Close()
	var user User
	if !rows.Next() {
		err := fmt.Errorf("User not found.")
		slog.Error(err.Error())
		c.Status(http.StatusNotFound)
		return nil, err
	}
	err = rows.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		slog.Error(err.Error())
		return nil, err
	}
	return &user, nil
}

func GetUserByID(c *fiber.Ctx, id int) (*User, error) {
	rows, err := db.Query("SELECT id,username, password FROM users where id = ($1)", id)
	if err != nil {
		slog.Error(err.Error())
		c.Status(http.StatusInternalServerError)
		return nil, err
	}
	defer rows.Close()
	var user User
	if !rows.Next() {
		err := fmt.Errorf("User not found.")
		slog.Error(err.Error())
		c.Status(http.StatusNotFound)
		return nil, err
	}
	err = rows.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		slog.Error(err.Error())
		return nil, err
	}
	return &user, nil
}
