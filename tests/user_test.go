package tests

import (
	"api-fiber-gorm/database"
	"api-fiber-gorm/model"
	"api-fiber-gorm/router"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	testCase string
	url      string
	method   string
	multiple bool
	expected model.User
}

type ResponseBody struct {
	Status  string
	Message string
	Data    model.User
}

type ResponseBodySlice struct {
	Status  string
	Message string
	Data    []model.User
}

type LoginResponse struct {
	Status  string
	Message string
	Data    struct {
		Username string
		Token    string
	}
}

func TestUser(t *testing.T) {
	app := fiber.New()
	app.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path} cookie ${reqHeaders}\n",
	}))
	router.SetupRoutes(app)

	if err := godotenv.Load("./.env.test"); err != nil {
		panic(err)
	}
	// データベースを初期化
	database.InitTestDB()
	db, _ := database.DB.DB()
	defer db.Close()

	expectBody := model.User{
		Username: "Alice",
		Email:    "alice@ex.com",
		Password: "password",
	}

	jsonBody, _ := json.Marshal(&expectBody)

	updateBody := model.User{
		Username: "AliceUpdated",
		Email:    "alice@ex.com",
		Password: "password",
	}
	_ = updateBody

	testCases := []TestCase{
		{"POST - User", "/api/user", http.MethodPost, false, expectBody},
		{"GET - User", "/api/user/1", http.MethodGet, false, expectBody},
		{"GET - Users", "/api/user", http.MethodGet, true, expectBody},
		{"POST - Login", "/api/auth/login", http.MethodPost, false, expectBody},
		// {"UPDATE - User", "/api/user/1", http.MethodPatch, false, updateBody}, //認証できず
		{"DELETE - User", "/api/user/1", http.MethodDelete, false, expectBody},
		{"GET - Deleted User", "/api/v1/user/1", http.MethodGet, false, model.User{}},
	}

	var token string

	for _, tc := range testCases {
		t.Run(tc.testCase, func(t *testing.T) {
			readerBody := bytes.NewBuffer(jsonBody)
			// updateするときはupdateBodyを送る
			if tc.method == "PATCH" {
				jsonBody, _ = json.Marshal(&updateBody)
				readerBody = bytes.NewBuffer(jsonBody)
			}

			req, _ := http.NewRequest(tc.method, tc.url, readerBody)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := app.Test(req)

			if tc.method == "DELETE" {
				// DELETEの場合はチェックを飛ばして、GETで確かめる
				return
			}

			// 全件取得の時はスライス
			if tc.multiple {
				var actual ResponseBodySlice
				json.NewDecoder(resp.Body).Decode(&actual)
				fmt.Print(actual.Data)
				checkError(t, err, 200, resp.StatusCode, tc.expected, actual.Data[0])
			} else if tc.testCase == "GET - Deleted User" {
				var actual ResponseBody
				json.NewDecoder(resp.Body).Decode(&actual)
				fmt.Print(actual)
				checkError(t, err, 404, resp.StatusCode, tc.expected, actual.Data)
			} else if tc.testCase == "POST - Login" {
				var actual LoginResponse
				json.NewDecoder(resp.Body).Decode(&actual)
				actualUser := model.User{
					Username: actual.Data.Username,
				}
				checkError(t, err, 200, resp.StatusCode, tc.expected, actualUser)
				token = actual.Data.Token
			} else {
				var actual ResponseBody
				json.NewDecoder(resp.Body).Decode(&actual)
				fmt.Print(actual)
				checkError(t, err, 200, resp.StatusCode, tc.expected, actual.Data)
			}
		})
	}
}

func checkError(
	t *testing.T,
	err error,
	expectedStatus int,
	actualStatus int,
	expected, actual model.User) {
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, actualStatus)
	assert.Equal(t, expected.Username, actual.Username)
}
