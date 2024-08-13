package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
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
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-agent", "test-password")
	err := client.SendStatus()
	if err != nil {
		t.Errorf("SendStatus returned an error: %v", err)
	}
}

func TestGetCommands(t *testing.T) {
	expectedResponse := GetCommandsResponse{
		State: State{
			Version:     "1.0", // Include this if you need the version in your test
			Status:      "200",
			Error:       "",
			ErrorUILink: "", // Include this if necessary
		},
		Commands: []Command{
			{Command: "shutdown"},
			{Command: "update"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/agents/commands" {
			t.Errorf("Expected to request '/agents/commands', got: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-agent", "test-password")
	commands, err := client.GetCommands()
	if err != nil {
		t.Fatalf("GetCommands returned an error: %v", err)
	}

	if !reflect.DeepEqual(commands, expectedResponse.Commands) {
		t.Errorf("Expected commands %+v, got %+v", expectedResponse.Commands, commands)
	}
}
