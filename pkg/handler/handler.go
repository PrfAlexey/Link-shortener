package handler

import (
	"LinkShortener/pkg"
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"net/url"
)

type URL struct {
	URL string `json:"URL"`
}

type Handler struct {
	services pkg.Service
}

func NewHandler(services pkg.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitHandler(e *echo.Echo) {

	e.POST("/link", h.SaveURL)
	e.GET("/:link", h.GetURL)
}

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

func (h *Handler) GetURL(c echo.Context) error {
	fmt.Println(c.Path())
	link := c.Get("link").(string)

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
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("This is not a valid URL."))
	}
	u, err := url.Parse(URL)
	if err != nil || u.Host == "" {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("This is not a valid URL."))
	}
	return nil
}
