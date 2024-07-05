package main

import (
	getuser "caching/features/getusers"
	"caching/redis"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
)

func main() {
	if err := setupApp(); err != nil {
		panic(errors.Wrap(err, "failed to setup app"))
	}

	// init api
	r := buildRouter()

	http.ListenAndServe(":3000", r)
	// add caching handler
}

func buildRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(time.Second * 30))
	addController(r)
	return r
}

func addController(r chi.Router) {
	r.Get("/users/{id}", getuser.GeUser)
}

func setupApp() error {
	if err := redis.SetUpdateRedis(); err != nil {
		return errors.Wrap(err, "failed setup cache")
	}
	return nil
}
