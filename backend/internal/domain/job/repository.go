package job

type Job struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenantId"`
	Type        string `json:"type"`
	Payload     string `json:"payload"`
	Status      string `json:"status"` // "pending", "running", "completed", "failed"
	Attempts    int    `json:"attempts"`
	MaxAttempts int    `json:"maxAttempts"`
	RunAt       string `json:"runAt"`
	LastError   string `json:"lastError,omitempty"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type Repository interface {
	Enqueue(j Job) (Job, error)
	FetchPending(limit int) ([]Job, error)
	MarkRunning(id string) error
	MarkCompleted(id string) error
	MarkFailed(id string, errMessage string, retryAfterSeconds int) error
	CountByStatus() (map[string]int, error)
}

