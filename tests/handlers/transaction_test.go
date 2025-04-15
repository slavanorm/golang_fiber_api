package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"rentincome/internal/middleware"
	"rentincome/internal/models"
	"rentincome/tests/testsuite"
)

type TransactionTestSuite struct {
    testsuite.BaseTestSuite
}

func (s *TransactionTestSuite) TestCreateTransaction() {
    // Create estate for test user
    estate := models.Estate{
        Type:    "flat_yearly",
        Price:   1000,
        Address: "Test Address",
        UserID:  s.TestUser.ID,
    }
    s.Transaction.Create(&estate)

    // Create transaction payload
    payload := map[string]any{
        "estate_id":   estate.ID,
        "amount":      500.50,
        "date_actual": time.Now().Format(time.RFC3339),
    }

    // Create and serve request
    req := s.CreateRequestJson("POST", "/transactions", payload)
    rr := httptest.NewRecorder()
    s.Router.ServeHTTP(rr, req)

    // Verify response
    s.Equal(http.StatusCreated, rr.Code)
    var response map[string]interface{}
    json.Unmarshal(rr.Body.Bytes(), &response)
    s.Equal(500.50, response["amount"])

    // Verify database entry
    var transaction models.Transaction
    s.Transaction.First(&transaction, uint(response["id"].(float64)))
    s.Equal(estate.ID, transaction.EstateID)
    s.Equal(s.TestUser.ID, transaction.UserID)
}

func (s *TransactionTestSuite) TestUpdateTransaction() {
    // Create test transaction
    estate := models.Estate{
        Type:    "flat_yearly",
        Price:   1000,
        Address: "Test Address",
        UserID:  s.TestUser.ID,
    }
    s.Transaction.Create(&estate)
    transaction := models.Transaction{
        UserID:      s.TestUser.ID,
        EstateID:    estate.ID,
        Amount:      500.50,
        DateActual:  time.Now(),
    }
    s.Transaction.Create(&transaction)

    // Update payload
    updatePayload := map[string]interface{}{
        "amount":   600.75,
        "comment1": "Updated comment",
    }

    // Create and serve request
    req := s.CreateRequestJson("PUT", "/transactions/"+strconv.Itoa(int(transaction.ID)), updatePayload)
    rr := httptest.NewRecorder()
    s.Router.ServeHTTP(rr, req)

    // Verify response and database update
    s.Equal(http.StatusOK, rr.Code)
    var updatedTransaction models.Transaction
    s.Transaction.First(&updatedTransaction, transaction.ID)
    s.Equal(600.75, updatedTransaction.Amount)
    s.Equal("Updated comment", updatedTransaction.Comment1)
}

func (s *TransactionTestSuite) TestDeleteTransaction() {
    // Create test transaction
    estate := models.Estate{
        Type:    "flat_yearly",
        Price:   1000,
        Address: "Test Address",
        UserID:  s.TestUser.ID,
    }
    s.Transaction.Create(&estate)
    transaction := models.Transaction{
        UserID:      s.TestUser.ID,
        EstateID:    estate.ID,
        Amount:      500.50,
        DateActual:  time.Now(),
    }
    s.Transaction.Create(&transaction)

    // Create and serve DELETE request
    req := s.CreateRequestString("DELETE", "/transactions/"+strconv.Itoa(int(transaction.ID)), "")
    rr := httptest.NewRecorder()
    s.Router.ServeHTTP(rr, req)

    // Verify deletion
    s.Equal(http.StatusNoContent, rr.Code)
    var deletedTransaction models.Transaction
    result := s.Transaction.First(&deletedTransaction, transaction.ID)
    s.Error(result.Error)
    s.Equal("record not found", result.Error.Error())
}

func (s *TransactionTestSuite) TestUnauthorizedAccess() {
    // Create test transaction for primary user
    estate := models.Estate{
        Type:    "flat_yearly",
        Price:   1000,
        Address: "Test Address",
        UserID:  s.TestUser.ID,
    }
    s.Transaction.Create(&estate)
    transaction := models.Transaction{
        UserID:      s.TestUser.ID,
        EstateID:    estate.ID,
        Amount:      500.50,
        DateActual:  time.Now(),
    }
    s.Transaction.Create(&transaction)

    // Create another user and generate token
    anotherUser := models.User{
        Username: "hacker",
        Email:    "hacker@example.com",
        Phone:    "1122334455",
        Password: "hackpass",
    }
    s.Transaction.Create(&anotherUser)
    unauthToken := middleware.GenerateAuthToken(anotherUser.ID)

    // Create unauthorized DELETE request
    req := httptest.NewRequest("DELETE", "/transactions/"+strconv.Itoa(int(transaction.ID)), nil)
    req.Header.Set("Authorization", "Bearer "+unauthToken)
    req = req.WithContext(context.WithValue(req.Context(), "user", anotherUser))
    rr := httptest.NewRecorder()
    s.Router.ServeHTTP(rr, req)

    // Verify access denied and transaction remains
    s.Equal(http.StatusForbidden, rr.Code)
    var existingTransaction models.Transaction
    s.Transaction.First(&existingTransaction, transaction.ID)
    s.Equal(transaction.ID, existingTransaction.ID)
}

func TestTransactionSuite(t *testing.T) {
    suite.Run(t, new(TransactionTestSuite))
}