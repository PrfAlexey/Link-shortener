package service

import (
	"LinkShortener/pkg"
	"math/rand"
	"time"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	DataBase    = true
)

type Service struct {
	repo   pkg.Repository
	dbrepo pkg.DBRepository
}

func NewService(repo pkg.Repository, dbrepo pkg.DBRepository) pkg.Service {
	return &Service{
		repo:   repo,
		dbrepo: dbrepo,
	}
}

func (s *Service) SaveURL(URL string) (string, error) {

	var link string
	if DataBase {
		if link, err := s.dbrepo.DBCheckURL(URL); err == nil {
			return link, err
		}

		for {
			link = GenerateLink(URL)
			if err := s.dbrepo.DBSaveURL(URL, link); err == nil {
				return link, nil
			}
		}
	}

	for {
		link = GenerateLink(URL)
		if shortURL, err := s.repo.SaveURL(URL, link); err == nil {
			return shortURL, nil
		}
	}

}

func (s *Service) GetURL(link string) (string, error) {
	if DataBase {
		URL, err := s.dbrepo.DBGetURL(link)
		return URL, err
	}
	URL, err := s.repo.GetURL(link)
	return URL, err
}

func GenerateLink(URL string) string {
	link := make([]byte, 10)
	rand.Seed(time.Now().UnixNano())
	for i := range link {
		link[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(link)
}
