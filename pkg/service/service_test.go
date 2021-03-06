package service

import (
	"LinkShortener/pkg"
	"LinkShortener/pkg/handler"
	mocks "LinkShortener/pkg/mocks"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	testURL  = handler.URL{URL: "https://github.com/test_URL/"}
	testLink = "1234567891"
)

func setUp(t *testing.T) (*mocks.MockRepository, *mocks.MockDBRepository, pkg.Service) {
	ctrl := gomock.NewController(t)

	rep := mocks.NewMockRepository(ctrl)
	repDb := mocks.NewMockDBRepository(ctrl)

	service := NewService(rep, repDb)
	return rep, repDb, service
}

func TestService_SaveURL(t *testing.T) {
	rep, repDb, service := setUp(t)
	if os.Getenv("DB") == "true" {
		repDb.EXPECT().DBCheckURL(testURL.URL).Return("", errors.New(""))
		repDb.EXPECT().DBSaveURL(testURL.URL, gomock.Any()).AnyTimes().Return(nil)

		_, err := service.SaveURL(testURL.URL)
		assert.Nil(t, err)

		return
	}

	rep.EXPECT().SaveURL(testURL.URL, gomock.Any()).Return(testLink, nil)
	_, err := service.SaveURL(testURL.URL)

	assert.Nil(t, err)
}

func TestService_GetURL(t *testing.T) {
	rep, repDb, service := setUp(t)
	if os.Getenv("DB") == "true" {
		repDb.EXPECT().DBGetURL(testLink).Return(testURL.URL, nil)
		_, err := service.GetURL(testLink)
		assert.Nil(t, err)
		return

	}

	rep.EXPECT().GetURL(testLink).Return(testURL.URL, nil)

	_, err := service.GetURL(testLink)
	assert.Nil(t, err)
}

func TestService_GetURLError(t *testing.T) {
	rep, repDb, service := setUp(t)
	if os.Getenv("DB") == "true" {
		repDb.EXPECT().DBGetURL(testLink).Return("", errors.New(""))
		_, err := service.GetURL(testLink)
		assert.NotNil(t, err)

		return
	}

	rep.EXPECT().GetURL(testLink).Return("", errors.New(""))

	_, err := service.GetURL(testLink)
	assert.NotNil(t, err)
}
