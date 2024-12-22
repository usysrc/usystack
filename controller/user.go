package controller

import (
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/usysrc/usystack/model"
)

var sessionStore *session.Store

func CreateSessionStore() {
	sessionStore = session.New()
}

func Login(c *fiber.Ctx) error {
	var loginData model.LoginData
	if err := c.BodyParser(&loginData); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		slog.Error(err.Error())
		return err
	}
	user, err := model.GetUserByName(c, loginData)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	if user.Username == loginData.Username && user.Password == loginData.Password {
		sess, err := sessionStore.Get(c)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		sess.Set("userID", user.ID)
		sess.Save()
		return c.Render("login", fiber.Map{
			"User": user,
		})
	}
	return nil
}

func Logout(c *fiber.Ctx) error {
	sess, err := sessionStore.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve session"})
	}
	sess.Destroy()
	return c.Render("loginform", fiber.Map{})
}
