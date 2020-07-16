package main

import (
	"net/http"
	"time"

	"github.com/summerKK/go-code-snippet-library/blog-service/internal/routers"
)

func main() {
	router := routers.NewRouter()

	s := &http.Server{
		Addr:           ":8008",
		Handler:        router,
		TLSConfig:      nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
