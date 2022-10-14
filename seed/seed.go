package seed

import (
	"api-fiber-gorm/model"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func readCsv(filepath string) ([][]string, error) {
	//読み込むCSVファイルを記載
	csvFile, err := os.Open(filepath)
	if err != nil {
		fmt.Println("csvは読みこめませんでした")
		return nil, err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1
	record, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return record, err
}

func UserSeed(db *gorm.DB, filename string) error {
	record, _ := readCsv(filename)
	for i := 1; i < len(record); i++ {
		password := record[i][2]
		bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)

		users := model.User{Username: record[i][0], Email: record[i][1], Password: string(bytes)}

		if err := db.Create(&users).Error; err != nil {
			fmt.Printf("%+v", err)
		}
	}
	return nil
}

func TweetSeed(db *gorm.DB, filename string) error {
	record, _ := readCsv(filename)
	for i := 1; i < len(record); i++ {

		userId, _ := strconv.Atoi(record[i][1])
		tweets := model.Tweet{Content: record[i][0], UserID: userId}

		if err := db.Create(&tweets).Error; err != nil {
			fmt.Printf("%+v", err)
		}
	}
	return nil
}

func FolloweesSeed(db *gorm.DB, filename string) error {
	record, _ := readCsv(filename)
	for i := 1; i < len(record); i++ {

		userId, _ := strconv.Atoi(record[i][0])
		followeeId, _ := strconv.Atoi(record[i][1])

		var user model.User
		var followee model.User
		db.First(&user, userId)
		db.First(&followee, followeeId)

		if err := db.Model(&user).Association("Followees").Append(&followee); err != nil {
			fmt.Printf("%+v", err)
		}
	}
	return nil
}
