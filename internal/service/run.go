package service

import (
	"context"
	"time"
)

type RunStatus string

const (
	RunStatusQueued    RunStatus = "queued"
	RunStatusRunning   RunStatus = "running"
	RunStatusCompleted RunStatus = "completed"
	RunStatusFailed    RunStatus = "failed"
)

type Run struct {
	ID        string    `json:"id"`
	Project   string    `json:"project"`
	Ref       string    `json:"ref,omitempty"`
	Commit    string    `json:"commit,omitempty"`
	Status    RunStatus `json:"status"`
	Report    string    `json:"report,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RunRequest struct {
	Project string `json:"project"`
	Ref     string `json:"ref"`
	Commit  string `json:"commit"`
}

type Job struct {
	RunID   string `json:"run_id"`
	Project string `json:"project"`
	Ref     string `json:"ref,omitempty"`
	Commit  string `json:"commit,omitempty"`
}

type RunRepository interface {
	CreateRun(context.Context, Run) error
	GetRun(context.Context, string) (Run, error)
	UpdateRun(context.Context, string, RunStatus, string) error
	Close() error
}

type Queue interface {
	Publish(context.Context, Job) error
	Consume(context.Context) (<-chan Job, error)
	Close() error
}

type StatusCache interface {
	SetStatus(context.Context, string, RunStatus) error
	Close() error
}
