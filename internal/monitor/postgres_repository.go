package monitor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type postgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a fully connected postgres tracking repository
func NewPostgresRepository(dsn string) (Repository, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Backoff mechanism to ensure backend explicitly waits for dockerized database readiness
	var pingErr error
	for i := 0; i < 10; i++ {
		pingErr = db.Ping()
		if pingErr == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if pingErr != nil {
		return nil, fmt.Errorf("failed to reach database after retries: %w", pingErr)
	}

	if err := initSchema(db); err != nil {
		return nil, err
	}

	return &postgresRepository{db: db}, nil
}

func initSchema(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS monitors (
		id TEXT PRIMARY KEY,
		user_id TEXT,
		url TEXT,
		interval BIGINT,
		last_checked TIMESTAMP,
		status_code INT,
		response_time BIGINT,
		is_healthy BOOLEAN,
		ai_explanation TEXT,
		is_running BOOLEAN
	)`
	_, err := db.Exec(query)
	return err
}

func (r *postgresRepository) Add(ctx context.Context, m *Monitor) error {
	query := `
	INSERT INTO monitors (id, user_id, url, interval, last_checked, status_code, response_time, is_healthy, ai_explanation, is_running)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		m.ID, m.UserID, m.URL, m.Interval, m.LastChecked, m.StatusCode, m.ResponseTime, m.IsHealthy, m.AIExplanation, m.IsRunning,
	)
	return err
}

func (r *postgresRepository) List(ctx context.Context, userID string) ([]*Monitor, error) {
	query := `SELECT id, user_id, url, interval, last_checked, status_code, response_time, is_healthy, ai_explanation, is_running FROM monitors WHERE user_id = $1`
	return r.queryMonitors(ctx, query, userID)
}

func (r *postgresRepository) GetAll(ctx context.Context) ([]*Monitor, error) {
	query := `SELECT id, user_id, url, interval, last_checked, status_code, response_time, is_healthy, ai_explanation, is_running FROM monitors`
	return r.queryMonitors(ctx, query)
}

func (r *postgresRepository) queryMonitors(ctx context.Context, query string, args ...interface{}) ([]*Monitor, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*Monitor
	for rows.Next() {
		var m Monitor
		if err := rows.Scan(
			&m.ID, &m.UserID, &m.URL, &m.Interval, &m.LastChecked, &m.StatusCode, &m.ResponseTime, &m.IsHealthy, &m.AIExplanation, &m.IsRunning,
		); err != nil {
			return nil, err
		}
		result = append(result, &m)
	}
	return result, rows.Err()
}

func (r *postgresRepository) UpdateStatus(ctx context.Context, id string, lastChecked time.Time, statusCode int, responseTime time.Duration, isHealthy bool, aiExplanation string) error {
	query := `
	UPDATE monitors 
	SET last_checked = $1, status_code = $2, response_time = $3, is_healthy = $4, ai_explanation = $5
	WHERE id = $6
	`
	res, err := r.db.ExecContext(ctx, query, lastChecked, statusCode, responseTime, isHealthy, aiExplanation, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("monitor not found")
	}

	return nil
}

func (r *postgresRepository) SetRunning(ctx context.Context, id string, isRunning bool) error {
	query := `UPDATE monitors SET is_running = $1 WHERE id = $2`
	res, err := r.db.ExecContext(ctx, query, isRunning, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("monitor not found")
	}

	return nil
}
