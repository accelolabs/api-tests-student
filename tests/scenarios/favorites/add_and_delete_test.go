package favorites

import (
	"api-tests-template/internal/managers/advertisements"

	"api-tests-template/tests/helpers"

	favManager "api-tests-template/internal/managers/favorites"

	"github.com/stretchr/testify/require"
)

// Positive: Проверяем добавление в избранное и удаление из него.
func (s *TestSuite) TestAddAndDelete() {
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

	s.Run("Добавленное объявление есть в избранном", func() {
		t := s.T()

		ids := favoriteAdvertisementIDs(t, s.loginData2.Token)

		require.Contains(t, ids, created.ID)
		t.Log("Объявление есть в избранном")
	})

	t.Logf("Удаляем %s из избранного", created.ID)
	favManager.DeleteFavorite(t, token2, created.ID)

	s.Run("Удаленное объявление отсутствует в избранном", func() {
		t := s.T()

		ids := favoriteAdvertisementIDs(t, s.loginData2.Token)

		require.NotContains(t, ids, created.ID)
		t.Log("Объявления нет в избранном")
	})
}
