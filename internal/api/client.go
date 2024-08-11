package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	baseURL    string
	uuid       string
	deviceID   string
	httpClient *http.Client
}

type StatusUpdate struct {
	UUID           string                 `json:"uuid"`
	HardwareID     string                 `json:"hardware_id"`
	AdditionalData map[string]interface{} `json:"additional_data"`
}

type CommandsRequest struct {
	UUID       string `json:"uuid"`
	HardwareID string `json:"hardware_id"`
}

type Command struct {
	Command string                 `json:"command"`
	Payload map[string]interface{} `json:"payload"`
}

type State struct {
	Version     string `json:"version"`
	Status      string `json:"status"`
	Error       string `json:"error"`
	ErrorUILink string `json:"error_ui_link"`
}

type GetCommandsResponse struct {
	State      State     `json:"state"`
	TotalCount int       `json:"total_count"`
	Commands   []Command `json:"commands"`
}

func (c Command) String() string {
	var params []string
	for key, value := range c.Payload {
		params = append(params, fmt.Sprintf("%s: %s", key, value))
	}

	paramsStr := strings.Join(params, ", ")
	if paramsStr != "" {
		return fmt.Sprintf("Action: %s, Params: {%s}", c.Command, paramsStr)
	}
	return fmt.Sprintf("Action: %s", c.Command)
}

func NewClient(baseURL, uuid, deviceID string) *Client {
	return &Client{
		baseURL:  baseURL,
		uuid:     uuid,
		deviceID: deviceID,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) SendStatus() error {
	url := fmt.Sprintf("%s/agents/status", c.baseURL)
	log.Debug("Sending ok status to API:", url)
	statusUpdate := StatusUpdate{UUID: c.uuid, HardwareID: c.deviceID}
	payload, err := json.Marshal(statusUpdate)
	if err != nil {
		return err
	}
	log.Debugf("Payload being sent: %s", string(payload))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Log the response body for more details on the error
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Debugf("Failed to send status: unexpected status code: %d, response: %s", resp.StatusCode, string(bodyBytes))
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) GetCommands() ([]Command, error) {
	url := fmt.Sprintf("%s/agents/commands", c.baseURL)
	commandsRequest := CommandsRequest{
		UUID:       c.uuid,
		HardwareID: c.deviceID,
	}

	payload, err := json.Marshal(commandsRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response GetCommandsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if response.State.Status != "200" {
		return nil, fmt.Errorf("API returned error: %s", response.State.Error)
	}

	return response.Commands, nil
}
