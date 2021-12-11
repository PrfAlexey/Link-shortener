package pkg

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go

//Repository interface in memory
type Repository interface {
	SaveURL(URL, link string) (string, error)
	GetURL(link string) (string, error)
}

//DBRepository interface for Postgres
type DBRepository interface {
	DBSaveURL(URL, link string) error
	DBGetURL(link string) (string, error)
	DBCheckURL(URL string) (string, error)
}
