package router

import (
	"api-fiber-gorm/handler"
	"api-fiber-gorm/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// Middleware
	api := app.Group("/api", logger.New())
	api.Get("/", handler.Hello)

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/logout", handler.Logout)

	// User
	user := api.Group("/user")
	user.Get("/:id", handler.GetUser)
	user.Get("/", handler.GetUsers)
	user.Get("/:userId/tweet", handler.GetUserTweet)
	user.Get("/:userId/isFollowing", middleware.Protected(), handler.GetIsFollowing)
	user.Get("/:userId/followings", handler.GetUserFollowing)
	user.Get("/:userId/followers", handler.GetUserFollowers)
	user.Post("/", handler.CreateUser)
	user.Patch("/:id", middleware.Protected(), handler.UpdateUser)
	user.Delete("/:id", middleware.Protected(), handler.DeleteUser)

	// Tweet
	tweet := api.Group("/tweet")
	tweet.Get("/", handler.GetAllTweets)
	tweet.Get("/:id", handler.GetTweet)
	tweet.Post("/", middleware.Protected(), handler.CreateTweet)
	tweet.Delete("/:id", middleware.Protected(), handler.DeleteTweet)

	// Relation
	relation := api.Group("/relation")
	relation.Post("/:followeeId", middleware.Protected(), handler.CreateRelation)
	relation.Delete("/:followeeId", middleware.Protected(), handler.DeleteRelation)
}
