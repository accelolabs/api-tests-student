package advertisements

import (
	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	commonModels "api-tests-template/internal/managers/common/models"

	"github.com/stretchr/testify/require"
)

func GetAdvertisementsList(t *testing.T, token string, text string) models.AdvertisementsList {
	raw := advertisements.HttpGetAdvertisements(t, token, text)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := models.AdvertisementsList{}
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}

func PostAdvertisement(t *testing.T, token string, request models.NewAdvertisement) models.Advertisement {
	t.Helper()

	raw := advertisements.HttpPostAdvertisement(t, token, request, http.StatusCreated)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := models.Advertisement{}
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}

func PostAdvertisementExpectError(
	t *testing.T,
	token string,
	request models.NewAdvertisement,
	expectedStatus int,
) commonModels.ErrorResponse {
	t.Helper()

	raw := advertisements.HttpPostAdvertisement(t, token, request, expectedStatus)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := commonModels.ErrorResponse{}
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}

func GetAdvertisement(t *testing.T, token string, id string) models.Advertisement {
	t.Helper()

	raw := advertisements.HttpGetAdvertisement(t, token, id)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := models.Advertisement{}
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}

func DeleteAdvertisement(t *testing.T, token string, id string) {
	t.Helper()

	raw := advertisements.HttpDeleteAdvertisement(t, token, id)
	defer raw.Body.Close()
}

func GetAdvertisementPhotos(t *testing.T, token string, id string) []models.Photo {
	t.Helper()

	raw := advertisements.HttpGetAdvertisementPhotos(t, token, id)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := []models.Photo{}
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}
