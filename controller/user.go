package controller

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/usysrc/usystack/model"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash verifies a password against a hashed password
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// the register page
func Register(c *fiber.Ctx) error {
	return c.Render("register", fiber.Map{}, "layout")
}

func RegisterUser(c *fiber.Ctx) error {
	var registerData model.RegisterData
	if err := c.BodyParser(&registerData); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		slog.Error(err.Error())
		return err
	}

	hashedPassword, err := hashPassword(registerData.Password)
	if err != nil {
		return c.Render("registerform", fiber.Map{})
	}
	registerData.Password = hashedPassword

	err = model.RegisterUser(c, registerData)
	if err != nil {
		return c.Render("registerform", fiber.Map{
			"ErrorDescription": "Could not register user: username already in use.",
		})
	}

	return c.Render("registerform", fiber.Map{})
}

// the login page
func Login(c *fiber.Ctx) error {
	sess, ok := c.Locals("session").(*session.Session)
	if !ok {
		return c.Render("loginform", fiber.Map{})
	}
	userID := sess.Get("userID")
	user := &model.User{}
	if userID != nil {
		id, err := strconv.Atoi(userID.(string))
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		user, err = model.GetUserByID(c, id)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		user.LoggedIn = true
	}

	return c.Render("login", fiber.Map{
		"User": user,
	}, "layout")
}

// handle login of user
func LoginUser(c *fiber.Ctx) error {
	var loginData model.LoginData
	if err := c.BodyParser(&loginData); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		slog.Error(err.Error())
		return err
	}

	user, err := model.GetUserByName(c, loginData.Username)
	if err != nil {
		slog.Error(err.Error())
		return c.Render("loginform", fiber.Map{
			"LoginFailed": true,
		})
	}

	if user.Username != loginData.Username || !checkPasswordHash(loginData.Password, user.Password) {
		return c.Render("loginform", fiber.Map{
			"LoginFailed": true,
		})
	}

	sess, ok := c.Locals("session").(*session.Session)
	if !ok {
		return c.Render("loginform", fiber.Map{})
	}
	sess.Set("userID", user.ID)
	sess.Save()
	user.LoggedIn = true

	return c.Render("loginform", fiber.Map{
		"User": user,
	})
}

func Logout(c *fiber.Ctx) error {
	if sess, ok := c.Locals("session").(*session.Session); ok {
		sess.Destroy()
	}
	return c.Render("loginform", fiber.Map{})
}
