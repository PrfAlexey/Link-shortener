package service

import (
	"LinkShortener/pkg"
	"math/rand"
	"os"
	"time"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

//Service struct
type Service struct {
	repo   pkg.Repository
	dbrepo pkg.DBRepository
}

//NewService initializes new Service
func NewService(repo pkg.Repository, dbrepo pkg.DBRepository) pkg.Service {
	return &Service{
		repo:   repo,
		dbrepo: dbrepo,
	}
}

//SaveURL transfers URL to the repository layer depending on the choice of repository
func (s *Service) SaveURL(URL string) (string, error) {

	var link string
	if os.Getenv("DB") == "true" {
		if link, err := s.dbrepo.DBCheckURL(URL); err == nil {
			return link, err
		}

		for {
			link = generateLink()
			if err := s.dbrepo.DBSaveURL(URL, link); err == nil {
				return link, nil
			}
		}
	}

	for {
		link = generateLink()
		if shortURL, err := s.repo.SaveURL(URL, link); err == nil {
			return shortURL, nil
		}
	}

}

//GetURL gets URL by link from the repository layer depending on the choice of repository
func (s *Service) GetURL(link string) (string, error) {
	if os.Getenv("DB") == "true" {
		URL, err := s.dbrepo.DBGetURL(link)
		return URL, err
	}
	URL, err := s.repo.GetURL(link)
	return URL, err
}

func generateLink() string {
	link := make([]byte, 10)
	rand.Seed(time.Now().UnixNano())
	for i := range link {
		link[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(link)
}
