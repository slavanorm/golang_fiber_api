package testsuite

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"rentincome/database"
	"rentincome/model"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type BaseTestSuite struct {
	suite.Suite
	App       *fiber.App
	db        *gorm.DB
	DB        *gorm.DB //transaction
	TestUser  model.User
	AuthToken string
}

// Setup test suite
func (s *BaseTestSuite) SetupSuite() {
	s.App = fiber.New()
	database.InitDB(":memory:")

}

func (s *BaseTestSuite) TearDownTest() {
	s.DB.Rollback()
}

func (s *BaseTestSuite) TearDownSuite() {
	database.CloseDB()
}

// Helper function to create test user
func createTestUser(s *BaseTestSuite) (model.User, string) {
	password := "testpass"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	pass := string(hashedPassword)
	user := model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "1234567890",
		Password: pass,
	}

	s.DB.Create(&user)
	return user, pass
}

func (s *BaseTestSuite) SetupTest() {
	s.DB = s.db.Begin()
	s.TestUser, s.AuthToken = createTestUser(s)
}

// Helper methods to create authenticated requests
func (s *BaseTestSuite) CreateRequestString(method, path string, body string) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+s.AuthToken)
	return req.WithContext(context.WithValue(req.Context(), "user", s.TestUser))
}

func (s *BaseTestSuite) CreateRequestJson(method, path string, body map[string]any) *http.Request {
	parsed, _ := json.Marshal(body)
	req := httptest.NewRequest(method, path, bytes.NewBuffer(parsed))
	req.Header.Set("Authorization", "Bearer "+s.AuthToken)
	return req.WithContext(context.WithValue(req.Context(), "user", s.TestUser))
}
