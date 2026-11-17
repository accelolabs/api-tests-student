package advertisementcreation

import (
	"api-tests-template/internal/managers/advertisements"

	"api-tests-template/tests/helpers"

	"github.com/stretchr/testify/require"
)

// Positive: Проверяем создание объявления с валидными полями
func (s *TestSuite) TestCreate() {
	request := helpers.NewAdvertisementRequest()

	t := s.T()
	token := s.loginData.Token

	t.Log("Создаём объявление")
	created := advertisements.PostAdvertisement(t, token, request)
	require.NotEmpty(t, created.ID)
	t.Logf("Создано объявление %s", created.ID)

	t.Cleanup(func() {
		t.Logf("Удаляем объявление %s", created.ID)
		advertisements.DeleteAdvertisement(t, token, created.ID)
	})

	s.Run("Объявление доступно по ID и содержит переданные данные", func() {
		t := s.T()

		t.Log("Получаем объявление по ID")
		actual := advertisements.GetAdvertisement(t, token, created.ID)

		require.Equal(t, created.ID, actual.ID)
		require.Equal(t, request.Title, actual.Title)
		require.Equal(t, request.Description, actual.Description)
		require.InDelta(t, request.Price, actual.Price, 0.001)
		require.Equal(t, request.Quantity, actual.Quantity)
		require.Equal(t, s.loginData.User.Id, actual.UserID)
		t.Log("Поля объявления совпадают")
	})

	s.Run("Объявление находится по полному названию", func() {
		t := s.T()

		t.Log("Ищем по полному названию")
		list := advertisements.GetAdvertisementsList(t, token, request.Title)

		found := false
		for _, item := range list.Items {
			if item.ID == created.ID {
				found = true
				break
			}
		}

		require.True(t, found, "Созданное объявление не найдено в поиске")
		t.Log("Объявление найдено")
	})

	s.Run("Фотографии доступны", func() {
		t := s.T()

		t.Log("Получаем фотографии")
		photos := advertisements.GetAdvertisementPhotos(t, token, created.ID)

		require.Len(t, photos, len(request.Photos))
		require.Len(t, created.Photos, len(request.Photos))

		for _, photo := range photos {
			require.NotEmpty(t, photo.ID)
			require.NotEmpty(t, photo.URL)
		}

		require.ElementsMatch(t, created.Photos, photos)
		t.Log("Фотографии совпадают")
	})
}
