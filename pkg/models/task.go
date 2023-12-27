package models

import "encoding/json"

type Status string

const (
	Open       Status = "open"
	InProgress Status = "in progress"
	Closed     Status = "closed"
	Blocked    Status = "blocked"
	Archived   Status = "archived"
)

type Task struct {
	Name          string   `json:"name,omitempty"`
	Description   string   `json:"description,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	Status        Status   `json:"status,omitempty"`
	Priority      int      `json:"priority,omitempty"`
	DueDate       int64    `json:"due_date,omitempty"`
	DueDateTime   bool     `json:"due_date_time,omitempty"`
	TimeEstimate  int32    `json:"time_estimate,omitempty"`
	StartDate     int64    `json:"start_date,omitempty"`
	StartDateTime bool     `json:"start_date_time,omitempty"`
}

func (t *Task) Unmarshal() {

}

func (t *Task) Scan(data []byte) error {
	return json.Unmarshal(data, t)
}
