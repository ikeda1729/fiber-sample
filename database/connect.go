package database

import (
	"api-fiber-gorm/model"
	"fmt"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDB connect to db
func ConnectDB() {
	var err error
	p := os.Getenv("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", os.Getenv("DB_HOST"), port, os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	print("dsnForDebug: ", dsn)

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")
	// DB.Migrator().DropTable(&model.Tweet{}, &model.User{}, "user_followees")
	DB.AutoMigrate(&model.Tweet{}, &model.User{})
	fmt.Println("Database Migrated")
	// if err := seed.UserSeed(DB, "./seed/users.csv"); err != nil {
	// 	fmt.Println(err)
	// }
	// if err := seed.TweetSeed(DB, "./seed/tweets.csv"); err != nil {
	// 	fmt.Println(err)
	// }
	// if err := seed.FolloweesSeed(DB, "./seed/user_followees.csv"); err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("Seeds added")
}

// ConnectDB connect to db
func InitTestDB() {
	var err error
	p := os.Getenv("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), port, os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")
	DB.Migrator().DropTable(&model.Tweet{}, &model.User{})
	DB.AutoMigrate(&model.Tweet{}, &model.User{})
	fmt.Println("Database Migrated")
}
