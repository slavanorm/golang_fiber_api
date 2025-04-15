// internal/handlers/estate_test.go
package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"rentincome/internal/models"
	"rentincome/tests/testsuite"

	"github.com/stretchr/testify/assert"
)


type TestSuite struct {
    testsuite.BaseTestSuite
}


// Table-driven test for input validation
func (s *TestSuite) TestCreateEstateValidation() {
    tests := []struct {
        name         string
        requestBody  string
        expectedCode int
        expectedMsg  string
    }{
        {
            "Missing Name",
            `{"address": "123 St", "price": 250000}`,
            http.StatusBadRequest,
            "Name is required",
        },
        {
            "Empty Address",
            `{"name": "Estate", "address": "", "price": 250000}`,
            http.StatusBadRequest,
            "Address is required",
        },
        {
            "Invalid Price",
            `{"name": "Estate", "address": "123 St", "price": -100}`,
            http.StatusBadRequest,
            "Price must be a positive number",
        },
    }

    for _, tt := range tests {
        s.T().Run(tt.name, func(t *testing.T) {
            req := s.CreateRequestString("POST", "/estates", tt.requestBody)
            rr := httptest.NewRecorder()
            s.Router.ServeHTTP(rr, req)
            
            assert.Equal(t, tt.expectedCode, rr.Code)
            assert.Contains(t, rr.Body.String(), tt.expectedMsg)
        })
    }
}

// CRUD test using test suite
func (s *TestSuite) TestEstateCrudOperations() {
    // Create estate
    createReq := []byte(`{
        "name": "Test Estate",
        "address": "123 Test St",
        "price": 250000
    }`)
    
    createRR := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/estates", bytes.NewBuffer(createReq))
    req.Header.Set("Authorization", "Bearer "+s.AuthToken)
    s.Router.ServeHTTP(createRR, req)
    
    assert.Equal(s.T(), http.StatusCreated, createRR.Code)
    
    var createdEstate models.Estate
    json.Unmarshal(createRR.Body.Bytes(), &createdEstate)
    
    // Update estate
    updateReq := []byte(`{"name": "Updated Estate"}`)
    updateRR := httptest.NewRecorder()
    req, _ = http.NewRequest("PUT", fmt.Sprintf("/estates/%d", createdEstate.ID), bytes.NewBuffer(updateReq))
    req.Header.Set("Authorization", "Bearer "+s.AuthToken)
    s.Router.ServeHTTP(updateRR, req)
    
    assert.Equal(s.T(), http.StatusOK, updateRR.Code)
    
    // Delete estate
    deleteRR := httptest.NewRecorder()
    req, _ = http.NewRequest("DELETE", fmt.Sprintf("/estates/%d", createdEstate.ID), nil)
    req.Header.Set("Authorization", "Bearer "+s.AuthToken)
    s.Router.ServeHTTP(deleteRR, req)
    
    assert.Equal(s.T(), http.StatusNoContent, deleteRR.Code)
}


func (s *TestSuite) TestGetEstate() {
    // Create test estate
    estate := models.Estate{
        Name:    "Test Property",
        Type:    "flat_yearly",
        Price:   2500,
        Address: "123 Main St",
        UserID:  s.TestUser.ID,
    }
    s.Transaction.Create(&estate)

    // Create request
    req := s.CreateRequestString("GET", "/estates/"+strconv.Itoa(int(estate.ID)), "")
    rr := httptest.NewRecorder()
    s.Router.ServeHTTP(rr, req)

    // Verify response
    s.Equal(http.StatusOK, rr.Code)
    var response models.Estate
    json.Unmarshal(rr.Body.Bytes(), &response)
    s.Equal(estate.Name, response.Name)
    s.Equal(estate.Price, response.Price)
}

func (s *TestSuite) TestGetEstateNotFound() {
    // Request non-existent estate
    req := s.CreateRequestString("GET", "/estates/999", "")
    rr := httptest.NewRecorder()
    s.Router.ServeHTTP(rr, req)
    s.Equal(http.StatusNotFound, rr.Code)
}

func (s *TestSuite) TestGetEstateUnauthorized() {
    // Create another user
    otherUser := models.User{
        Username: "otheruser",
        Email:    "other@example.com",
        Password: "password123",
    }
    s.Transaction.Create(&otherUser)

    // Create estate for other user
    otherEstate := models.Estate{
        Name:    "Other Property",
        Type:    "flat_yearly",
        Price:   3000,
        Address: "456 Oak St",
        UserID:  otherUser.ID,
    }
    s.Transaction.Create(&otherEstate)

    // Try to access as primary user
    req := s.CreateRequestString("GET", "/estates/"+strconv.Itoa(int(otherEstate.ID)), "")
    rr := httptest.NewRecorder()
    s.Router.ServeHTTP(rr, req)
    s.Equal(http.StatusNotFound, rr.Code)
}