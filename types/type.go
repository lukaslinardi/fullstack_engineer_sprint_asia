package types

import "time"

type TaskRequest struct {
	TaskName     string `json:"task_name"`
	ParentId     *int   `json:"parent_id"`
	DeadlineTask string `json:"deadline_task"`
}

type TaskResponse struct {
	ID           int       `json:"id" db:"id"`
	TaskName     string    `json:"task_name" db:"task_name"`
	TaskStatus   int       `json:"task_status" db:"task_status"`
	ParentId     *int      `json:"parent_id" db:"parent_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	DeadlineTask time.Time `json:"deadline_task" db:"deadline_task"`
	Percentage   *float32  `json:"completion_percentage" json:"completion_percentage"`
}

type TaskUpdate struct {
	ID           int
	TaskName     string `json:"task_name" db:"task_name"`
	DeadlineTask string `json:"deadline_task"`
	ParentId     *int   `json:"parent_id" db:"parent_id"`
}
