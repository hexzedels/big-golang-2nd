package handler

import (
	"context"
	"scheduler/scheduler/internal/cases"
	"scheduler/scheduler/internal/input/http/gen"
)

var _ gen.StrictServerInterface = (*Server)(nil)

type Server struct {
	schedulerCase *cases.SchedulerCase
}

func NewServer(schCase *cases.SchedulerCase) *Server {
	return &Server{
		schedulerCase: schCase,
	}
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
	panic("not implemented") // TODO: Implement
}

// Delete a job
// (DELETE /jobs/{job_id})
func (r *Server) DeleteJobsJobId(ctx context.Context, request gen.DeleteJobsJobIdRequestObject) (gen.DeleteJobsJobIdResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Get job details
// (GET /jobs/{job_id})
func (r *Server) GetJobsJobId(ctx context.Context, request gen.GetJobsJobIdRequestObject) (gen.GetJobsJobIdResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Get job executions
// (GET /jobs/{job_id}/executions)
func (r *Server) GetJobsJobIdExecutions(ctx context.Context, request gen.GetJobsJobIdExecutionsRequestObject) (gen.GetJobsJobIdExecutionsResponseObject, error) {
	panic("not implemented") // TODO: Implement
}
