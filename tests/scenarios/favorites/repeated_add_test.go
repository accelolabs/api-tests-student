package favorites

import (
	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"net/http"

	"api-tests-template/tests/helpers"

	favManager "api-tests-template/internal/managers/favorites"
)

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
