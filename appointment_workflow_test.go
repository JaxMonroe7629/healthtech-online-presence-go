package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNotifyOnCallPublishesWhenSomeoneIsOnline(t *testing.T) {
	published := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/realtime/presence/get/healthtech-ops" {
			w.Write([]byte(`{"ok":true,"data":{"users":[{"id":"nurse-7","status":"online"}]},"error":null,"metadata":{}}`))
			return
		}
		if r.URL.Path == "/v1/realtime/publish" && r.Method == http.MethodPost {
			published = true
			w.Write([]byte(`{"ok":true,"data":{},"error":null,"metadata":{}}`))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()
	sent, err := NotifyOnCall(&Client{BaseURL: server.URL, APIKey: "test", HTTP: server.Client()}, "healthtech-ops", Appointment{ID: "apt-1", Status: "ready"}, "clinic-ops", "appointment-apt-1")
	if err != nil || !sent || !published {
		t.Fatalf("sent=%v published=%v err=%v", sent, published, err)
	}
}
func TestNotifyOnCallSkipsEmptyRoster(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/v1/realtime/presence/get/") {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		w.Write([]byte(`{"ok":true,"data":{"users":[]},"error":null,"metadata":{}}`))
	}))
	defer server.Close()
	sent, err := NotifyOnCall(&Client{BaseURL: server.URL, APIKey: "test", HTTP: server.Client()}, "healthtech-ops", Appointment{ID: "apt-2", Status: "ready"}, "clinic-ops", "appointment-apt-2")
	if err != nil || sent {
		t.Fatalf("sent=%v err=%v", sent, err)
	}
}
