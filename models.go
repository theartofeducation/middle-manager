package main

// Webhook holds the information for a Webhook from ClickUp.
type Webhook struct {
	ID     string `json:"webhook_id"`
	Event  string
	TaskID string `json:"task_id"`
}

// Task holds the information for a Task on ClickUp.
type Task struct {
	ID string
}
