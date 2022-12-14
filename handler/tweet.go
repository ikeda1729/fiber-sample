package handler

import (
	"api-fiber-gorm/database"
	"api-fiber-gorm/model"
	"api-fiber-gorm/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
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

// GetTweet query tweet
func GetUserTweet(c *fiber.Ctx) error {
	userId := c.Params("userId")
	db := database.DB
	var user model.User
	db.Preload("Tweets", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC").Scopes(Paginate(c))
	}).Find(&user, userId)

	var count int64
	db.Model(&model.Tweet{}).Where("user_id = ?", userId).Count(&count)

	return c.JSON(fiber.Map{"status": "success", "message": "Tweet found", "totalCount": count, "data": user})
}

// GetTweet query tweet
func GetTimeline(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))
	db := database.DB
	type Result struct {
		model.Tweet
		Username string `json:"username"`
	}
	var results []Result
	db.Model(&model.Tweet{}).Distinct("tweets.id").Select("tweets.*, users.username").
		Joins("left join user_followees on user_followees.followee_id = tweets.user_id").
		Joins("left join users on users.id = tweets.user_id").
		Where("user_followees.user_id = ?", userId).
		Or("tweets.user_id = ?", userId).
		Order("tweets.created_at desc").Scopes(Paginate(c)).Scan(&results)

	var count int64
	db.Model(&model.Tweet{}).Distinct("tweets.id").
		Joins("left join user_followees on user_followees.followee_id = tweets.user_id").
		Where("user_followees.user_id = ?", userId).Or("tweets.user_id = ?", userId).Count(&count)

	return c.JSON(fiber.Map{"status": "success", "message": "Tweet found", "totalCount": count, "data": results})
}

// CreateTweet new tweet
func CreateTweet(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))
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

	tweet := model.Tweet{Content: newTweet.Content, UserID: userId}
	db.Create(&tweet)
	return c.JSON(fiber.Map{"status": "success", "message": "Created tweet", "data": tweet})
}

// DeleteTweet delete tweet
func DeleteTweet(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))

	var tweet model.Tweet
	db.First(&tweet, id)
	if tweet.Content == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No tweet found with ID", "data": nil})
	}
	if tweet.UserID != userId {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Unauthorized", "data": nil})
	}
	db.Delete(&tweet)
	return c.JSON(fiber.Map{"status": "success", "message": "Tweet successfully deleted", "data": nil})
}
