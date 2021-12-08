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

type Repository struct {
	storage map[string]string
}

type DBRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(storage map[string]string) pkg.Repository {
	return &Repository{
		storage: storage,
	}
}

func NewDBRepository(pool *pgxpool.Pool) pkg.DBRepository {
	return &DBRepository{
		pool: pool,
	}
}

func (r *Repository) SaveURL(URL, link string) error {
	if _, inMap := r.storage[link]; !inMap {
		r.storage[link] = URL
		return nil
	}
	return errors.New("Duplicate link")
}

func (r *Repository) GetURL(link string) (string, error) {
	if _, inMap := r.storage[link]; inMap {
		return r.storage[link], nil
	}
	return "", errors.New("There is no URL for this link")
}

func (r *DBRepository) DBSaveURL(URL, link string) error {
	_, err := r.pool.Exec(context.Background(), `INSERT INTO urlandlinks (url, link) values ($1, $2)`, URL, link)
	if err != nil {
		fmt.Println(err.Error())
		return errors.New("Duplicate link")
	}
	return nil
}

func (r *DBRepository) DBGetURL(link string) (string, error) {
	var URL []string
	println(link)
	err := pgxscan.Select(context.Background(), r.pool, &URL,
		`SELECT url FROM urlandlinks WHERE link = $1`, link)
	if errors.As(err, &pgx.ErrNoRows) {
		return "", errors.New("Invalid link111")
	}
	return URL[0], nil
}
