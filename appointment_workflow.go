package main

import "fmt"

type Appointment struct {
	ID     string
	Status string
}

func NotifyOnCall(client *Client, channel string, appointment Appointment, accountID, requestID string) (bool, error) {
	presence, err := client.Presence(channel)
	if err != nil {
		return false, err
	}
	if len(presence.Users) == 0 {
		return false, nil
	}
	data := map[string]string{"appointment_id": appointment.ID, "status": appointment.Status, "message": fmt.Sprintf("Appointment %s is ready for team review", appointment.ID)}
	if err := client.Publish(channel, "appointment.updated", data, accountID, requestID); err != nil {
		return false, err
	}
	return true, nil
}
