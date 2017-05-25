package buildserver

import (
	"encoding/json"
	"net/http"

	"github.com/concourse/atc/api/present"
	"github.com/concourse/atc/dbng"
)

func (s *Server) CreateRebuild(build dbng.Build) http.Handler {
	logger := s.logger.Session("create-rebuild")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := build.Reset()
		if err != nil {
			logger.Error("failed-to-reset-build", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		scheduled, err := build.Schedule()
		if err != nil {
			logger.Error("failed-to-schedule-build", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		if !scheduled {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		//p, ok, err := build.Pipeline()
		//		if !ok {
		//			if err != nil {
		//				w.WriteHeader(http.StatusInternalServerError)
		//				w.Write([]byte(err.Error()))
		//				return
		//			}
		//			w.WriteHeader(http.StatusNotFound)
		//			return
		//		}
		//
		//		inputs, outputs, err := build.Resources()
		//		if err != nil {
		//			w.WriteHeader(http.StatusInternalServerError)
		//			w.Write([]byte(err.Error()))
		//			return
		//		}
		//
		//		job, ok, err := p.Job(build.JobName())
		//		if !ok {
		//			if err != nil {
		//				w.WriteHeader(http.StatusInternalServerError)
		//				w.Write([]byte(err.Error()))
		//				return
		//			}
		//			w.WriteHeader(http.StatusNotFound)
		//			return
		//		}
		//
		//		engineBuild, err := s.engine.CreateBuild(hLog, build, plan)
		//		if err != nil {
		//			hLog.Error("failed-to-start-build", err)
		//			w.WriteHeader(http.StatusInternalServerError)
		//			return
		//		}
		//
		//		go engineBuild.Resume(logger)

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(present.Build(build))
	})
}
