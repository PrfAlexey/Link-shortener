package repository

import (
	"LinkShortener/pkg"
	mocks "LinkShortener/pkg/mocks"
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgmock"
	"github.com/jackc/pgproto3/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"strings"
	"testing"
	"time"
)

const DBConnect = "user=postgres dbname=postgres1 password=4444 host=localhost port=5432 sslmode=disable pool_max_conns=50"

var (
	testLink        = "1234567890"
	testURL         = "https://github.com/"
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

func setUpDB(t *testing.T) *pgxpool.Pool {
	script := &pgmock.Script{
		Steps: pgmock.AcceptUnauthenticatedConnRequestSteps(),
	}
	script.Steps = append(script.Steps, pgmock.ExpectMessage(&pgproto3.Query{String: ""}))
	script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.RowDescription{
		Fields: []pgproto3.FieldDescription{
			pgproto3.FieldDescription{
				Name:                 []byte("URL"),
				TableOID:             0,
				TableAttributeNumber: 0,
				DataTypeOID:          23,
				DataTypeSize:         60,
				TypeModifier:         -1,
				Format:               0,
			},
			pgproto3.FieldDescription{
				Name:                 []byte("Link"),
				TableOID:             0,
				TableAttributeNumber: 0,
				DataTypeOID:          23,
				DataTypeSize:         60,
				TypeModifier:         -1,
				Format:               0,
			},
		},
	}))
	script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.DataRow{
		Values: [][]byte{[]byte("1")},
	}))
	script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.CommandComplete{CommandTag: []byte("SELECT 1")}))
	script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.ReadyForQuery{TxStatus: 'I'}))
	script.Steps = append(script.Steps, pgmock.ExpectMessage(&pgproto3.Terminate{}))

	ln, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)
	defer ln.Close()

	serverErrChan := make(chan error, 1)
	go func() {
		defer close(serverErrChan)

		conn, err := ln.Accept()
		if err != nil {
			serverErrChan <- err
			return
		}
		defer conn.Close()

		err = conn.SetDeadline(time.Now().Add(time.Second))
		if err != nil {
			serverErrChan <- err
			return
		}

		err = script.Run(pgproto3.NewBackend(pgproto3.NewChunkReader(conn), conn))
		if err != nil {
			serverErrChan <- err
			return
		}
	}()

	parts := strings.Split(ln.Addr().String(), ":")
	host := parts[0]
	port := parts[1]
	connStr := fmt.Sprintf("sslmode=disable host=%s port=%s", host, port)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	pool, err := pgxpool.Connect(ctx, connStr)
	require.NoError(t, err)
	return pool
}

func TestDBRepository_DBCheckURLError(t *testing.T) {
	pool := setUpDB(t)
	h := NewDBRepository(pool)
	_, err := h.DBCheckURL(testInvalidURL)

	assert.NotNil(t, err)
}

func TestDBRepository_DBGetURLError(t *testing.T) {
	pool := setUpDB(t)
	h := NewDBRepository(pool)
	_, err := h.DBGetURL(testInvalidLink)

	assert.NotNil(t, err)
}

func TestDBRepository_DBSaveURLError(t *testing.T) {
	pool := setUpDB(t)
	h := NewDBRepository(pool)

	err := h.DBSaveURL(testURL, testLink)

	assert.NotNil(t, err)
}
