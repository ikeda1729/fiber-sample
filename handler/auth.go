package handler

import (
	"api-fiber-gorm/config"
	"api-fiber-gorm/database"
	"api-fiber-gorm/model"
	"api-fiber-gorm/utils"
	"errors"
	"net/mail"
	"os"
	"time"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// CheckPasswordHash compare password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUserByEmail(e string) (*model.User, error) {
	db := database.DB
	var user model.User
	if err := db.Where("email = ?", e).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func getUserByUsername(u string) (*model.User, error) {
	db := database.DB
	var user model.User
	if err := db.Where("username = ?", u).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Login get user and password
func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	type UserData struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	input := new(LoginInput)

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
	}

	// validation
	if err := utils.ValidateStruct(*input); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "validation error", "message": "Review your input", "data": err})
	}

	email := input.Email
	pass := input.Password
	user, err := getUserByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Error on email", "data": err})
	}

	var ud UserData
	if user != nil {
		ud = UserData{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Password: user.Password,
		}
	}

	if !CheckPasswordHash(pass, ud.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid password", "data": nil})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = ud.Username
	claims["user_id"] = ud.ID
	exp := time.Now().Add(time.Hour * 72)
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Cookieに保存
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    t,
		Expires:  exp,
		SameSite: "none",
		Secure:   true,
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	type LoginResponse struct {
		UserID   uint
		Username string
		Token    string
	}

	data := LoginResponse{UserID: ud.ID, Username: ud.Username, Token: t}

	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": data})
}

func Logout(ctx *fiber.Ctx) error {
	// cookieをクリアする
	cookie := fiber.Cookie{
		Name:    "jwt",
		Value:   "",
		Expires: time.Now().Add(-time.Hour * 24), // -を指定
		Domain:  os.Getenv("DOMAIN"),
	}

	ctx.Cookie(&cookie)
	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "Success logout",
	})
}
