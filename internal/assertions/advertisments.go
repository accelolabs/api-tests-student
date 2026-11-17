package assertions

import (
	"api-tests-template/internal/managers/advertisements/models"
	"testing"
)

// ContainsAdvertisementByTitle проверяет наличие объявления с точным названием.
func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	t.Helper()

	for _, item := range items {
		if item.Title == title {
			return
		}
	}

	t.Fatalf("Объявление с названием %q не найдено", title)
}
