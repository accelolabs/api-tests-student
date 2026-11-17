package advertisementcreation

import (
	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"net/http"

	"api-tests-template/tests/helpers"
)

// Negative: Проверяем 400 ошибку при создании объявления с отрицательной ценой
func (s *TestSuite) TestCreatePriceNegative() {
	request := helpers.NewAdvertisementRequest()
	request.Price = -request.Price

	t := s.T()
	token := s.loginData.Token

	t.Log("Отправляем запрос с ожидаемой ошибкой")
	resp := advertisements.PostAdvertisementExpectError(t, token, request, http.StatusBadRequest)

	s.Run("Объявление с отрицательной ценой не должно быть создано", func() {
		t := s.T()

		assertions.APIErrorEquals(t, resp, helpers.InvalidRequestError, "price must be between 0 and 1000000 inclusive")
	})
}
