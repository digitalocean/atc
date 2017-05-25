package buildserver

import (
	"encoding/json"
	"net/http"

	"code.cloudfoundry.org/lager"

	"github.com/concourse/atc"
	"github.com/concourse/atc/api/present"
	"github.com/concourse/atc/dbng"
)

func (s *Server) CreateRebuild(build dbng.Build) http.Handler {
	hLog := s.logger.Session("create-rebuild")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var plan atc.Plan
		err := json.NewDecoder(r.Body).Decode(&plan)
		if err != nil {
			hLog.Info("malformed-request", lager.Data{"error": err.Error()})
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = build.Reset()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		engineBuild, err := s.engine.CreateBuild(hLog, build, plan)
		if err != nil {
			hLog.Error("failed-to-start-build", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		engineBuild.Resume(hLog)

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(present.Build(build))
	})
}
