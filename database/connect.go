package database

import (
	"api-fiber-gorm/config"
	"api-fiber-gorm/model"
	"api-fiber-gorm/seed"
	"fmt"
	"os"
	"strconv"

	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDB connect to db
func ConnectDB() {
	url := os.Getenv("DATABASE_URL")
	connection, err := pq.ParseURL(url)
	if err != nil {
		panic(err.Error())
	}
	connection += " sslmode=require"
	DB, err = gorm.Open(postgres.Open(connection), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")
	// DB.Migrator().DropTable(&model.Tweet{}, &model.User{}, "user_followees")
	DB.AutoMigrate(&model.Tweet{}, &model.User{})
	fmt.Println("Database Migrated")
	if err := seed.UserSeed(DB, "./seed/users.csv"); err != nil {
		fmt.Println(err)
	}
	if err := seed.TweetSeed(DB, "./seed/tweets.csv"); err != nil {
		fmt.Println(err)
	}
	if err := seed.FolloweesSeed(DB, "./seed/user_followees.csv"); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Seeds added")
}

// ConnectDB connect to db
func InitTestDB() {
	var err error
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Config("DB_HOST"), port, config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_NAME"))
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")
	DB.Migrator().DropTable(&model.Tweet{}, &model.User{})
	DB.AutoMigrate(&model.Tweet{}, &model.User{})
	fmt.Println("Database Migrated")
}
