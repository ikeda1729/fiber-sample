package handler

import (
	"api-fiber-gorm/database"
	"api-fiber-gorm/model"
	"api-fiber-gorm/utils"
	"sort"

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
	var user model.User
	userId := int(claims["user_id"].(float64))
	db := database.DB
	db.Preload("Tweets").Preload("Followees.Tweets").Find(&user, userId)

	type Timeline struct {
		Tweet    model.Tweet
		Username string
	}

	var tl Timeline
	var result []Timeline
	// 自分のtweetとusernameを取り出す
	for _, tweet := range user.Tweets {
		tl.Tweet = tweet
		tl.Username = user.Username
		result = append(result, tl)
	}
	// followeesのtweetとusernameを取り出す
	for _, followee := range user.Followees {
		for _, tweet := range followee.Tweets {
			tl.Tweet = tweet
			tl.Username = followee.Username
			result = append(result, tl)
		}
	}
	// timeline全体をsort
	sort.Slice(result[:], func(i, j int) bool {
		return result[i].Tweet.CreatedAt.After(result[j].Tweet.CreatedAt)
	})

	return c.JSON(fiber.Map{"status": "success", "message": "Tweet found", "data": result})
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
