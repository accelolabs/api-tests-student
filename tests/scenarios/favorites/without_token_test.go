package favorites

import (
	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"net/http"

	"api-tests-template/tests/helpers"

	favManager "api-tests-template/internal/managers/favorites"
)

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
