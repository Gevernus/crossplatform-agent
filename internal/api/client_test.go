package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/agents/status" {
			t.Errorf("Expected to request '/agents/status', got: %s", r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json header, got: %s", r.Header.Get("Content-Type"))
		}
		var statusUpdate StatusUpdate
		err := json.NewDecoder(r.Body).Decode(&statusUpdate)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}
		if statusUpdate.Status != "active" {
			t.Errorf("Expected status 'active', got: %s", statusUpdate.Status)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-agent", "test-password")
	err := client.SendStatus("active")
	if err != nil {
		t.Errorf("SendStatus returned an error: %v", err)
	}
}

func TestGetCommands(t *testing.T) {
	expectedCommands := []Command{
		{Action: "shutdown", Params: map[string]string{"delay": "30"}},
		{Action: "update", Params: map[string]string{"version": "1.2.0"}},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/agents/commands" {
			t.Errorf("Expected to request '/agents/commands', got: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedCommands)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-agent", "test-password")
	commands, err := client.GetCommands()
	if err != nil {
		t.Errorf("GetCommands returned an error: %v", err)
	}

	if len(commands) != len(expectedCommands) {
		t.Errorf("Expected %d commands, got %d", len(expectedCommands), len(commands))
	}

	for i, cmd := range commands {
		if cmd.Action != expectedCommands[i].Action {
			t.Errorf("Expected command action '%s', got '%s'", expectedCommands[i].Action, cmd.Action)
		}
		for k, v := range expectedCommands[i].Params {
			if cmd.Params[k] != v {
				t.Errorf("Expected param '%s' to be '%s', got '%s'", k, v, cmd.Params[k])
			}
		}
	}
}