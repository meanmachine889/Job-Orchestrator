package orchestrator

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type JobCreate struct {
	ID             uuid.UUID       `json:"job_id"`
	Type           string          `json:"type"`
	Payload        json.RawMessage `json:"payload"`
	Status         string          `json:"status"`
	RetryCount     int             `json:"retry_count"`
	MaxRetries     int             `json:"max_retries"`
	TimeoutSeconds int             `json:"timeout_seconds"`
}

type Client struct {
	baseUrl string
}

func New(baseUrl string) *Client {
	return &Client{baseUrl: baseUrl}
}

func (c *Client) RegisterWorker(hostname string) (string, error) {
	body, _ := json.Marshal(map[string]string{
		"hostname": hostname,
	})

	resp, err := http.Post(c.baseUrl+"/workers/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		ID string `json:"id"`
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return "", err
	}

	return res.ID, nil
}

func (c *Client) SendHeartbeat(workerID string) error {
	body, _ := json.Marshal(map[string]string{
		"id": workerID,
	})

	_, err := http.Post(c.baseUrl+"/workers/heartbeat", "application/json", bytes.NewBuffer(body))
	return err
}

func (c *Client) FetchJob(workerID string) (*JobCreate, error) {
	body, _ := json.Marshal(map[string]string{
		"worker_id": workerID,
	})

	resp, err := http.Post(c.baseUrl+"/jobs/next", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	var job JobCreate
	err = json.NewDecoder(resp.Body).Decode(&job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (c *Client) ReportJobResult(jobID string, status string, errMsg string) error {
	body, _ := json.Marshal(map[string]string{
		"job_id": jobID,
		"status": status,
		"error":  errMsg,
	})

	_, err := http.Post(c.baseUrl+"/jobs/report", "application/json", bytes.NewBuffer(body))
	return err
}
