package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"rentincome/model"
	"rentincome/tests/testsuite"

	"github.com/stretchr/testify/suite"
)

type LocalTestSuite struct { testsuite.BaseTestSuite }

func TestSuiteRunner(t *testing.T) {
    suite.Run(t, new(LocalTestSuite))
}

func (s *LocalTestSuite) TestRegister_Success() {
    input := map[string]any{
        "username": "newuser",
        "email":    "new@example.com",
        "phone":    "123456789",
        "password": "password123",
    }
    req := s.CreateRequestJson("POST", "/register", input)

    resp, err := s.App.Test(req, -1)


    var user model.User
    json.Unmarshal(resp.Body, &user)

    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        t.Fatalf("Failed to read response body: %v", err)
    }    
    
    expectedBody := `{"message":"Hello, World!"}`
    assert.JSONEq(t, expectedBody, string(body))

    s.Equal(http.StatusCreated, resp.StatusCode )
    s.Equal("newuser", user.Username)
}

func (s *LocalTestSuite) TestLogin_Success() {
    input := map[string]any{
        "username": "testuser",
        "password": "testpass",
    }

    req := s.CreateRequestJson("POST", "/login", input)
    rr := httptest.NewRecorder()

    s.Router.ServeHTTP(rr, req)

    var response map[string]string
    json.Unmarshal(rr.Body.Bytes(), &response)

    s.Equal(http.StatusOK, rr.Code)
    s.Contains(response, "token")
}
