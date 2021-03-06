package buildserver

import (
	"errors"
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/auth"
	"github.com/concourse/atc/dbng"
)

type scopedHandlerFactory struct {
	logger lager.Logger
}

func NewScopedHandlerFactory(
	logger lager.Logger,
) *scopedHandlerFactory {
	return &scopedHandlerFactory{
		logger: logger,
	}
}

func (f *scopedHandlerFactory) HandlerFor(buildScopedHandler func(dbng.Build) http.Handler) http.HandlerFunc {
	logger := f.logger.Session("scoped-build-factory")

	return func(w http.ResponseWriter, r *http.Request) {
		build, ok := r.Context().Value(auth.BuildContextKey).(dbng.Build)
		if !ok {
			logger.Error("build-is-not-in-context", errors.New("build-is-not-in-context"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		buildScopedHandler(build).ServeHTTP(w, r)
	}
}
