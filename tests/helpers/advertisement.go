package helpers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
)

const (
	testTitlePrefix = "Лемуры "
	testDescription = "Сидят на ветке и что-то замышляют на будущее."
	testPrice       = 1337.67
	testQuantity    = 2
	testPhotoPath   = "../../testdata/photo.jpg"

	InvalidRequestError          = "invalid_request"
	UnauthorizedError            = "unauthorized"
	AuthorizationRequiredMessage = "Authorization header required"
)

// NewAdvertisementRequest готовит валидные данные с новым названием и отдельным списком фотографий.
func NewAdvertisementRequest() models.NewAdvertisement {
	return models.NewAdvertisement{
		Title:       testTitlePrefix + time.Now().Format("2006-01-02 15:04:05.000000000"),
		Description: testDescription,
		Price:       testPrice,
		Quantity:    testQuantity,
		Photos:      []string{testPhotoPath},
	}
}

// CreateTestAdvertisement создаёт объявление для сценария; удаление выполняет вызывающий тест.
func CreateTestAdvertisement(t *testing.T, token string) models.Advertisement {
	t.Helper()
	t.Log("Создаём тестовое объявление")
	created := advertisements.PostAdvertisement(t, token, NewAdvertisementRequest())
	require.NotEmpty(t, created.ID, "Сервер не вернул ID созданного объявления")
	t.Logf("Создано объявление %s", created.ID)
	return created
}
