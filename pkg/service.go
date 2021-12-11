package pkg

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go

//Service interface
type Service interface {
	SaveURL(URL string) (string, error)
	GetURL(link string) (string, error)
}
