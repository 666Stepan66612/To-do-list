package models

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

// ============================================================================
// ТЕСТЫ ДЛЯ NewTaskRepository
// ============================================================================

func TestNewTaskRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)
	if repo == nil {
		t.Fatal("NewTaskRepository вернул nil")
	}
	if repo.DB != db {
		t.Error("DB не установлена правильно")
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ CreateTask
// ============================================================================

func TestCreateTaskSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	task := &Task{
		UserID: 1,
		Name:   "Test Task",
		Text:   "Test Description",
	}

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "text", "complete", "create_time", "complete_at"}).
		AddRow(1, 1, "Test Task", "Test Description", false, now, nil)

	mock.ExpectQuery(`INSERT INTO tasks`).
		WithArgs(1, "Test Task", "Test Description").
		WillReturnRows(rows)

	err = repo.CreateTask(task)
	if err != nil {
		t.Errorf("CreateTask вернул ошибку: %v", err)
	}

	if task.ID != 1 {
		t.Errorf("Ожидался ID 1, получен %d", task.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Не выполнены ожидания mock: %v", err)
	}
}

func TestCreateTaskError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	task := &Task{
		UserID: 1,
		Name:   "Test Task",
		Text:   "Test Description",
	}

	mock.ExpectQuery(`INSERT INTO tasks`).
		WithArgs(1, "Test Task", "Test Description").
		WillReturnError(sql.ErrConnDone)

	err = repo.CreateTask(task)
	if err == nil {
		t.Error("CreateTask должен вернуть ошибку")
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ GetAllTasksByUser
// ============================================================================

func TestGetAllTasksByUserSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "text", "complete", "create_time", "complete_at"}).
		AddRow(1, 1, "Task 1", "Description 1", false, now, nil).
		AddRow(2, 1, "Task 2", "Description 2", true, now, &now)

	mock.ExpectQuery(`SELECT id, user_id, name, text, complete, create_time, complete_at FROM tasks`).
		WithArgs(1).
		WillReturnRows(rows)

	tasks, err := repo.GetAllTasksByUser(1)
	if err != nil {
		t.Errorf("GetAllTasksByUser вернул ошибку: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Ожидалось 2 задачи, получено %d", len(tasks))
	}

	if tasks[0].ID != 1 || tasks[1].ID != 2 {
		t.Error("Неправильные ID задач")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Не выполнены ожидания mock: %v", err)
	}
}

func TestGetAllTasksByUserEmpty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "text", "complete", "create_time", "complete_at"})

	mock.ExpectQuery(`SELECT id, user_id, name, text, complete, create_time, complete_at FROM tasks`).
		WithArgs(1).
		WillReturnRows(rows)

	tasks, err := repo.GetAllTasksByUser(1)
	if err != nil {
		t.Errorf("GetAllTasksByUser вернул ошибку: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("Ожидался пустой список, получено %d задач", len(tasks))
	}
}

func TestGetAllTasksByUserError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	mock.ExpectQuery(`SELECT id, user_id, name, text, complete, create_time, complete_at FROM tasks`).
		WithArgs(1).
		WillReturnError(sql.ErrConnDone)

	_, err = repo.GetAllTasksByUser(1)
	if err == nil {
		t.Error("GetAllTasksByUser должен вернуть ошибку")
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ GetCompletedTasksByUser
// ============================================================================

func TestGetCompletedTasksByUserSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "text", "complete", "create_time", "complete_at"}).
		AddRow(1, 1, "Completed Task", "Description", true, now, &now)

	mock.ExpectQuery(`SELECT id, user_id, name, text, complete, create_time, complete_at FROM tasks`).
		WithArgs(1).
		WillReturnRows(rows)

	tasks, err := repo.GetCompletedTasksByUser(1)
	if err != nil {
		t.Errorf("GetCompletedTasksByUser вернул ошибку: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Ожидалась 1 задача, получено %d", len(tasks))
	}

	if !tasks[0].Complete {
		t.Error("Задача должна быть завершена")
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ GetUncompletedTasksByUser
// ============================================================================

func TestGetUncompletedTasksByUserSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "text", "complete", "create_time", "complete_at"}).
		AddRow(1, 1, "Uncompleted Task", "Description", false, now, nil)

	mock.ExpectQuery(`SELECT id, user_id, name, text, complete, create_time, complete_at FROM tasks`).
		WithArgs(1).
		WillReturnRows(rows)

	tasks, err := repo.GetUncompletedTasksByUser(1)
	if err != nil {
		t.Errorf("GetUncompletedTasksByUser вернул ошибку: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Ожидалась 1 задача, получено %d", len(tasks))
	}

	if tasks[0].Complete {
		t.Error("Задача не должна быть завершена")
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ DeleteTaskByUser
// ============================================================================

func TestDeleteTaskByUserSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	mock.ExpectExec(`DELETE FROM tasks WHERE id`).
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteTaskByUser(1, 1)
	if err != nil {
		t.Errorf("DeleteTaskByUser вернул ошибку: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Не выполнены ожидания mock: %v", err)
	}
}

func TestDeleteTaskByUserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	mock.ExpectExec(`DELETE FROM tasks WHERE id`).
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.DeleteTaskByUser(1, 1)
	if err == nil {
		t.Error("DeleteTaskByUser должен вернуть ошибку когда задача не найдена")
	}
}

func TestDeleteTaskByUserError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	mock.ExpectExec(`DELETE FROM tasks WHERE id`).
		WithArgs(1, 1).
		WillReturnError(sql.ErrConnDone)

	err = repo.DeleteTaskByUser(1, 1)
	if err == nil {
		t.Error("DeleteTaskByUser должен вернуть ошибку")
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ CompleteTaskByUser
// ============================================================================

func TestCompleteTaskByUserSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	mock.ExpectExec(`UPDATE tasks`).
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.CompleteTaskByUser(1, 1)
	if err != nil {
		t.Errorf("CompleteTaskByUser вернул ошибку: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Не выполнены ожидания mock: %v", err)
	}
}

func TestCompleteTaskByUserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	mock.ExpectExec(`UPDATE tasks`).
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.CompleteTaskByUser(1, 1)
	if err == nil {
		t.Error("CompleteTaskByUser должен вернуть ошибку когда задача не найдена")
	}
}

func TestCompleteTaskByUserError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	mock.ExpectExec(`UPDATE tasks`).
		WithArgs(1, 1).
		WillReturnError(sql.ErrConnDone)

	err = repo.CompleteTaskByUser(1, 1)
	if err == nil {
		t.Error("CompleteTaskByUser должен вернуть ошибку")
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ GetTaskByID
// ============================================================================

func TestGetTaskByIDSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "text", "complete", "create_time", "complete_at"}).
		AddRow(1, "Task 1", "Description", false, now, nil)

	mock.ExpectQuery(`SELECT \* FROM tasks`).
		WithArgs(1).
		WillReturnRows(rows)

	task, err := repo.GetTaskByID(1)
	if err != nil {
		t.Errorf("GetTaskByID вернул ошибку: %v", err)
	}

	if task.ID != 1 {
		t.Errorf("Ожидался ID 1, получен %d", task.ID)
	}
}

func TestGetTaskByIDNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	mock.ExpectQuery(`SELECT \* FROM tasks`).
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetTaskByID(1)
	if err == nil {
		t.Error("GetTaskByID должен вернуть ошибку когда задача не найдена")
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ GetIDByName
// ============================================================================

func TestGetIDByNameSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "text", "complete", "create_time", "complete_at"}).
		AddRow(1, "Task 1", "Description", false, now, nil)

	mock.ExpectQuery(`SELECT \* FROM tasks`).
		WithArgs("Task 1").
		WillReturnRows(rows)

	task, err := repo.GetIDByName("Task 1")
	if err != nil {
		t.Errorf("GetIDByName вернул ошибку: %v", err)
	}

	if task.Name != "Task 1" {
		t.Errorf("Ожидалось имя 'Task 1', получено %s", task.Name)
	}
}

func TestGetIDByNameNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	mock.ExpectQuery(`SELECT \* FROM tasks`).
		WithArgs("NonExistent").
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetIDByName("NonExistent")
	if err == nil {
		t.Error("GetIDByName должен вернуть ошибку когда задача не найдена")
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ GetAllTasks (без user_id)
// ============================================================================

func TestGetAllTasksSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "text", "complete", "createtime", "completeat"}).
		AddRow(1, "Task 1", "Description 1", false, now, nil).
		AddRow(2, "Task 2", "Description 2", true, now, &now)

	mock.ExpectQuery(`SELECT \* FROM tasks`).
		WillReturnRows(rows)

	tasks, err := repo.GetAllTasks()
	if err != nil {
		t.Errorf("GetAllTasks вернул ошибку: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Ожидалось 2 задачи, получено %d", len(tasks))
	}
}

func TestGetAllTasksError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	mock.ExpectQuery(`SELECT \* FROM tasks`).
		WillReturnError(sql.ErrConnDone)

	_, err = repo.GetAllTasks()
	if err == nil {
		t.Error("GetAllTasks должен вернуть ошибку")
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ GetCompletedTasks (без user_id)
// ============================================================================

func TestGetCompletedTasksSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "text", "complete", "createtime", "completeat"}).
		AddRow(1, "Completed Task", "Description", true, now, &now)

	mock.ExpectQuery(`SELECT \* FROM tasks`).
		WillReturnRows(rows)

	tasks, err := repo.GetCompletedTasks()
	if err != nil {
		t.Errorf("GetCompletedTasks вернул ошибку: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Ожидалась 1 задача, получено %d", len(tasks))
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ GetUncompletedTasks (без user_id)
// ============================================================================

func TestGetUncompletedTasksSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "text", "complete", "createtime", "completeat"}).
		AddRow(1, "Uncompleted Task", "Description", false, now, nil)

	mock.ExpectQuery(`SELECT \* FROM tasks`).
		WillReturnRows(rows)

	tasks, err := repo.GetUncompletedTasks()
	if err != nil {
		t.Errorf("GetUncompletedTasks вернул ошибку: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Ожидалась 1 задача, получено %d", len(tasks))
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ DeleteTask (без user_id)
// ============================================================================

func TestDeleteTaskSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	mock.ExpectExec(`DELETE FROM tasks WHERE id`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteTask(1)
	if err != nil {
		t.Errorf("DeleteTask вернул ошибку: %v", err)
	}
}

func TestDeleteTaskError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	mock.ExpectExec(`DELETE FROM tasks WHERE id`).
		WithArgs(1).
		WillReturnError(sql.ErrConnDone)

	err = repo.DeleteTask(1)
	if err == nil {
		t.Error("DeleteTask должен вернуть ошибку")
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ CompleteTask (без user_id)
// ============================================================================

func TestCompleteTaskSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	mock.ExpectExec(`UPDATE tasks`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.CompleteTask(1)
	if err != nil {
		t.Errorf("CompleteTask вернул ошибку: %v", err)
	}
}

func TestCompleteTaskNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	mock.ExpectExec(`UPDATE tasks`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.CompleteTask(1)
	if err == nil {
		t.Error("CompleteTask должен вернуть ошибку когда задача не найдена")
	}
}

func TestCompleteTaskError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	mock.ExpectExec(`UPDATE tasks`).
		WithArgs(1).
		WillReturnError(sql.ErrConnDone)

	err = repo.CompleteTask(1)
	if err == nil {
		t.Error("CompleteTask должен вернуть ошибку")
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ CreateTable
// ============================================================================

func TestCreateTableSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS users`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS tasks`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.CreateTable()
	if err != nil {
		t.Errorf("CreateTable вернул ошибку: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Не выполнены ожидания mock: %v", err)
	}
}

func TestCreateTableUsersError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS users`).
		WillReturnError(sql.ErrConnDone)

	err = repo.CreateTable()
	if err == nil {
		t.Error("CreateTable должен вернуть ошибку при ошибке создания таблицы users")
	}
}

func TestCreateTableTasksError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS users`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS tasks`).
		WillReturnError(sql.ErrConnDone)

	err = repo.CreateTable()
	if err == nil {
		t.Error("CreateTable должен вернуть ошибку при ошибке создания таблицы tasks")
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ SCAN ERRORS
// ============================================================================

func TestGetAllTasksByUserScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	// Неправильное количество колонок
	rows := sqlmock.NewRows([]string{"id", "user_id", "name"}).
		AddRow(1, 1, "Task 1")

	mock.ExpectQuery(`SELECT id, user_id, name, text, complete, create_time, complete_at FROM tasks`).
		WithArgs(1).
		WillReturnRows(rows)

	_, err = repo.GetAllTasksByUser(1)
	if err == nil {
		t.Error("GetAllTasksByUser должен вернуть ошибку при ошибке Scan")
	}
}

func TestGetCompletedTasksByUserScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	rows := sqlmock.NewRows([]string{"id", "user_id", "name"}).
		AddRow(1, 1, "Task 1")

	mock.ExpectQuery(`SELECT id, user_id, name, text, complete, create_time, complete_at FROM tasks`).
		WithArgs(1).
		WillReturnRows(rows)

	_, err = repo.GetCompletedTasksByUser(1)
	if err == nil {
		t.Error("GetCompletedTasksByUser должен вернуть ошибку при ошибке Scan")
	}
}

func TestGetUncompletedTasksByUserScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	rows := sqlmock.NewRows([]string{"id", "user_id", "name"}).
		AddRow(1, 1, "Task 1")

	mock.ExpectQuery(`SELECT id, user_id, name, text, complete, create_time, complete_at FROM tasks`).
		WithArgs(1).
		WillReturnRows(rows)

	_, err = repo.GetUncompletedTasksByUser(1)
	if err == nil {
		t.Error("GetUncompletedTasksByUser должен вернуть ошибку при ошибке Scan")
	}
}

func TestGetAllTasksScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Task 1")

	mock.ExpectQuery(`SELECT \* FROM tasks`).
		WillReturnRows(rows)

	_, err = repo.GetAllTasks()
	if err == nil {
		t.Error("GetAllTasks должен вернуть ошибку при ошибке Scan")
	}
}

func TestGetCompletedTasksScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Task 1")

	mock.ExpectQuery(`SELECT \* FROM tasks`).
		WillReturnRows(rows)

	_, err = repo.GetCompletedTasks()
	if err == nil {
		t.Error("GetCompletedTasks должен вернуть ошибку при ошибке Scan")
	}
}

func TestGetUncompletedTasksScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock: %v", err)
	}
	defer db.Close()

	repo := NewTaskRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Task 1")

	mock.ExpectQuery(`SELECT \* FROM tasks`).
		WillReturnRows(rows)

	_, err = repo.GetUncompletedTasks()
	if err == nil {
		t.Error("GetUncompletedTasks должен вернуть ошибку при ошибке Scan")
	}
}
