package orchestrator

import (
	"bytes"
	"encoding/json"
	"net/http"
)

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
