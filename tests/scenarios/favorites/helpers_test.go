package favorites

import (
	"testing"

	favManager "api-tests-template/internal/managers/favorites"
)

// favoriteAdvertisementIDs возвращает ID всех объявлений в избранном при лимите API в 50 записей.
func favoriteAdvertisementIDs(t *testing.T, token string) []string {
	t.Helper()
	t.Log("Получаем избранное")
	list := favManager.GetFavoritesList(t, token, favoritesPage, favoritesPageSize)
	ids := make([]string, 0, len(list.Items))

	for _, item := range list.Items {
		ids = append(ids, item.Advertisement.ID)
	}

	return ids
}
