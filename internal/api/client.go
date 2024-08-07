package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	agentID    string
	password   string
	httpClient *http.Client
}

type StatusUpdate struct {
	Status string `json:"status"`
}

type Command struct {
	Action string            `json:"action"`
	Params map[string]string `json:"params,omitempty"`
}

func NewClient(baseURL, agentID, password string) *Client {
	return &Client{
		baseURL:  baseURL,
		agentID:  agentID,
		password: password,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) SendStatus(status string) error {
	log.Debug("Sending status to API:", status)
	url := fmt.Sprintf("%s/agents/status", c.baseURL)
	statusUpdate := StatusUpdate{Status: status}
	payload, err := json.Marshal(statusUpdate)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.agentID, c.password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) GetCommands() ([]Command, error) {
	url := fmt.Sprintf("%s/agents/commands", c.baseURL)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.agentID, c.password)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var commands []Command
	if err := json.NewDecoder(resp.Body).Decode(&commands); err != nil {
		return nil, err
	}

	return commands, nil
}
