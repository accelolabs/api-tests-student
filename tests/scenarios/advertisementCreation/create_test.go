package advertisementcreation

import (
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"os"
	"testing"
	"time"

	advModels "api-tests-template/internal/managers/advertisements/models"
	base "api-tests-template/tests"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя с кредами из переменных окружения и получение его параметров")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TestCreate() {
	title := "Лемуры " + time.Now().Format("2006-01-02 15:04:05.000000000")
	description := "Сидят на ветке и что-то замышляют на будущее."
	price := 1337.67
	quantity := 2

	request := advModels.NewAdvertisement{
		Title:       &title,
		Description: &description,
		Price:       &price,
		Quantity:    &quantity,
		Photos:      []string{"../../testdata/photo.jpg"},
	}

	t := s.T()
	token := s.loginData.Token

	created := advertisements.PostAdvertisement(t, token, request)
	require.NotEmpty(t, created.ID)

	t.Cleanup(func() {
		advertisements.DeleteAdvertisement(t, token, created.ID)
	})

	s.Run("Объявление доступно по ID и содержит переданные данные", func() {
		t := s.T()

		actual := advertisements.GetAdvertisement(t, token, created.ID)

		require.Equal(t, created.ID, actual.ID)
		require.Equal(t, title, actual.Title)
		require.Equal(t, description, actual.Description)
		require.InDelta(t, price, actual.Price, 0.001)
		require.Equal(t, quantity, actual.Quantity)
		require.Equal(t, s.loginData.User.Id, actual.UserID)
	})

	s.Run("Объявление находится по полному названию", func() {
		t := s.T()

		list := advertisements.GetAdvertisementsList(t, token, title)

		found := false
		for _, item := range list.Items {
			if item.ID == created.ID {
				found = true
				break
			}
		}

		require.True(t, found, "Созданное объявление не найдено в поиске")
	})

	s.Run("Фотографии доступны", func() {
		t := s.T()

		photos := advertisements.GetAdvertisementPhotos(t, token, created.ID)

		require.Len(t, photos, len(request.Photos))
		require.Len(t, created.Photos, len(request.Photos))

		for _, photo := range photos {
			require.NotEmpty(t, photo.ID)
			require.NotEmpty(t, photo.URL)
		}

		require.ElementsMatch(t, created.Photos, photos)
	})
}
