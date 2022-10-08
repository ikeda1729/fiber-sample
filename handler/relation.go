package handler

import (
	"api-fiber-gorm/database"
	"api-fiber-gorm/model"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type RelationResponse struct {
	UserID     uint
	FolloweeID uint
}

func GetUserFollowing(c *fiber.Ctx) error {
	userId := c.Params("userId")
	db := database.DB
	var user model.User
	db.First(&user, userId)
	var users []model.UserResponse
	db.Model(&user).Association("Followees").Find(&users)

	return c.JSON(fiber.Map{"status": "success", "message": "Following", "data": users})
}

func GetUserFollowers(c *fiber.Ctx) error {
	userId := c.Params("userId")
	db := database.DB
	var users []model.UserResponse
	db.Model(&model.User{}).Joins("inner join user_followees on user_followees.user_id = users.id").
		Where("user_followees.followee_id = ?", userId).Find(&users)

	return c.JSON(fiber.Map{"status": "success", "message": "Following", "data": users})
}

func CreateRelation(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))
	followeeId := c.Params("followeeId")
	db := database.DB
	var user model.User
	var followee model.User
	db.First(&user, userId)
	db.First(&followee, followeeId)

	db.Model(&user).Association("Followees").Append(&followee)

	data := RelationResponse{UserID: user.ID, FolloweeID: followee.ID}
	return c.JSON(fiber.Map{"status": "success", "message": "Followed", "data": data})
}

func DeleteRelation(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))
	followeeId := c.Params("followeeId")
	db := database.DB
	var user model.User
	var followee model.User
	db.First(&user, userId)
	db.First(&followee, followeeId)

	db.Model(&user).Association("Followees").Delete(&followee)

	data := RelationResponse{UserID: user.ID, FolloweeID: followee.ID}
	return c.JSON(fiber.Map{"status": "success", "message": "Unfollowed", "data": data})
}
