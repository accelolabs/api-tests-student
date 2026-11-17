package advertisementcreation

import (
	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"net/http"
	"os"
	"strings"
	"testing"

	base "api-tests-template/tests"
	"api-tests-template/tests/helpers"

	"github.com/stretchr/testify/require"
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

// Positive: Проверяем создание объявления с валидными полями
func (s *TestSuite) TestCreate() {
	request := helpers.NewAdvertisementRequest()

	t := s.T()
	token := s.loginData.Token

	created := advertisements.PostAdvertisement(t, token, request)
	require.NotEmpty(t, created.ID)

	t.Cleanup(func() {
		advertisements.DeleteAdvertisement(t, token, created.ID)
	})

	s.Run("Объявление доступно по ID и содержит переданные данные", func() {
		t := s.T()

		actual := advertisements.GetAdvertisement(t, token, created.ID)

		require.Equal(t, created.ID, actual.ID)
		require.Equal(t, request.Title, actual.Title)
		require.Equal(t, request.Description, actual.Description)
		require.InDelta(t, request.Price, actual.Price, 0.001)
		require.Equal(t, request.Quantity, actual.Quantity)
		require.Equal(t, s.loginData.User.Id, actual.UserID)
	})

	s.Run("Объявление находится по полному названию", func() {
		t := s.T()

		list := advertisements.GetAdvertisementsList(t, token, request.Title)

		found := false
		for _, item := range list.Items {
			if item.ID == created.ID {
				found = true
				break
			}
		}

		require.True(t, found, "Созданное объявление не найдено в поиске")
	})

	s.Run("Фотографии доступны", func() {
		t := s.T()

		photos := advertisements.GetAdvertisementPhotos(t, token, created.ID)

		require.Len(t, photos, len(request.Photos))
		require.Len(t, created.Photos, len(request.Photos))

		for _, photo := range photos {
			require.NotEmpty(t, photo.ID)
			require.NotEmpty(t, photo.URL)
		}

		require.ElementsMatch(t, created.Photos, photos)
	})
}

// Negative: Проверяем 400 ошибку при создании объявления с пустым названием
func (s *TestSuite) TestCreateTitleEmpty() {
	request := helpers.NewAdvertisementRequest()
	request.Title = ""

	t := s.T()
	token := s.loginData.Token

	resp := advertisements.PostAdvertisementExpectError(t, token, request, http.StatusBadRequest)

	s.Run("Объявление без названия не должно быть создано", func() {
		t := s.T()

		assertions.APIErrorEquals(t, resp, helpers.InvalidRequestError, "title is required")
	})
}

// Negative: Проверяем 400 ошибку при создании объявления со слишком длинным названием
func (s *TestSuite) TestCreateTitleTooLong() {
	request := helpers.NewAdvertisementRequest()
	request.Title = strings.Repeat("Лемур", 11)

	t := s.T()
	token := s.loginData.Token

	resp := advertisements.PostAdvertisementExpectError(t, token, request, http.StatusBadRequest)

	s.Run("Объявление со слишком длинным названием не должно быть создано", func() {
		t := s.T()

		assertions.APIErrorEquals(t, resp, helpers.InvalidRequestError, "title must not exceed 50 characters provided 55")
	})
}

// Negative: Проверяем 400 ошибку при создании объявления с отрицательной ценой
func (s *TestSuite) TestCreatePriceNegative() {
	request := helpers.NewAdvertisementRequest()
	request.Price = -request.Price

	t := s.T()
	token := s.loginData.Token

	resp := advertisements.PostAdvertisementExpectError(t, token, request, http.StatusBadRequest)

	s.Run("Объявление с отрицательной ценой не должно быть создано", func() {
		t := s.T()

		assertions.APIErrorEquals(t, resp, helpers.InvalidRequestError, "price must be between 0 and 1000000 inclusive")
	})
}

// Negative: Проверяем 401 ошибку при попытке создать объявление без токена
func (s *TestSuite) TestCreateTokenEmpty() {
	request := helpers.NewAdvertisementRequest()

	t := s.T()
	resp := advertisements.PostAdvertisementExpectError(t, "", request, http.StatusUnauthorized)

	s.Run("Объявление нельзя создать без токена", func() {
		t := s.T()

		assertions.APIErrorEquals(t, resp, helpers.UnauthorizedError, helpers.AuthorizationRequiredMessage)
	})
}

// Negative: Проверяем 401 ошибку при попытке создать объявление с невалидным токеном
func (s *TestSuite) TestCreateTokenInvalid() {
	request := helpers.NewAdvertisementRequest()

	t := s.T()
	resp := advertisements.PostAdvertisementExpectError(t, "invalidtoken", request, http.StatusUnauthorized)

	s.Run("Объявление нельзя создать с неверным токеном", func() {
		t := s.T()

		assertions.APIErrorEquals(t, resp, helpers.UnauthorizedError, "Invalid or expired token")
	})
}
