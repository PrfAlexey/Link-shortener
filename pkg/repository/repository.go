package repository

import (
	"LinkShortener/pkg"
	"database/sql"
	"errors"
	"fmt"
)

type link struct {
	link string
}

type URL struct {
	URL string
}

//Repository in memory
type Repository struct {
	storage map[string]string
}

//DBRepository in Postgres DB
type DBRepository struct {
	db *sql.DB
}

//NewRepository initializes new in-memory repository
func NewRepository(storage map[string]string) pkg.Repository {
	return &Repository{
		storage: storage,
	}
}

//NewDBRepository initializes new Postgres repository
func NewDBRepository(db *sql.DB) pkg.DBRepository {
	return &DBRepository{
		db: db,
	}
}

//SaveURL saves URL and link in Map
func (r *Repository) SaveURL(URL, link string) (string, error) {
	if link, inMap := r.storage[URL]; inMap {
		return link, nil
	}
	if _, inMap := r.storage[link]; !inMap {
		r.storage[link] = URL
		r.storage[URL] = link
		return link, nil
	}
	return "", errors.New("Duplicate link")
}

//GetURL gets URL from Map
func (r *Repository) GetURL(link string) (string, error) {
	if _, inMap := r.storage[link]; inMap {
		return r.storage[link], nil
	}
	return "", errors.New("there is no URL for this link")
}

//DBSaveURL saves URL in Postgres DB
func (r *DBRepository) DBSaveURL(URL, link string) error {
	_, err := r.db.Exec(`INSERT INTO urlandlinks (url, link) values ($1, $2)`, URL, link)

	if err != nil {
		fmt.Println(err.Error())
		return errors.New("duplicate link")
	}
	return nil
}

//DBGetURL gets URL from Postgres DB
func (r *DBRepository) DBGetURL(link string) (string, error) {
	URL := &URL{}

	row := r.db.QueryRow(`SELECT url FROM urlandlinks WHERE link = $1`, link)
	err := row.Scan(&URL.URL)

	if err != nil {
		return "", errors.New("Invalid link")
	}

	return URL.URL, nil
}

//DBCheckURL checks URL in Postgres DB
func (r *DBRepository) DBCheckURL(url string) (string, error) {
	link := &link{}

	row := r.db.QueryRow(`SELECT link FROM urlandlinks WHERE url = $1`, url)

	err := row.Scan(&link.link)

	if err != nil {
		return "", errors.New("there is no such URL")
	}

	return link.link, nil
}
