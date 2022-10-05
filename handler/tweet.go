package handler

import (
	"api-fiber-gorm/database"
	"api-fiber-gorm/model"
	"api-fiber-gorm/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// GetAllTweets query all tweets
func GetAllTweets(c *fiber.Ctx) error {
	db := database.DB
	var tweets []model.Tweet
	db.Find(&tweets)
	return c.JSON(fiber.Map{"status": "success", "message": "All tweets", "data": tweets})
}

// GetTweet query tweet
func GetTweet(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var tweet model.Tweet
	db.Find(&tweet, id)
	if tweet.Content == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No tweet found with ID", "data": nil})

	}
	return c.JSON(fiber.Map{"status": "success", "message": "Tweet found", "data": tweet})
}

// CreateTweet new tweet
func CreateTweet(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	db := database.DB

	type NewTweet struct {
		Content string `json:"content" validate:"required,lte=280"`
	}
	newTweet := new(NewTweet)
	if err := c.BodyParser(newTweet); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	// validation
	if err := utils.ValidateStruct(*newTweet); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "validation error", "message": "Review your input", "data": err})
	}

	var user model.User
	db.Where("username = ?", username).First(&user)

	tweet := model.Tweet{Content: "test", UserID: username, User: user}
	db.Create(&tweet)
	return c.JSON(fiber.Map{"status": "success", "message": "Created tweet", "data": tweet})
}

// DeleteTweet delete tweet
func DeleteTweet(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var tweet model.Tweet
	db.First(&tweet, id)
	if tweet.Content == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No tweet found with ID", "data": nil})

	}
	db.Delete(&tweet)
	return c.JSON(fiber.Map{"status": "success", "message": "Tweet successfully deleted", "data": nil})
}
