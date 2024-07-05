package httperrors

import (
	"caching/util"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

var (
	ErrNotFound = ErrResponse{
		HTTPStatusCode: http.StatusNotFound,
		ErrorResponse: ErrorResponse{
			StatusText: util.Ptr("Resource not found"),
		},
	}

	ErrInvalidRequest = ErrResponse{
		HTTPStatusCode: http.StatusBadRequest,
		ErrorResponse: ErrorResponse{
			StatusText: util.Ptr("Invalid request"),
		},
	}

	ErrInternal = ErrResponse{
		HTTPStatusCode: http.StatusInternalServerError,
		ErrorResponse: ErrorResponse{
			StatusText: util.Ptr("Internal error"),
		},
	}

	ErrUnauthorize = ErrResponse{
		HTTPStatusCode: http.StatusUnauthorized,
		ErrorResponse: ErrorResponse{
			StatusText: util.Ptr("Unauthorize error"),
		},
	}
)

type ErrorResponse struct {
	AppCode    *int    `json:"appCode,omitempty"`
	ErrorText  *string `json:"errorText,omitempty"`
	StatusText *string `json:"statusText,omitempty"`
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	ErrorResponse
}

func (e ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func (e ErrResponse) Error() string {
	return util.SafeValue(e.ErrorText)
}

func (e ErrResponse) WithErrorText(err error) render.Renderer {
	e.ErrorText = util.Ptr(err.Error())
	e.Err = err
	return e
}

func NewFailureRender(err error) render.Renderer {
	errRender := ErrInternal
	if errors.As(err, &errRender) {
		errRender.ErrorText = util.Ptr(err.Error())
		return errRender
	}
	return errRender.WithErrorText(err)
}
