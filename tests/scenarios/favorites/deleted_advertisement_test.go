package favorites

import (
	"api-tests-template/internal/managers/advertisements"

	"api-tests-template/tests/helpers"

	favManager "api-tests-template/internal/managers/favorites"

	"github.com/stretchr/testify/require"
)

// Positive: Проверяем исчезновение удалённого объявления из избранного.
func (s *TestSuite) TestDeletedAdvertisement() {
	t := s.T()
	token := s.loginData.Token

	created := helpers.CreateTestAdvertisement(t, token)

	token2 := s.loginData2.Token
	t.Logf("Добавляем %s в избранное", created.ID)
	favManager.PostFavorite(t, token2, created.ID)

	s.Run("Добавленное объявление есть в избранном", func() {
		t := s.T()

		ids := favoriteAdvertisementIDs(t, s.loginData2.Token)

		require.Contains(t, ids, created.ID)
		t.Log("Объявление есть в избранном")
	})

	t.Logf("Удаляем объявление %s", created.ID)
	advertisements.DeleteAdvertisement(t, token, created.ID)

	s.Run("Удаленное автором объявление отсутствует в избранном", func() {
		t := s.T()

		ids := favoriteAdvertisementIDs(t, s.loginData2.Token)

		require.NotContains(t, ids, created.ID)
		t.Log("Объявления нет в избранном")
	})
}
