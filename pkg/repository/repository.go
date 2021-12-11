package repository

import (
	"LinkShortener/pkg"
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

//Repository in memory
type Repository struct {
	storage map[string]string
}

//DBRepository ...
type DBRepository struct {
	pool *pgxpool.Pool
}

//NewRepository ...
func NewRepository(storage map[string]string) pkg.Repository {
	return &Repository{
		storage: storage,
	}
}

//NewDBRepository ...
func NewDBRepository(pool *pgxpool.Pool) pkg.DBRepository {
	return &DBRepository{
		pool: pool,
	}
}

//SaveURL ...
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

//GetURL ...
func (r *Repository) GetURL(link string) (string, error) {
	if _, inMap := r.storage[link]; inMap {
		return r.storage[link], nil
	}
	return "", errors.New("there is no URL for this link")
}

//DBSaveURL ...
func (r *DBRepository) DBSaveURL(URL, link string) error {
	_, err := r.pool.Exec(context.Background(), `INSERT INTO urlandlinks (url, link) values ($1, $2)`, URL, link)
	if err != nil {
		fmt.Println(err.Error())
		return errors.New("duplicate link")
	}
	return nil
}

//DBGetURL ...
func (r *DBRepository) DBGetURL(link string) (string, error) {
	var URL []string
	err := pgxscan.Select(context.Background(), r.pool, &URL,
		`SELECT url FROM urlandlinks WHERE link = $1`, link)
	if errors.As(err, &pgx.ErrNoRows) {
		return "", errors.New("Invalid link")
	}
	return URL[0], nil
}

//DBCheckURL ...
func (r *DBRepository) DBCheckURL(URL string) (string, error) {
	var link []string

	err := pgxscan.Select(context.Background(), r.pool, &link,
		`SELECT link FROM urlandlinks WHERE url = $1`, URL)

	if errors.As(err, &pgx.ErrNoRows) || len(link) == 0 {
		return "", errors.New("there is no such URL")
	}
	return link[0], nil
}
