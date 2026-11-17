package favorites

import (
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"os"
	"testing"
	"time"

	base "api-tests-template/tests"

	advModels "api-tests-template/internal/managers/advertisements/models"
	favManager "api-tests-template/internal/managers/favorites"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	loginData  models.LoginOkResponse
	loginData2 models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователей с кредами из переменных окружения и получение их параметров")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
	s.loginData2 = auth.Login(s.T(), os.Getenv("TEST_LOGIN_2"), os.Getenv("TEST_PASSWORD_2"))
}

func (s *TestSuite) TestAddAndDelete() {
	t := s.T()
	token := s.loginData.Token

	requestCreateAdv := advModels.NewAdvertisement{
		Title:       "Лемуры " + time.Now().Format("2006-01-02 15:04:05.000000000"),
		Description: "Сидят на ветке и что-то замышляют на будущее.",
		Price:       1337.67,
		Quantity:    2,
		Photos:      []string{"../../testdata/photo.jpg"},
	}

	created := advertisements.PostAdvertisement(t, token, requestCreateAdv)
	require.NotEmpty(t, created.ID)

	t.Cleanup(func() {
		advertisements.DeleteAdvertisement(t, token, created.ID)
	})

	token2 := s.loginData2.Token
	favManager.PostFavorite(t, token2, created.ID)

	s.Run("Добавленное объявление есть в избранном", func() {
		t := s.T()

		list := favManager.GetFavoritesList(t, s.loginData2.Token, 1, 50)

		ids := make([]string, 0, len(list.Items))
		for _, item := range list.Items {
			ids = append(ids, item.Advertisement.ID)
		}

		require.Contains(t, ids, created.ID)
	})

	favManager.DeleteFavorite(t, token2, created.ID)

	s.Run("Удаленное объявление отсутствует в избранном", func() {
		t := s.T()

		list := favManager.GetFavoritesList(t, s.loginData2.Token, 1, 50)

		ids := make([]string, 0, len(list.Items))
		for _, item := range list.Items {
			ids = append(ids, item.Advertisement.ID)
		}

		require.NotContains(t, ids, created.ID)
	})
}

func (s *TestSuite) TestDeletedAdvertisement() {
	t := s.T()
	token := s.loginData.Token

	requestCreateAdv := advModels.NewAdvertisement{
		Title:       "Лемуры " + time.Now().Format("2006-01-02 15:04:05.000000000"),
		Description: "Сидят на ветке и что-то замышляют на будущее.",
		Price:       1337.67,
		Quantity:    2,
		Photos:      []string{"../../testdata/photo.jpg"},
	}

	created := advertisements.PostAdvertisement(t, token, requestCreateAdv)
	require.NotEmpty(t, created.ID)

	token2 := s.loginData2.Token
	favManager.PostFavorite(t, token2, created.ID)

	s.Run("Добавленное объявление есть в избранном", func() {
		t := s.T()

		list := favManager.GetFavoritesList(t, s.loginData2.Token, 1, 50)

		ids := make([]string, 0, len(list.Items))
		for _, item := range list.Items {
			ids = append(ids, item.Advertisement.ID)
		}

		require.Contains(t, ids, created.ID)
	})

	advertisements.DeleteAdvertisement(t, token, created.ID)

	s.Run("Удаленное автором объявление отсутствует в избранном", func() {
		t := s.T()

		list := favManager.GetFavoritesList(t, s.loginData2.Token, 1, 50)

		ids := make([]string, 0, len(list.Items))
		for _, item := range list.Items {
			ids = append(ids, item.Advertisement.ID)
		}

		require.NotContains(t, ids, created.ID)
	})
}
