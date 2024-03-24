package repository

import (
	"context"
	"errors"
	"github.com/MaximInnopolis/ProductCatalog/internal/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MaximInnopolis/ProductCatalog/internal/model"
)

func TestCreateCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	categoryDB := repository.NewCategoryDatabase(db)

	ctx := context.Background()

	// Mocking successful insert
	mock.ExpectExec("INSERT INTO categories").WithArgs("TestCategory").WillReturnResult(sqlmock.NewResult(1, 1))

	category := &model.Category{Name: "TestCategory"}
	_, err = categoryDB.CreateCategory(ctx, category)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Mocking error during insert
	mock.ExpectExec("INSERT INTO categories").WithArgs("TestCategory").WillReturnError(errors.New("some error"))

	_, err = categoryDB.CreateCategory(ctx, category)
	if err == nil || err.Error() != "some error" {
		t.Errorf("expected 'some error' error, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllCategories(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	categoryDB := repository.NewCategoryDatabase(db)

	ctx := context.Background()

	// Mocking successful query
	rows := sqlmock.NewRows([]string{"name"}).AddRow("Category1").AddRow("Category2")
	mock.ExpectQuery("SELECT name FROM categories").WillReturnRows(rows)

	categories, err := categoryDB.GetAllCategories(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(categories) != 2 || categories[0] != "Category1" || categories[1] != "Category2" {
		t.Errorf("expected categories not returned: %v", categories)
	}

	// Mocking error during query
	mock.ExpectQuery("SELECT name FROM categories").WillReturnError(errors.New("some error"))

	_, err = categoryDB.GetAllCategories(ctx)
	if err == nil || err.Error() != "some error" {
		t.Errorf("expected 'some error' error, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	categoryDB := repository.NewCategoryDatabase(db)

	ctx := context.Background()

	// Mocking successful update
	mock.ExpectExec("UPDATE categories").WithArgs("NewCategoryName", "OldCategoryName").WillReturnResult(sqlmock.NewResult(1, 1))

	err = categoryDB.UpdateCategory(ctx, "OldCategoryName", &model.Category{Name: "NewCategoryName"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Mocking error during update
	mock.ExpectExec("UPDATE categories").WithArgs("NewCategoryName", "OldCategoryName").WillReturnError(errors.New("some error"))

	err = categoryDB.UpdateCategory(ctx, "OldCategoryName", &model.Category{Name: "NewCategoryName"})
	if err == nil || err.Error() != "some error" {
		t.Errorf("expected 'some error' error, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	categoryDB := repository.NewCategoryDatabase(db)

	ctx := context.Background()

	// Mocking successful deletion
	mock.ExpectExec("DELETE FROM categories").WithArgs("CategoryName").WillReturnResult(sqlmock.NewResult(1, 1))

	err = categoryDB.DeleteCategory(ctx, "CategoryName")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Mocking error during deletion
	mock.ExpectExec("DELETE FROM categories").WithArgs("CategoryName").WillReturnError(errors.New("some error"))

	err = categoryDB.DeleteCategory(ctx, "CategoryName")
	if err == nil || err.Error() != "some error" {
		t.Errorf("expected 'some error' error, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
