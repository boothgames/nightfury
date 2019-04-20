package nightfury_test

import (
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"gitlab.com/jskswamy/nightfury/pkg/db"
	mocks "gitlab.com/jskswamy/nightfury/pkg/internal/mocks/db"
	"gitlab.com/jskswamy/nightfury/pkg/nightfury"
	"testing"
)

func TestSecurityIncidentSave(t *testing.T) {
	t.Run("should be able to save security Incident", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		securityIncident := nightfury.SecurityIncident{
			Title:    "title",
			Content:  "content",
			Tag:      "web",
			Takeaway: "dont do this",
		}
		repository.EXPECT().Save("securityIncidents", securityIncident)

		err := securityIncident.Save(repository)

		assert.NoError(t, err)
	})

	t.Run("should return error returned by repository save", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		securityIncident := nightfury.SecurityIncident{Title: "title"}
		repository.EXPECT().Save("securityIncidents", securityIncident).Return(fmt.Errorf("unable to save"))

		err := securityIncident.Save(repository)

		if assert.Error(t, err) {
			assert.Equal(t, "unable to save", err.Error())
		}
	})
}

func TestNewSecurityIncidentFromRepoWithName(t *testing.T) {
	t.Run("should fetch the security incident from db", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		repository.EXPECT().Fetch("securityIncidents", "one", gomock.Any()).Return(true, nil)

		actual, err := nightfury.NewSecurityIncidentFromRepoWithName(repository, "one")

		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})

	t.Run("should fail to fetch the client from db", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		repository.EXPECT().Fetch("securityIncidents", "one", gomock.Any()).Return(false, nil)

		actual, err := nightfury.NewSecurityIncidentFromRepoWithName(repository, "one")

		if assert.Error(t, err) {
			assert.Equal(t, "securityIncident with name one doesn't exists", err.Error())
		}
		assert.Equal(t, nightfury.SecurityIncident{}, actual)
	})

	t.Run("should return error returned while fetching data from repo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		repository.EXPECT().Fetch("securityIncidents", "one", gomock.Any()).Return(false, fmt.Errorf("unable to fetch"))

		actual, err := nightfury.NewSecurityIncidentFromRepoWithName(repository, "one")

		if assert.Error(t, err) {
			assert.Equal(t, "unable to fetch", err.Error())
		}
		assert.Equal(t, nightfury.SecurityIncident{}, actual)
	})
}

func TestNewSecurityIncidentsFromRepo(t *testing.T) {
	t.Run("should be able to get all the security incidents", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expected := nightfury.SecurityIncidents{
			"title": nightfury.SecurityIncident{
				Title:    "title",
				Content:  "content",
				Tag:      "web",
				Takeaway: "dont do this",
			},
		}
		repository := mocks.NewMockRepository(ctrl)
		repository.EXPECT().FetchAll("securityIncidents", gomock.Any()).DoAndReturn(
			func(bucketName string, modelFn func(data []byte) (db.Model, error)) (interface{}, error) {
				data, _ := json.Marshal(nightfury.SecurityIncident{
					Title:    "title",
					Content:  "content",
					Tag:      "web",
					Takeaway: "dont do this",
				})
				model, err := modelFn(data)
				if err != nil {
					return nil, err
				}
				return nightfury.SecurityIncidents{model.ID(): model.(nightfury.SecurityIncident)}, nil
			})

		securityIncidents, err := nightfury.NewSecurityIncidentsFromRepo(repository)

		assert.NoError(t, err)
		if !cmp.Equal(expected, securityIncidents) {
			assert.Fail(t, cmp.Diff(nightfury.SecurityIncident{}, securityIncidents))
		}
	})
}

func TestSecurityIncidentDelete(t *testing.T) {
	t.Run("should be able to delete security incident", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		securityIncident := nightfury.SecurityIncident{Title: "securityIncident"}
		repository.EXPECT().Delete("securityIncidents", securityIncident)

		err := securityIncident.Delete(repository)

		assert.NoError(t, err)
	})

	t.Run("should return error returned by repository delete", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		securityIncident := nightfury.SecurityIncident{Title: "securityIncident"}
		repository.EXPECT().Delete("securityIncidents", securityIncident).Return(fmt.Errorf("unable to delete"))

		err := securityIncident.Delete(repository)

		if assert.Error(t, err) {
			assert.Equal(t, "unable to delete", err.Error())
		}
	})
}
