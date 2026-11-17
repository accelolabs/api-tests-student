package favorites

import (
	"api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	commonModels "api-tests-template/internal/managers/common/models"

	"github.com/stretchr/testify/require"
)

// GetFavoritesList возвращает страницу избранного пользователя.
func GetFavoritesList(t *testing.T, token string, page int, pageSize int) models.FavoritesList {
	t.Helper()

	raw := favorites.HttpGetFavorites(t, token, page, pageSize, http.StatusOK)

	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := models.FavoritesList{}
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}

// PostFavorite добавляет объявление в избранное и возвращает созданную запись.
func PostFavorite(t *testing.T, token string, advertisementID string) models.FavoriteItem {
	t.Helper()

	raw := favorites.HttpPostFavorites(t, token, advertisementID, http.StatusCreated)

	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := models.FavoriteItem{}
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}

// DeleteFavorite удаляет объявление из избранного пользователя.
func DeleteFavorite(t *testing.T, token string, advertisementID string) {
	t.Helper()

	raw := favorites.HttpDeleteFavorites(t, token, advertisementID, http.StatusNoContent)
	defer raw.Body.Close()
}

// GetFavoritesListExpectError проверяет статус отказа получения избранного и возвращает ошибку API.
func GetFavoritesListExpectError(
	t *testing.T,
	token string,
	page int,
	pageSize int,
	expectedStatus int,
) commonModels.ErrorResponse {
	t.Helper()

	raw := favorites.HttpGetFavorites(t, token, page, pageSize, expectedStatus)

	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := commonModels.ErrorResponse{}
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}

// PostFavoriteExpectError проверяет статус отказа добавления в избранное и возвращает ошибку API.
func PostFavoriteExpectError(
	t *testing.T,
	token string,
	advertisementID string,
	expectedStatus int,
) commonModels.ErrorResponse {
	t.Helper()

	raw := favorites.HttpPostFavorites(t, token, advertisementID, expectedStatus)

	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := commonModels.ErrorResponse{}
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}

// DeleteFavoriteExpectError проверяет статус отказа удаления из избранного и возвращает ошибку API.
func DeleteFavoriteExpectError(
	t *testing.T,
	token string,
	advertisementID string,
	expectedStatus int,
) commonModels.ErrorResponse {
	t.Helper()

	raw := favorites.HttpDeleteFavorites(t, token, advertisementID, expectedStatus)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := commonModels.ErrorResponse{}
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}
