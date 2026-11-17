package favorites

import (
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"os"
	"testing"

	base "api-tests-template/tests"

	"github.com/stretchr/testify/suite"
)

const (
	favoritesPage     = 1
	favoritesPageSize = 50
)

type TestSuite struct {
	suite.Suite
	loginData  models.LoginOkResponse
	loginData2 models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователей с кредами из переменных окружения и получение их параметров")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
	s.loginData2 = auth.Login(s.T(), os.Getenv("TEST_LOGIN_2"), os.Getenv("TEST_PASSWORD_2"))
}
