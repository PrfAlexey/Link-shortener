package service

import (
	"LinkShortener/pkg"
	"math/rand"
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
	if DataBase {
		var link string
		for {
			link = GenerateLink(URL)
			if err := s.dbrepo.DBSaveURL(URL, link); err == nil {
				return link, nil
			}
		}
	}
	var link string
	for {
		link = GenerateLink(URL)
		if err := s.repo.SaveURL(URL, link); err == nil {
			return link, nil
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
	for i := range link {
		link[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(link)
}
