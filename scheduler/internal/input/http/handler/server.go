package handler

import (
	"context"
	"net/http"
	"scheduler/scheduler/internal/cases"
	"scheduler/scheduler/internal/input/http/gen"

	"github.com/go-chi/chi/v5"
)

var _ gen.StrictServerInterface = (*Server)(nil)
var _ http.Handler = (*Server)(nil)

type Server struct {
	schedulerCase *cases.SchedulerCase
	router        *chi.Mux
}

func NewServer(schCase *cases.SchedulerCase) *Server {
	s := &Server{
		schedulerCase: schCase,
		router:        chi.NewRouter(),
	}

	strictHandler := gen.NewStrictHandler(s, nil)

	gen.HandlerFromMux(strictHandler, s.router)

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Create a new job
// (POST /jobs)
func (r *Server) PostJobs(ctx context.Context, request gen.PostJobsRequestObject) (gen.PostJobsResponseObject, error) {
	jobID, err := r.schedulerCase.Create(ctx, toEntityJob(request.Body))
	if err != nil {
		return nil, err // 500
	}

	return gen.PostJobs201JSONResponse(jobID), nil
}

// List jobs
// (GET /jobs)
func (r *Server) GetJobs(ctx context.Context, request gen.GetJobsRequestObject) (gen.GetJobsResponseObject, error) {
	// TODO: вернуть задание по id
	return gen.GetJobs200JSONResponse{}, nil

}

// Delete a job
// (DELETE /jobs/{job_id})
func (r *Server) DeleteJobsJobId(ctx context.Context, request gen.DeleteJobsJobIdRequestObject) (gen.DeleteJobsJobIdResponseObject, error) {
	// TODO: удалить задание
	return gen.DeleteJobsJobId204Response{}, nil
}

// Get job details
// (GET /jobs/{job_id})
func (r *Server) GetJobsJobId(ctx context.Context, request gen.GetJobsJobIdRequestObject) (gen.GetJobsJobIdResponseObject, error) {
	// TODO: реализовать поиск работы по идентификатору
	return gen.GetJobsJobId200JSONResponse{}, nil
}

// Get job executions
// (GET /jobs/{job_id}/executions)
func (r *Server) GetJobsJobIdExecutions(ctx context.Context, request gen.GetJobsJobIdExecutionsRequestObject) (gen.GetJobsJobIdExecutionsResponseObject, error) {
	// реализовать поиск выполненных задач
	return gen.GetJobsJobIdExecutions200JSONResponse{}, nil
}
