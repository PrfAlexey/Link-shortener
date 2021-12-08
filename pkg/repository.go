package pkg

type Repository interface {
	SaveURL(URL, link string) error
	GetURL(link string) (string, error)
}

type DBRepository interface {
	DBSaveURL(URL, link string) error
	DBGetURL(link string) (string, error)
}
