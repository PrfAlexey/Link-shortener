package repository

import (
	"LinkShortener/pkg"
	mocks "LinkShortener/pkg/mocks"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

const dbConnect = "user=postgres dbname=postgres1 password=4444 host=localhost port=5432 sslmode=disable pool_max_conns=50"

var (
	testLink        = "Q2zmAaE9Cy"
	testURL         = "https://ru.stackoverflow.com/"
	testInvalidLink = "1224567890"
	testInvalidURL  = "github"
)

func setUp(t *testing.T) (*mocks.MockRepository, pkg.Repository) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := map[string]string{
		testLink: testURL,
		testURL:  testLink,
	}
	store := mocks.NewMockRepository(ctrl)
	repo := NewRepository(m)
	return store, repo
}

func TestRepository_GetURL(t *testing.T) {
	store, repo := setUp(t)
	store.EXPECT().GetURL(testLink).Return(testURL, nil)

	_, err := repo.GetURL(testLink)

	assert.Nil(t, err)
}

func TestRepository_GetURLError(t *testing.T) {
	store, repo := setUp(t)
	store.EXPECT().GetURL(testInvalidLink).Return("", errors.New(""))

	_, err := repo.GetURL(testInvalidLink)

	assert.NotNil(t, err)
}

func TestRepository_SaveURL(t *testing.T) {
	store, repo := setUp(t)
	store.EXPECT().SaveURL(testURL, testLink).Return(testLink, nil)

	_, err := repo.SaveURL(testURL, testLink)

	assert.Nil(t, err)
}

func TestRepository_SaveURLError(t *testing.T) {
	store, repo := setUp(t)
	store.EXPECT().SaveURL(testInvalidURL, testLink).Return("", errors.New(""))

	_, err := repo.SaveURL(testInvalidURL, testLink)

	assert.NotNil(t, err)
}

func TestDBRepository_DBCheckURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewDBRepository(db)

	rows := sqlmock.NewRows([]string{"link"}).AddRow(testLink)

	mock.ExpectQuery("SELECT link").WithArgs(testURL).WillReturnRows(rows)

	_, err1 := repo.DBCheckURL(testURL)
	assert.Nil(t, err1)
}

func TestDBRepository_DBCheckURLError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewDBRepository(db)

	mock.ExpectQuery("SELECT link").WithArgs(testURL).WillReturnError(errors.New(""))

	_, err1 := repo.DBCheckURL(testURL)
	assert.NotNil(t, err1)
}

func TestDBRepository_DBSaveURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewDBRepository(db)

	mock.ExpectExec(`INSERT INTO urlandlinks`).WithArgs(testURL, testLink).WillReturnResult(sqlmock.NewResult(1, 1))

	err1 := repo.DBSaveURL(testURL, testLink)
	assert.Nil(t, err1)
}

func TestDBRepository_DBSaveURLError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewDBRepository(db)

	mock.ExpectExec(`INSERT INTO urlandlinks`).WithArgs(testURL, testLink).WillReturnError(errors.New(""))

	err1 := repo.DBSaveURL(testURL, testLink)
	assert.NotNil(t, err1)
}

func TestDBRepository_DBGetURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewDBRepository(db)
	rows := sqlmock.NewRows([]string{"url"}).AddRow(testURL)

	mock.ExpectQuery("SELECT url").WithArgs(testLink).WillReturnRows(rows)

	_, err1 := repo.DBGetURL(testLink)
	assert.Nil(t, err1)
}

func TestDBRepository_DBGetURLError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewDBRepository(db)

	mock.ExpectQuery("SELECT url").WithArgs(testLink).WillReturnError(errors.New(""))

	_, err1 := repo.DBGetURL(testLink)
	assert.NotNil(t, err1)
}
