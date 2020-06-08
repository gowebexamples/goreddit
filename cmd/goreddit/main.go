package main

import (
	"log"
	"net/http"

	"github.com/gowebexamples/goreddit/postgres"
	"github.com/gowebexamples/goreddit/web"
)

func main() {
	store, err := postgres.NewStore("postgres://postgres:secret@localhost/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	csrfKey := []byte("01234567890123456789012345678901")
	h := web.NewHandler(store, csrfKey)
	http.ListenAndServe(":3000", h)
}
