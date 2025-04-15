package database_test

import (
    "testing"

    "rentincome/tests/testsuite"
	"github.com/stretchr/testify/suite"
    "rentincome/model"
)

type LocalTestSuite struct { testsuite.BaseTestSuite }


func TestSuiteRunner(t *testing.T) {
    suite.Run(t, new(LocalTestSuite))}


func (s *LocalTestSuite) TestUserMigration() {
    user := model.User{
        Username: "migrationtest",
        Email:    "migration@test.com",
    }
    result := s.DB.Create(&user)
    s.NoError(result.Error)
    s.NotZero(user.ID)
}
