package favorites

import (
	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"net/http"

	"api-tests-template/tests/helpers"

	favManager "api-tests-template/internal/managers/favorites"
)

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
