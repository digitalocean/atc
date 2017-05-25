package buildserver

import (
	"encoding/json"
	"net/http"

	"github.com/concourse/atc/api/present"
	"github.com/concourse/atc/dbng"
)

func (s *Server) CreateRebuild(build dbng.Build) http.Handler {
	hLog := s.logger.Session("create-rebuild")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := build.Reset()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//go func() {
		//			nextPendingBuilds, err := s.Pipeline.GetPendingBuildsForJob(build.JobName())
		//			if err != nil {
		//				logger.Error("failed-to-get-next-pending-build-for-job", err)
		//				return
		//			}
		//
		//			err = s.BuildStarter.TryStartPendingBuildsForJob(logger, jobConfig, resourceConfigs, resourceTypes, nextPendingBuilds)
		//			if err != nil {
		//				logger.Error("failed-to-start-next-pending-build-for-job", err, lager.Data{"job-name": jobConfig.Name})
		//				return
		//			}
		//		}()

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(present.Build(build))
	})
}
