package advertisements

import (
	"fmt"
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisements/models"
)

func HttpGetAdvertisements(t *testing.T, token string, textForSearch string) *http.Response {
	t.Helper()

	return apiRunner.
		GetRunner().
		Auth(token).
		Create().
		Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Expect(t).
		Status(http.StatusOK).
		End().Response
}

func HttpPostAdvertisement(
	t *testing.T,
	token string,
	advertisement models.NewAdvertisement,
	expectedStatus int,
) *http.Response {
	t.Helper()

	request := apiRunner.GetRunner().Auth(token).Create().Post(path.AdvertisementPath)

	if advertisement.Title != nil {
		request.MultipartFormData("title", *advertisement.Title)
	}

	if advertisement.Description != nil {
		request.MultipartFormData("description", *advertisement.Description)
	}

	if advertisement.Price != nil {
		request.MultipartFormData("price", fmt.Sprint(*advertisement.Price))
	}

	if advertisement.Quantity != nil {
		request.MultipartFormData("quantity", fmt.Sprint(*advertisement.Quantity))
	}

	if len(advertisement.Photos) > 0 {
		request.MultipartFile("photos", advertisement.Photos...)
	}

	return request.Expect(t).Status(expectedStatus).End().Response
}

func HttpGetAdvertisement(t *testing.T, token string, id string) *http.Response {
	t.Helper()

	return apiRunner.
		GetRunner().
		Auth(token).
		Create().
		Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		Status(http.StatusOK).
		End().Response
}

func HttpDeleteAdvertisement(t *testing.T, token string, id string) *http.Response {
	t.Helper()

	return apiRunner.
		GetRunner().
		Auth(token).
		Create().
		Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		Status(http.StatusNoContent).
		End().Response
}

func HttpGetAdvertisementPhotos(t *testing.T, token string, id string) *http.Response {
	t.Helper()

	return apiRunner.
		GetRunner().
		Auth(token).
		Create().
		Get(fmt.Sprintf(path.AdvertisementPhotosPath, id)).
		Expect(t).
		Status(http.StatusOK).
		End().Response
}
