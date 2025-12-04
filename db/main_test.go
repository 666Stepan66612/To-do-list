package main

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// ============================================================================
// ТЕСТЫ ДЛЯ runMigrations
// ============================================================================

func TestRunMigrationsSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	// Users table
	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS users`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Tasks table
	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS tasks`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Add user_id column
	mock.ExpectExec(`DO \$\$`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Create index
	mock.ExpectExec(`CREATE INDEX IF NOT EXISTS idx_tasks_user_id`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = runMigrations(db)
	if err != nil {
		t.Errorf("runMigrations вернул ошибку: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Не выполнены ожидания mock: %v", err)
	}
}

func TestRunMigrationsUsersTableError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS users`).
		WillReturnError(sql.ErrConnDone)

	err = runMigrations(db)
	if err == nil {
		t.Error("runMigrations должен вернуть ошибку при ошибке создания users table")
	}
}

func TestRunMigrationsTasksTableError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS users`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS tasks`).
		WillReturnError(sql.ErrConnDone)

	err = runMigrations(db)
	if err == nil {
		t.Error("runMigrations должен вернуть ошибку при ошибке создания tasks table")
	}
}

func TestRunMigrationsAddColumnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS users`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS tasks`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec(`DO \$\$`).
		WillReturnError(sql.ErrConnDone)

	err = runMigrations(db)
	if err == nil {
		t.Error("runMigrations должен вернуть ошибку при ошибке добавления user_id column")
	}
}

func TestRunMigrationsCreateIndexError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS users`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS tasks`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec(`DO \$\$`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec(`CREATE INDEX IF NOT EXISTS idx_tasks_user_id`).
		WillReturnError(sql.ErrConnDone)

	err = runMigrations(db)
	if err == nil {
		t.Error("runMigrations должен вернуть ошибку при ошибке создания индекса")
	}
}
