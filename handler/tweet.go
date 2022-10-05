package handler

import (
	"api-fiber-gorm/database"
	"api-fiber-gorm/model"

	"github.com/gofiber/fiber/v2"
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
	if tweet.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No tweet found with ID", "data": nil})

	}
	return c.JSON(fiber.Map{"status": "success", "message": "Tweet found", "data": tweet})
}

// CreateTweet new tweet
func CreateTweet(c *fiber.Ctx) error {
	db := database.DB
	tweet := new(model.Tweet)
	if err := c.BodyParser(tweet); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create tweet", "data": err})
	}
	db.Create(&tweet)
	return c.JSON(fiber.Map{"status": "success", "message": "Created tweet", "data": tweet})
}

// DeleteTweet delete tweet
func DeleteTweet(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var tweet model.Tweet
	db.First(&tweet, id)
	if tweet.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No tweet found with ID", "data": nil})

	}
	db.Delete(&tweet)
	return c.JSON(fiber.Map{"status": "success", "message": "Tweet successfully deleted", "data": nil})
}
