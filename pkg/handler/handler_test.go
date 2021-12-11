package handler

import (
	mocks "LinkShortener/pkg/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testURL         = URL{URL: "https://github.com/test_URL"}
	testLink        = "1234567891"
	testInvalidLink = "4567891"
	testInvalidURL  = URL{URL: "github.com"}
)

func setUp(t *testing.T, URL, method string, testingURL URL) (echo.Context, Handler, *mocks.MockService) {
	e := echo.New()
	r := e.Router()
	r.Add(method, URL, func(ctx echo.Context) error { return nil })

	ctrl := gomock.NewController(t)
	service := mocks.NewMockService(ctrl)
	handler := Handler{
		services: service,
	}

	var req *http.Request
	switch method {
	case http.MethodPost:
		f, _ := json.Marshal(testingURL)
		req = httptest.NewRequest(http.MethodPost, URL, bytes.NewBuffer(f))
	case http.MethodGet:
		req = httptest.NewRequest(http.MethodGet, URL, nil)
	}

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(URL)

	return c, handler, service
}

func TestHandler_GetURL(t *testing.T) {
	c, h, service := setUp(t, "/:link", http.MethodGet, testURL)

	c.SetParamNames("link")
	c.SetParamValues(testLink)
	service.EXPECT().GetURL(testLink).Return(testURL.URL, nil)
	err := h.GetURL(c)

	assert.Nil(t, err)
}

func TestHandler_GetURLError(t *testing.T) {
	c, h, service := setUp(t, "/:link", http.MethodGet, testURL)
	c.SetParamNames("link")
	c.SetParamValues(testLink)

	service.EXPECT().GetURL(testLink).Return("", errors.New(""))

	err := h.GetURL(c)

	assert.NotNil(t, err)
}

func TestHandler_GetURLInvalidLink(t *testing.T) {
	c, h, _ := setUp(t, "/:link", http.MethodGet, testURL)

	err1 := isValidLink(testInvalidLink)
	err := h.GetURL(c)

	assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, err1.Error()), err)
}

func TestHandler_SaveURL(t *testing.T) {
	c, h, service := setUp(t, "/link", http.MethodPost, testURL)
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	service.EXPECT().SaveURL(testURL.URL).Return(testLink, nil)

	err := h.SaveURL(c)

	assert.Nil(t, err)
}

func TestHandler_SaveURLError(t *testing.T) {
	c, h, service := setUp(t, "/link", http.MethodPost, testURL)
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	service.EXPECT().SaveURL(testURL.URL).Return(testLink, errors.New(""))

	err := h.SaveURL(c)
	assert.NotNil(t, err)
}

func TestHandler_SaveURLErrorBind(t *testing.T) {
	c, h, _ := setUp(t, "/link", http.MethodPost, testURL)
	var URL URL
	err1 := c.Bind(&URL)

	err := h.SaveURL(c)

	assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, err1.Error()), err)
}

func TestHandler_GetURLInvalidURL(t *testing.T) {
	c, h, _ := setUp(t, "/:link", http.MethodGet, testURL)

	err1 := isValidLink(testInvalidURL.URL)
	err := h.GetURL(c)

	assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, err1.Error()), err)
}
