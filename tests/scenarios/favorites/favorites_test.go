package favorites

import (
	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"net/http"
	"os"
	"testing"

	base "api-tests-template/tests"
	"api-tests-template/tests/helpers"

	favManager "api-tests-template/internal/managers/favorites"

	"github.com/stretchr/testify/require"
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

// Positive: Проверяем добавление в избранное и удаление из него.
func (s *TestSuite) TestAddAndDelete() {
	t := s.T()
	token := s.loginData.Token

	created := helpers.CreateTestAdvertisement(t, token)

	t.Cleanup(func() {
		t.Logf("Удаляем объявление %s", created.ID)
		advertisements.DeleteAdvertisement(t, token, created.ID)
	})

	token2 := s.loginData2.Token
	t.Logf("Добавляем %s в избранное", created.ID)
	favManager.PostFavorite(t, token2, created.ID)

	s.Run("Добавленное объявление есть в избранном", func() {
		t := s.T()

		ids := favoriteAdvertisementIDs(t, s.loginData2.Token)

		require.Contains(t, ids, created.ID)
		t.Log("Объявление есть в избранном")
	})

	t.Logf("Удаляем %s из избранного", created.ID)
	favManager.DeleteFavorite(t, token2, created.ID)

	s.Run("Удаленное объявление отсутствует в избранном", func() {
		t := s.T()

		ids := favoriteAdvertisementIDs(t, s.loginData2.Token)

		require.NotContains(t, ids, created.ID)
		t.Log("Объявления нет в избранном")
	})
}

// Positive: Проверяем исчезновение удалённого объявления из избранного.
func (s *TestSuite) TestDeletedAdvertisement() {
	t := s.T()
	token := s.loginData.Token

	created := helpers.CreateTestAdvertisement(t, token)

	token2 := s.loginData2.Token
	t.Logf("Добавляем %s в избранное", created.ID)
	favManager.PostFavorite(t, token2, created.ID)

	s.Run("Добавленное объявление есть в избранном", func() {
		t := s.T()

		ids := favoriteAdvertisementIDs(t, s.loginData2.Token)

		require.Contains(t, ids, created.ID)
		t.Log("Объявление есть в избранном")
	})

	t.Logf("Удаляем объявление %s", created.ID)
	advertisements.DeleteAdvertisement(t, token, created.ID)

	s.Run("Удаленное автором объявление отсутствует в избранном", func() {
		t := s.T()

		ids := favoriteAdvertisementIDs(t, s.loginData2.Token)

		require.NotContains(t, ids, created.ID)
		t.Log("Объявления нет в избранном")
	})
}

// Negative: Проверяем запрет добавления собственного объявления в избранное.
func (s *TestSuite) TestAddOwn() {
	t := s.T()
	token := s.loginData.Token

	created := helpers.CreateTestAdvertisement(t, token)

	t.Cleanup(func() {
		t.Logf("Удаляем объявление %s", created.ID)
		advertisements.DeleteAdvertisement(t, token, created.ID)
	})

	t.Log("Отправляем запрос с ожидаемой ошибкой")
	resp := favManager.PostFavoriteExpectError(t, token, created.ID, http.StatusForbidden)

	s.Run("Собственное объявление нельзя добавить в избранное", func() {
		t := s.T()

		assertions.APIErrorEquals(t, resp, "forbidden", "Cannot add own advertisement to favorites")
	})
}

// Negative: Проверяем запрет повторного добавления в избранное.
func (s *TestSuite) TestRepeatedAdd() {
	t := s.T()
	token := s.loginData.Token

	created := helpers.CreateTestAdvertisement(t, token)

	t.Cleanup(func() {
		t.Logf("Удаляем объявление %s", created.ID)
		advertisements.DeleteAdvertisement(t, token, created.ID)
	})

	token2 := s.loginData2.Token
	t.Logf("Добавляем %s в избранное", created.ID)
	favManager.PostFavorite(t, token2, created.ID)

	t.Cleanup(func() {
		t.Logf("Удаляем %s из избранного", created.ID)
		favManager.DeleteFavorite(t, token2, created.ID)
	})

	t.Log("Отправляем запрос с ожидаемой ошибкой")
	resp := favManager.PostFavoriteExpectError(t, token2, created.ID, http.StatusConflict)

	s.Run("Объявление нельзя добавить в избранное дважды", func() {
		t := s.T()

		assertions.APIErrorEquals(t, resp, "already_exists", "Advertisement is already in favorites")
	})
}

// Negative: Проверяем добавление и удаление несуществующего объявления из избранного.
func (s *TestSuite) TestAddAndDeleteNonexistentID() {
	t := s.T()
	token := s.loginData.Token

	created := helpers.CreateTestAdvertisement(t, token)

	t.Logf("Удаляем объявление %s", created.ID)
	advertisements.DeleteAdvertisement(t, token, created.ID)

	token2 := s.loginData2.Token
	t.Log("Отправляем запрос с ожидаемой ошибкой")
	resp := favManager.PostFavoriteExpectError(t, token2, created.ID, http.StatusNotFound)

	s.Run("Нельзя добавить несуществующее объявление", func() {
		t := s.T()

		assertions.APIErrorEquals(t, resp, "not_found", "Advertisement not found")
	})

	t.Log("Отправляем запрос с ожидаемой ошибкой")
	resp = favManager.DeleteFavoriteExpectError(t, token2, created.ID, http.StatusNotFound)

	s.Run("Нельзя удалить несуществующее объявление", func() {
		t := s.T()

		assertions.APIErrorEquals(t, resp, "not_found", "Favorite entry not found")
	})
}

// Negative: Проверяем запрет операций с избранным без токена.
func (s *TestSuite) TestFavoritesWithoutToken() {
	t := s.T()
	token := s.loginData.Token

	created := helpers.CreateTestAdvertisement(t, token)

	t.Cleanup(func() {
		t.Logf("Удаляем объявление %s", created.ID)
		advertisements.DeleteAdvertisement(t, token, created.ID)
	})

	s.Run("Объявление нельзя добавить без токена", func() {
		t := s.T()
		t.Log("Отправляем запрос с ожидаемой ошибкой")
		resp := favManager.PostFavoriteExpectError(t, "", created.ID, http.StatusUnauthorized)

		assertions.APIErrorEquals(t, resp, helpers.UnauthorizedError, helpers.AuthorizationRequiredMessage)
	})

	s.Run("Список избранного нельзя получить без токена", func() {
		t := s.T()
		t.Log("Отправляем запрос с ожидаемой ошибкой")
		resp := favManager.GetFavoritesListExpectError(t, "", favoritesPage, favoritesPageSize, http.StatusUnauthorized)

		assertions.APIErrorEquals(t, resp, helpers.UnauthorizedError, helpers.AuthorizationRequiredMessage)
	})

	s.Run("Объявление нельзя удалить из избранного без токена", func() {
		t := s.T()
		t.Log("Отправляем запрос с ожидаемой ошибкой")
		resp := favManager.DeleteFavoriteExpectError(t, "", created.ID, http.StatusUnauthorized)

		assertions.APIErrorEquals(t, resp, helpers.UnauthorizedError, helpers.AuthorizationRequiredMessage)
	})
}

// favoriteAdvertisementIDs возвращает ID всех объявлений в избранном при лимите API в 50 записей.
func favoriteAdvertisementIDs(t *testing.T, token string) []string {
	t.Helper()
	t.Log("Получаем избранное")
	list := favManager.GetFavoritesList(t, token, favoritesPage, favoritesPageSize)
	ids := make([]string, 0, len(list.Items))

	for _, item := range list.Items {
		ids = append(ids, item.Advertisement.ID)
	}

	return ids
}
