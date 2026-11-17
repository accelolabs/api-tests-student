package advertisementcreation

import (
	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"net/http"
	"strings"

	"api-tests-template/tests/helpers"
)

// Negative: Проверяем 400 ошибку при создании объявления со слишком длинным названием
func (s *TestSuite) TestCreateTitleTooLong() {
	request := helpers.NewAdvertisementRequest()
	request.Title = strings.Repeat("Лемур", 11)

	t := s.T()
	token := s.loginData.Token

	t.Log("Отправляем запрос с ожидаемой ошибкой")
	resp := advertisements.PostAdvertisementExpectError(t, token, request, http.StatusBadRequest)

	s.Run("Объявление со слишком длинным названием не должно быть создано", func() {
		t := s.T()

		assertions.APIErrorEquals(t, resp, helpers.InvalidRequestError, "title must not exceed 50 characters provided 55")
	})
}
