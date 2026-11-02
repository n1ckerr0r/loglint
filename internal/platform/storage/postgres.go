package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/n1ckerr0r/loglint/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(ctx context.Context, databaseURL string) (*PostgresRepository, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}
	repo := &PostgresRepository{db: db}
	if err := repo.migrate(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return repo, nil
}

func (repo *PostgresRepository) migrate(ctx context.Context) error {
	_, err := repo.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS loglint_runs (
	id TEXT PRIMARY KEY,
	project TEXT NOT NULL,
	ref TEXT NOT NULL DEFAULT '',
	commit_sha TEXT NOT NULL DEFAULT '',
	status TEXT NOT NULL,
	report TEXT NOT NULL DEFAULT '',
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL
)`)
	return err
}

func (repo *PostgresRepository) CreateRun(ctx context.Context, run service.Run) error {
	_, err := repo.db.ExecContext(
		ctx,
		`INSERT INTO loglint_runs (id, project, ref, commit_sha, status, report, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		run.ID,
		run.Project,
		run.Ref,
		run.Commit,
		string(run.Status),
		run.Report,
		run.CreatedAt,
		run.UpdatedAt,
	)
	return err
}

func (repo *PostgresRepository) GetRun(ctx context.Context, id string) (service.Run, error) {
	var run service.Run
	var status string
	err := repo.db.QueryRowContext(
		ctx,
		`SELECT id, project, ref, commit_sha, status, report, created_at, updated_at
FROM loglint_runs WHERE id = $1`,
		id,
	).Scan(
		&run.ID,
		&run.Project,
		&run.Ref,
		&run.Commit,
		&status,
		&run.Report,
		&run.CreatedAt,
		&run.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return service.Run{}, ErrRunNotFound
	}
	if err != nil {
		return service.Run{}, err
	}
	run.Status = service.RunStatus(status)
	return run, nil
}

func (repo *PostgresRepository) UpdateRun(ctx context.Context, id string, status service.RunStatus, report string) error {
	result, err := repo.db.ExecContext(
		ctx,
		`UPDATE loglint_runs SET status = $1, report = $2, updated_at = $3 WHERE id = $4`,
		string(status),
		report,
		time.Now().UTC(),
		id,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrRunNotFound
	}
	return nil
}

func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}
