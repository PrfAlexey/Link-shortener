package pkg

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go

type Repository interface {
	SaveURL(URL, link string) (string, error)
	GetURL(link string) (string, error)
}

type DBRepository interface {
	DBSaveURL(URL, link string) error
	DBGetURL(link string) (string, error)
	DBCheckURL(URL string) (string, error)
}
