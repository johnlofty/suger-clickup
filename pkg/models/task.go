package models

import (
	"errors"
	"time"
)

type Status string

const (
	Open       Status = "open"
	InProgress Status = "in progress"
	Closed     Status = "closed"
	Blocked    Status = "blocked"
	Archived   Status = "archived"

	MaxDescription = 255

	OrgNameField = "81a95d61-ec0e-4bff-b16f-4b5609d6945a"
	OrgIDField   = "90a98593-c851-460f-8c3a-41b0360a9b67"
	CreatorField = "4e4abba4-27d1-49af-bf1a-f7ba83564fb9"
)

type Task struct {
	ID            string             `json:"id"`
	Name          string             `json:"name,omitempty"`
	Description   string             `json:"description,omitempty"`
	Tags          []string           `json:"tags,omitempty"`
	Status        string             `json:"status,omitempty"`
	Priority      int                `json:"priority,omitempty"`
	DueDate       int64              `json:"due_date,omitempty"`
	DueDateTime   bool               `json:"due_date_time,omitempty"`
	TimeEstimate  int32              `json:"time_estimate,omitempty"`
	StartDate     int64              `json:"start_date,omitempty"`
	StartDateTime bool               `json:"start_date_time,omitempty"`
	CustomFields  []TaskCustomFields `json:"custom_fields,omitempty"`
}

func (t *Task) AddOrgId(orgID string) *Task {
	t.CustomFields = append(t.CustomFields, TaskCustomFields{
		ID:    OrgIDField,
		Value: orgID,
	})
	return t
}
func (t *Task) AddOrgName(orgName string) *Task {
	t.CustomFields = append(t.CustomFields, TaskCustomFields{
		ID:    OrgNameField,
		Value: orgName,
	})
	return t
}
func (t *Task) AddCreator(creator string) *Task {
	t.CustomFields = append(t.CustomFields, TaskCustomFields{
		ID:    CreatorField,
		Value: creator,
	})
	return t
}

type TaskCustomFields struct {
	ID    string `json:"id,omitempty"`
	Value string `json:"value,omitempty"`
}

type CreateTaskRequest struct {
	Name        string
	Description string
	DueTime     int64
}

func (r CreateTaskRequest) Validate() error {
	if len(r.Name) == 0 {
		return errors.New("invalid name")
	}
	if len(r.Description) == 0 {
		return errors.New("invalid description")
	}
	if len(r.Description) > MaxDescription {
		return errors.New("description too long")
	}
	return nil
}

type Ticket struct {
	TicketID  string    `db:"ticket_id" json:"ticket_id,omitempty"`
	OrgID     int32     `db:"org_id" json:"org_id,omitempty"`
	UserID    int32     `db:"user_id" json:"user_id,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at,omitempty"`
}

type ClickupTask struct {
	ID          string
	Name        string
	Description string
	Status      ClickupTaskStatus
	DateCreated string
	DateUpdated string
}

type ClickupTaskStatus struct {
	ID     string `json:"id,omitempty"`
	Status string `json:"status,omitempty"`
	Color  string `json:"color,omitempty"`
	Type   string `json:"type,omitempty"`
}
