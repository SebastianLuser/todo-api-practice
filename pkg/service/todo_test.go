package service_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"todo-api/pkg/domain"
	"todo-api/pkg/service"
)

const (
	validUUID        = "123e4567-e89b-12d3-a456-426614174000"
	nonExistentID    = "00000000-0000-0000-0000-000000000000"
	validTitle       = "Test Todo"
	validDescription = "Test Description"
)

var fixedTime = time.Date(2026, 1, 28, 10, 30, 0, 0, time.UTC)

func TestService_Get_ReturnsListSuccessfully(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "title", "description", "status", "priority", "created_at", "updated_at"}).
		AddRow(validUUID, validTitle, validDescription, domain.StatusPending, domain.PriorityMedium, fixedTime, fixedTime)
	mock.ExpectQuery("SELECT").WithArgs(nil, nil).WillReturnRows(rows)
	svc := service.New(db)

	result, err := svc.Get(context.Background(), service.Filters{})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 todo, got %d", len(result))
	}
	if result[0].ID != validUUID {
		t.Errorf("expected ID %s, got %s", validUUID, result[0].ID)
	}
}

func TestService_Get_WithStatusFilter(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "title", "description", "status", "priority", "created_at", "updated_at"}).
		AddRow(validUUID, validTitle, validDescription, domain.StatusCompleted, domain.PriorityMedium, fixedTime, fixedTime)
	status := domain.StatusCompleted
	mock.ExpectQuery("SELECT").WithArgs(string(status), nil).WillReturnRows(rows)
	svc := service.New(db)

	result, err := svc.Get(context.Background(), service.Filters{Status: &status})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 todo, got %d", len(result))
	}
}

func TestService_Get_WithPriorityFilter(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "title", "description", "status", "priority", "created_at", "updated_at"}).
		AddRow(validUUID, validTitle, validDescription, domain.StatusPending, domain.PriorityHigh, fixedTime, fixedTime)
	priority := domain.PriorityHigh
	mock.ExpectQuery("SELECT").WithArgs(nil, string(priority)).WillReturnRows(rows)
	svc := service.New(db)

	result, err := svc.Get(context.Background(), service.Filters{Priority: &priority})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 todo, got %d", len(result))
	}
}

func TestService_Get_ReturnsEmptyList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "title", "description", "status", "priority", "created_at", "updated_at"})
	mock.ExpectQuery("SELECT").WithArgs(nil, nil).WillReturnRows(rows)
	svc := service.New(db)

	result, err := svc.Get(context.Background(), service.Filters{})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 todos, got %d", len(result))
	}
}

func TestService_Get_ReturnsErrorOnQueryFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()
	expectedErr := errors.New("database error")
	mock.ExpectQuery("SELECT").WithArgs(nil, nil).WillReturnError(expectedErr)
	svc := service.New(db)

	_, err = svc.Get(context.Background(), service.Filters{})

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestService_GetByID_ReturnsTodoSuccessfully(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "title", "description", "status", "priority", "created_at", "updated_at"}).
		AddRow(validUUID, validTitle, validDescription, domain.StatusPending, domain.PriorityMedium, fixedTime, fixedTime)
	mock.ExpectQuery("SELECT").WithArgs(validUUID).WillReturnRows(rows)
	svc := service.New(db)

	result, err := svc.GetByID(context.Background(), validUUID)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.ID != validUUID {
		t.Errorf("expected ID %s, got %s", validUUID, result.ID)
	}
	if result.Title != validTitle {
		t.Errorf("expected title %s, got %s", validTitle, result.Title)
	}
}

func TestService_GetByID_ReturnsErrTodoNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()
	mock.ExpectQuery("SELECT").WithArgs(nonExistentID).WillReturnError(sql.ErrNoRows)
	svc := service.New(db)

	_, err = svc.GetByID(context.Background(), nonExistentID)

	if !errors.Is(err, domain.ErrTodoNotFound) {
		t.Errorf("expected ErrTodoNotFound, got %v", err)
	}
}

func TestService_GetByID_ReturnsErrorOnQueryFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()
	expectedErr := errors.New("database error")
	mock.ExpectQuery("SELECT").WithArgs(validUUID).WillReturnError(expectedErr)

	svc := service.New(db)

	_, err = svc.GetByID(context.Background(), validUUID)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestService_Create_ReturnsTodoSuccessfully(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "title", "description", "status", "priority", "created_at", "updated_at"}).
		AddRow(validUUID, validTitle, validDescription, domain.StatusPending, domain.PriorityMedium, fixedTime, fixedTime)
	desc := validDescription
	mock.ExpectQuery("INSERT").
		WithArgs(validTitle, sqlmock.AnyArg(), domain.StatusPending, domain.PriorityMedium).
		WillReturnRows(rows)
	svc := service.New(db)
	input := service.CreateInput{
		Title:       validTitle,
		Description: &desc,
		Status:      domain.StatusPending,
		Priority:    domain.PriorityMedium,
	}

	result, err := svc.Create(context.Background(), input)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.ID != validUUID {
		t.Errorf("expected ID %s, got %s", validUUID, result.ID)
	}
}

func TestService_Create_WithoutDescription(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "title", "description", "status", "priority", "created_at", "updated_at"}).
		AddRow(validUUID, validTitle, nil, domain.StatusPending, domain.PriorityMedium, fixedTime, fixedTime)
	mock.ExpectQuery("INSERT").
		WithArgs(validTitle, sqlmock.AnyArg(), domain.StatusPending, domain.PriorityMedium).
		WillReturnRows(rows)
	svc := service.New(db)
	input := service.CreateInput{
		Title:    validTitle,
		Status:   domain.StatusPending,
		Priority: domain.PriorityMedium,
	}

	result, err := svc.Create(context.Background(), input)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Description != "" {
		t.Errorf("expected empty description, got %s", result.Description)
	}
}

func TestService_Create_ReturnsErrorOnQueryFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()
	expectedErr := errors.New("database error")
	mock.ExpectQuery("INSERT").
		WithArgs(validTitle, sqlmock.AnyArg(), domain.StatusPending, domain.PriorityMedium).
		WillReturnError(expectedErr)
	svc := service.New(db)
	input := service.CreateInput{
		Title:    validTitle,
		Status:   domain.StatusPending,
		Priority: domain.PriorityMedium,
	}

	_, err = svc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestService_Update_ReturnsTodoSuccessfully(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()
	updatedTitle := "Updated Title"
	rows := sqlmock.NewRows([]string{"id", "title", "description", "status", "priority", "created_at", "updated_at"}).
		AddRow(validUUID, updatedTitle, validDescription, domain.StatusPending, domain.PriorityMedium, fixedTime, fixedTime)

	mock.ExpectQuery("UPDATE").
		WithArgs(validUUID, updatedTitle, nil, nil, nil).
		WillReturnRows(rows)
	svc := service.New(db)
	input := service.UpdateInput{Title: &updatedTitle}

	result, err := svc.Update(context.Background(), validUUID, input)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Title != updatedTitle {
		t.Errorf("expected title %s, got %s", updatedTitle, result.Title)
	}
}

func TestService_Update_ReturnsErrTodoNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()
	updatedTitle := "Updated Title"
	mock.ExpectQuery("UPDATE").
		WithArgs(nonExistentID, updatedTitle, nil, nil, nil).
		WillReturnError(sql.ErrNoRows)
	svc := service.New(db)
	input := service.UpdateInput{Title: &updatedTitle}

	_, err = svc.Update(context.Background(), nonExistentID, input)

	if !errors.Is(err, domain.ErrTodoNotFound) {
		t.Errorf("expected ErrTodoNotFound, got %v", err)
	}
}

func TestService_Update_ReturnsErrorOnQueryFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()
	expectedErr := errors.New("database error")
	updatedTitle := "Updated Title"
	mock.ExpectQuery("UPDATE").
		WithArgs(validUUID, updatedTitle, nil, nil, nil).
		WillReturnError(expectedErr)
	svc := service.New(db)
	input := service.UpdateInput{Title: &updatedTitle}

	_, err = svc.Update(context.Background(), validUUID, input)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestService_Delete_Successfully(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()
	mock.ExpectExec("DELETE").WithArgs(validUUID).WillReturnResult(sqlmock.NewResult(0, 1))
	svc := service.New(db)

	err = svc.Delete(context.Background(), validUUID)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestService_Delete_ReturnsErrTodoNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()
	mock.ExpectExec("DELETE").WithArgs(nonExistentID).WillReturnResult(sqlmock.NewResult(0, 0))
	svc := service.New(db)

	err = svc.Delete(context.Background(), nonExistentID)

	if !errors.Is(err, domain.ErrTodoNotFound) {
		t.Errorf("expected ErrTodoNotFound, got %v", err)
	}
}

func TestService_Delete_ReturnsErrorOnExecFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()
	expectedErr := errors.New("database error")
	mock.ExpectExec("DELETE").WithArgs(validUUID).WillReturnError(expectedErr)
	svc := service.New(db)

	err = svc.Delete(context.Background(), validUUID)

	if err == nil {
		t.Error("expected error, got nil")
	}
}
