package advertisementcreation

import (
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"os"
	"testing"

	base "api-tests-template/tests"

	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя с кредами из переменных окружения и получение его параметров")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}
