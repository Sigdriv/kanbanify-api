package model

import "time"

type Issue struct {
	KanbanID    string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	Variant     Variant   `json:"variant"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Status string

const (
	StatusBacklog    Status = "backlog"
	StatusInProgress Status = "inProgress"
	StatusDone       Status = "done"
)

type Variant string

const (
	VariantTask  Variant = "task"
	VariantBug   Variant = "bug"
	VariantChore Variant = "chore"
)
