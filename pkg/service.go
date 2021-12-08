package pkg

type Service interface {
	SaveURL(URL string) (string, error)
	GetURL(link string) (string, error)
}
