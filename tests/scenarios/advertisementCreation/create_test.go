package advertisementcreation

import (
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	advModels "api-tests-template/internal/managers/advertisements/models"
	base "api-tests-template/tests"

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
	request := advModels.NewAdvertisement{
		Title:       "Лемуры " + time.Now().Format("2006-01-02 15:04:05.000000000"),
		Description: "Сидят на ветке и что-то замышляют на будущее.",
		Price:       1337.67,
		Quantity:    2,
		Photos:      []string{"../../testdata/photo.jpg"},
	}

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
	request := advModels.NewAdvertisement{
		Title:       "",
		Description: "Сидят на ветке и что-то замышляют на будущее.",
		Price:       1337.67,
		Quantity:    2,
		Photos:      []string{"../../testdata/photo.jpg"},
	}

	t := s.T()
	token := s.loginData.Token

	resp := advertisements.PostAdvertisementExpectError(t, token, request, http.StatusBadRequest)

	s.Run("Объявление без названия не должно быть создано", func() {
		t := s.T()

		require.Equal(t, "invalid_request", resp.Error)
		require.Equal(t, "title is required", resp.Message)
	})
}

// Negative: Проверяем 400 ошибку при создании объявления с слишком длинным названием
func (s *TestSuite) TestCreateTitleTooLong() {
	request := advModels.NewAdvertisement{
		Title:       strings.Repeat("Лемур", 11),
		Description: "Сидят на ветке и что-то замышляют на будущее.",
		Price:       1337.67,
		Quantity:    2,
		Photos:      []string{"../../testdata/photo.jpg"},
	}

	t := s.T()
	token := s.loginData.Token

	resp := advertisements.PostAdvertisementExpectError(t, token, request, http.StatusBadRequest)

	s.Run("Объявление со слишком длинным названием не должно быть создано", func() {
		t := s.T()

		require.Equal(t, "invalid_request", resp.Error)
		require.Equal(t, "title must not exceed 50 characters provided 55", resp.Message)
	})
}

// Negative: Проверяем 400 ошибку при создании объявления с отрицательной ценой
func (s *TestSuite) TestCreatePriceNegative() {
	request := advModels.NewAdvertisement{
		Title:       "Лемур",
		Description: "Сидят на ветке и что-то замышляют на будущее.",
		Price:       -1337.67,
		Quantity:    2,
		Photos:      []string{"../../testdata/photo.jpg"},
	}

	t := s.T()
	token := s.loginData.Token

	resp := advertisements.PostAdvertisementExpectError(t, token, request, http.StatusBadRequest)

	s.Run("Объявление с отрицательной ценой не должно быть создано", func() {
		t := s.T()

		require.Equal(t, "invalid_request", resp.Error)
		require.Equal(t, "price must be between 0 and 1000000 inclusive", resp.Message)
	})
}

// Negative: Проверяем 401 ошибку при попытке создать объявление без токена
func (s *TestSuite) TestCreateTokenEmpty() {
	request := advModels.NewAdvertisement{
		Title:       "Лемуры " + time.Now().Format("2006-01-02 15:04:05.000000000"),
		Description: "Сидят на ветке и что-то замышляют на будущее.",
		Price:       1337.67,
		Quantity:    2,
		Photos:      []string{"../../testdata/photo.jpg"},
	}

	t := s.T()
	resp := advertisements.PostAdvertisementExpectError(t, "", request, http.StatusUnauthorized)

	s.Run("Объявление нельзя создать без токена", func() {
		t := s.T()

		require.Equal(t, "unauthorized", resp.Error)
		require.Equal(t, "Authorization header required", resp.Message)
	})
}

// Negative: Проверяем 401 ошибку при попытке создать объявление с невалидным токеном
func (s *TestSuite) TestCreateTokenInvalid() {
	request := advModels.NewAdvertisement{
		Title:       "Лемуры " + time.Now().Format("2006-01-02 15:04:05.000000000"),
		Description: "Сидят на ветке и что-то замышляют на будущее.",
		Price:       1337.67,
		Quantity:    2,
		Photos:      []string{"../../testdata/photo.jpg"},
	}

	t := s.T()
	resp := advertisements.PostAdvertisementExpectError(t, "invalidtoken", request, http.StatusUnauthorized)

	s.Run("Объявление нельзя создать с неверным токеном", func() {
		t := s.T()

		require.Equal(t, "unauthorized", resp.Error)
		require.Equal(t, "Invalid or expired token", resp.Message)
	})
}
