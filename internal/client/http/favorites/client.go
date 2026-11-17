package favorites

import (
	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/favorites/models"
	"fmt"
	"net/http"
	"testing"
)

// HttpGetFavorites получает страницу избранного и проверяет ожидаемый статус.
func HttpGetFavorites(
	t *testing.T,
	token string,
	page int,
	pageSize int,
	expectedStatus int,
) *http.Response {
	t.Helper()

	return apiRunner.
		GetRunner().
		Auth(token).
		Create().
		Get(path.FavoritesPath).
		Query("page", fmt.Sprint(page)).
		Query("page_size", fmt.Sprint(pageSize)).
		Expect(t).
		Status(expectedStatus).
		End().Response
}

// HttpPostFavorites добавляет объявление в избранное и проверяет ожидаемый статус.
func HttpPostFavorites(t *testing.T, token string, id string, expectedStatus int) *http.Response {
	t.Helper()

	return apiRunner.
		GetRunner().
		Auth(token).
		Create().
		Post(path.FavoritesPath).
		ContentType("application/json").
		JSON(models.AddFavoriteRequest{AdvertisementID: id}).
		Expect(t).
		Status(expectedStatus).
		End().Response
}

// HttpDeleteFavorites удаляет объявление из избранного и проверяет ожидаемый статус.
func HttpDeleteFavorites(t *testing.T, token string, id string, expectedStatus int) *http.Response {
	t.Helper()

	return apiRunner.
		GetRunner().
		Auth(token).
		Create().
		Delete(fmt.Sprintf(path.FavoritesDeletePath, id)).
		Expect(t).
		Status(expectedStatus).
		End().Response
}
