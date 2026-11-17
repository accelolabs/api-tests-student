package advertisementcreation

import (
	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"net/http"

	"api-tests-template/tests/helpers"
)

// Negative: Проверяем 401 ошибку при попытке создать объявление без токена
func (s *TestSuite) TestCreateTokenEmpty() {
	request := helpers.NewAdvertisementRequest()

	t := s.T()
	t.Log("Отправляем запрос с ожидаемой ошибкой")
	resp := advertisements.PostAdvertisementExpectError(t, "", request, http.StatusUnauthorized)

	s.Run("Объявление нельзя создать без токена", func() {
		t := s.T()

		assertions.APIErrorEquals(t, resp, helpers.UnauthorizedError, helpers.AuthorizationRequiredMessage)
	})
}
