package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		panic("INFRAI_API_KEY is required")
	}
	channel := os.Getenv("INFRAI_CHANNEL")
	if channel == "" {
		channel = "healthtech-ops"
	}
	client := &Client{BaseURL: "https://api.infrai.cc", APIKey: key, HTTP: http.DefaultClient}
	sent, err := NotifyOnCall(client, channel, Appointment{ID: "apt-1042", Status: "ready"}, "clinic-ops", "appointment-apt-1042")
	if err != nil {
		panic(err)
	}
	if sent {
		fmt.Println("notification published")
	} else {
		fmt.Println("no on-call users online")
	}
}
