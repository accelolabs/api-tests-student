package favorites

import (
	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"net/http"

	"api-tests-template/tests/helpers"

	favManager "api-tests-template/internal/managers/favorites"
)

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
