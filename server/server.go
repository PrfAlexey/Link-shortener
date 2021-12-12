package server

import (
	"LinkShortener/pkg/handler"
	"LinkShortener/pkg/repository"
	"LinkShortener/pkg/service"
	"database/sql"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/labstack/echo"
	"log"
	"os"
)

const (
	dbConnect = "user=admin dbname=link-shortener password=4444 host=localhost port=5432 sslmode=disable"
)

//Server struct
type Server struct {
	e *echo.Echo
}

//NewServer creates a server
func NewServer() *Server {
	var server Server
	e := echo.New()

	if os.Getenv("DB") == "true" {

		db, err := sql.Open("pgx", dbConnect)
		if err != nil {
			log.Fatalln("cant parse config", err)
		}
		err = db.Ping() // вот тут будет первое подключение к базе
		if err != nil {
			log.Fatalln(err)
		}
		db.SetMaxOpenConns(10)

		repos := repository.NewDBRepository(db)
		services := service.NewService(nil, repos)
		handler := handler.NewHandler(services)
		handler.InitHandler(e)
		server.e = e
	} else {
		storage := make(map[string]string)
		repos := repository.NewRepository(storage)
		services := service.NewService(repos, nil)
		handler := handler.NewHandler(services)
		handler.InitHandler(e)
		server.e = e
	}

	return &server
}

//ListenAndServe starts the server
func (s Server) ListenAndServe() {
	s.e.Start(":8000")
}
