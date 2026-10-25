package repository

import (
	"GoWatch/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var NF = errors.New("Not found")

type PostgresRepository interface {
	Create(*domain.Monitor) error
	GetAllMon() ([]*domain.Monitor, error)
	DelMon(int64) error
	GetMon(int64) (*domain.Monitor, error)
	GetDueMonitors(ctx context.Context) ([]*domain.Monitor, error)
	Close() error
}

type PostRepo struct {
	db *sql.DB
}

func (r *PostRepo) SaveAndUpdate(ctx context.Context, check domain.CheckResult) error {
	queryInsert := `INSERT INTO checks (monitor_id, status_code, response_time_ms, error_msg, checked_at)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err // Возможно стоит заменить ошибку БД на свою
	}

	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, queryInsert, check.MonitorID, check.StatusCode, check.ResponseTimeMs, check.ErrorMsg,
		check.CheckedAt).Scan(&check.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostRepo) GetDueMonitors(ctx context.Context) ([]*domain.Monitor, error) {
	res := make([]*domain.Monitor, 0)

	query := `WITH due AS (
	SELECT id FROM monitors
	WHERE next_check_at <= NOW()
	FOR UPDATE SKIP LOCKED)
	UPDATE monitors m SET next_check_at = NOW() + (m.interval_sec|| ' seconds')::interval
	FROM due WHERE m.id = due.id
	RETURNING m.id, m.name, m.url, m.interval_sec, m.created_at, m.next_check_at
	`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var info domain.Monitor

		err = rows.Scan(&info.ID, &info.Name, &info.URL, &info.IntervalSec, &info.CreatedAt, &info.NextCheckAt)
		if err != nil {
			log.Printf("scanning error: %v", err)
			return nil, fmt.Errorf("internal error")
		}

		res = append(res, &info)
	}

	if rows.Err() != nil {
		log.Println(rows.Err())
		return nil, fmt.Errorf("internal error")
	}

	return res, tx.Commit()
}

func NewPostgresRepo(dsn string) (*PostRepo, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("cant open database")
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("cant connect to database")
	}

	return &PostRepo{db: db}, nil
}

func (r *PostRepo) Close() error {
	return r.db.Close()
}

func (r *PostRepo) Create(data *domain.Monitor) error {
	query := `INSERT INTO monitors (name, url, interval_sec, next_check_at)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at
	`

	return r.db.QueryRow(query, data.Name, data.URL, data.IntervalSec,
		time.Now().Add(time.Duration(data.IntervalSec)*time.Second)).Scan(&data.ID, &data.CreatedAt)
}

func (r *PostRepo) GetAllMon() ([]*domain.Monitor, error) {
	query := `SELECT id, name, url, interval_sec, created_at FROM monitors`

	res := make([]*domain.Monitor, 0)

	rows, err := r.db.Query(query)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("internal error")
	}

	defer rows.Close()

	for rows.Next() {
		var monitor domain.Monitor

		err := rows.Scan(&monitor.ID, &monitor.Name, &monitor.URL, &monitor.IntervalSec, &monitor.CreatedAt)
		if err != nil {
			log.Printf("scanning error: %v", err)
			return nil, fmt.Errorf("internal error")
		}
		res = append(res, &monitor)
	}

	if rows.Err() != nil {
		log.Println(rows.Err())
		return nil, fmt.Errorf("internal error")
	}

	return res, nil
}

func (r *PostRepo) DelMon(id int64) error {
	query := `DELETE FROM monitors WHERE id = $1`

	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	deletedStr, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if deletedStr == 0 {
		return NF
	}

	return nil
}

func (r *PostRepo) GetMon(id int64) (*domain.Monitor, error) {
	var info domain.Monitor
	query := `SELECT id, name, url, interval_sec, created_at FROM monitors WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(&info.ID, &info.Name, &info.URL, &info.IntervalSec, &info.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NF
		}

		return nil, err
	}

	return &info, nil
}
