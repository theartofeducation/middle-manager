package main

// ClickUpStatus represents a status in ClickUp.
type ClickUpStatus string

const clickUpStatusReadyForDevelopment ClickUpStatus = "ready for development"

// Webhook holds the information for a Webhook from ClickUp.
type Webhook struct {
	ID     string `json:"webhook_id"`
	Event  string
	TaskID string `json:"task_id"`
}

// Task holds the information for a Task on ClickUp.
type Task struct {
	ID     string
	Name   string
	Status TaskStatus
}

// TaskStatus holds the information for a Task's status on ClickUp.
type TaskStatus struct {
	Status ClickUpStatus
}
