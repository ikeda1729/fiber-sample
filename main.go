package main

import (
	"api-fiber-gorm/database"
	"api-fiber-gorm/router"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	// CORSの設定
	app.Use(cors.New(cors.Config{
		// 認証にcookieなどの情報を必要とするかどうか
		AllowCredentials: true,
	}))

	database.ConnectDB()

	router.SetupRoutes(app)
	port := os.Getenv("PORT")
	log.Fatal(app.Listen(":" + port))
}
