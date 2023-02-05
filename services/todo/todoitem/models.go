package todoitem

import "time"

type TodoItem struct {
	Id        *string    `json:"id,omitempty"`
	Created   *time.Time `json:"created,omitempty"`
	Updated   *time.Time `json:"updatedomitempty"`
	Deleted   *bool      `json:"deleted,omitempty"`
	Completed *bool      `json:"completed,omitempty"`
	Summary   *string    `json:"summary,omitempty"`
}
