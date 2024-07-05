package getuser

import (
	"caching/features/getusers/internal"
	"caching/httperrors"
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type service interface {
	GetUsers(ctx context.Context, id string) (*internal.User, error)
}

func newService() (service, error) {
	return internal.NewService()
}

func GeUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	svc, err := newService()
	if err != nil {
		render.Render(w, r, httperrors.NewFailureRender(err))
		return
	}

	u, err := svc.GetUsers(r.Context(), id)
	if err != nil {
		render.Render(w, r, httperrors.NewFailureRender(err))
		return
	}
	switch u.CacheHit {
	case true:
		w.Header().Add("CacheHit", "true")
	case false:
		w.Header().Add("CacheHit", "false")
	}

	render.JSON(w, r, u)
}
