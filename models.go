package main

// Webhook holds the information for a webhook from ClickUp.
type Webhook struct {
	ID     string `json:"webhook_id"`
	Event  string
	TaskID string `json:"task_id"`
}
