package entity

type JobCompletion struct {
	JobID        string  `json:"jobId"`
	Status       string  `json:"status"` // "completed" or "failed"
	FinishedAt   int64   `json:"finishedAt"`
	ErrorMessage *string `json:"errorMessage,omitempty"`
}
