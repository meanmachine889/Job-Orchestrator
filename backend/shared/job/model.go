package job

type Status string

const (
	Pending  Status = "PENDING"
	Running  Status = "RUNNING"
	Success  Status = "SUCCESS"
	Failed   Status = "FAILED"
	Retrying Status = "RETRYING"
	Dead     Status = "DEAD"
)

