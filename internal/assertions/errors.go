package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/common/models"
)

// APIErrorEquals проверяет код и сообщение ошибки в JSON-ответе API.
func APIErrorEquals(t *testing.T, actual models.ErrorResponse, expectedError, expectedMessage string) {
	t.Helper()
	require.Equal(t, expectedError, actual.Error, "Неожиданный код ошибки API")
	require.Equal(t, expectedMessage, actual.Message, "Неожиданное сообщение ошибки API")
	t.Logf("Получена ожидаемая ошибка: %s, %s", actual.Error, actual.Message)
}
