package main

import "LinkShortener/server"

func main() {
	srv := server.NewServer()
	srv.ListenAndServe()
}
