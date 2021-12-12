package handler

import (
	"LinkShortener/pkg"
	"errors"
	"github.com/labstack/echo"
	"net/http"
	"net/url"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

//URL struct for parsing JSON
type URL struct {
	URL string `json:"URL"`
}

//Handler struct
type Handler struct {
	services pkg.Service
}

//NewHandler create a new Handler
func NewHandler(services pkg.Service) *Handler {
	return &Handler{
		services: services,
	}
}

//InitHandler initializes handlers
func (h *Handler) InitHandler(e *echo.Echo) {

	e.POST("/link", h.SaveURL)
	e.GET("/:link", h.GetURL)
}

//SaveURL parses body, checks it and transfers to the service layer
func (h *Handler) SaveURL(c echo.Context) error {
	var URL URL

	if err := c.Bind(&URL); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err := isValidURL(URL.URL); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	link, err := h.services.SaveURL(URL.URL)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"link": link,
	})
	return nil

}

//GetURL gets the "link" parameter, checks it and transfers to the service layer
func (h *Handler) GetURL(c echo.Context) error {

	link := c.Param("link")

	if err := isValidLink(link); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	URL, err := h.services.GetURL(link)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = c.JSON(http.StatusOK, map[string]interface{}{
		"URL": URL,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func isValidURL(URL string) error {
	if _, err := url.ParseRequestURI(URL); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("this is not a valid URL"))
	}
	u, err := url.Parse(URL)
	if err != nil || u.Host == "" {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("this is not a valid URL"))
	}
	return nil
}

func isValidLink(link string) error {
	rs := []rune(link)

	if len(rs) != 10 {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("this is not a valid link"))
	}
	m := make(map[rune]struct{})
	for _, ch := range []rune(letterBytes) {
		m[ch] = struct{}{}
	}

	for _, ch := range rs {
		if _, inMap := m[ch]; !inMap {
			return echo.NewHTTPError(http.StatusBadRequest, errors.New("this is not a valid link"))
		}
	}
	return nil
}
