package advertisementcreation

import (
	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"net/http"

	"api-tests-template/tests/helpers"
)

// Negative: Проверяем 400 ошибку при создании объявления с пустым названием
func (s *TestSuite) TestCreateTitleEmpty() {
	request := helpers.NewAdvertisementRequest()
	request.Title = ""

	t := s.T()
	token := s.loginData.Token

	t.Log("Отправляем запрос с ожидаемой ошибкой")
	resp := advertisements.PostAdvertisementExpectError(t, token, request, http.StatusBadRequest)

	s.Run("Объявление без названия не должно быть создано", func() {
		t := s.T()

		assertions.APIErrorEquals(t, resp, helpers.InvalidRequestError, "title is required")
	})
}
