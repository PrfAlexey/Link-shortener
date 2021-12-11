package server

import (
	"LinkShortener/pkg/handler"
	"LinkShortener/pkg/repository"
	"LinkShortener/pkg/service"
	"database/sql"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/labstack/echo"
	"log"
)

const (
	dbConnect = "user=postgres dbname=postgres1 password=4444 host=localhost port=5432 sslmode=disable"
	dataBase  = true
)

//Server ...
type Server struct {
	e *echo.Echo
}

//NewServer ...
func NewServer() *Server {
	var server Server
	e := echo.New()
	if dataBase {

		db, err := sql.Open("pgx", dbConnect)
		if err != nil {
			log.Fatalln("cant parse config", err)
		}
		err = db.Ping() // вот тут будет первое подключение к базе
		if err != nil {
			log.Fatalln(err)
		}
		db.SetMaxOpenConns(10)
		//pool, err := pgxpool.Connect(context.Background(), dbConnect)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//err = pool.Ping(context.Background())
		//if err != nil {
		//	log.Fatal(err)
		//}
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

//ListenAndServe ...
func (s Server) ListenAndServe() {
	s.e.Start(":8000")
}
