package repository

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MaximInnopolis/ProductCatalog/internal/model"
	"github.com/MaximInnopolis/ProductCatalog/internal/repository"
	"testing"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	authDB := repository.NewAuthDatabase(db)

	ctx := context.Background()

	// Mocking successful insert
	mock.ExpectQuery("SELECT COUNT(.+)").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec("INSERT INTO users").WithArgs("username", "password").WillReturnResult(sqlmock.NewResult(1, 1))

	user := &model.User{
		Username: "username",
		Password: "password",
	}
	err = authDB.CreateUser(ctx, user)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Mocking username already exists
	mock.ExpectQuery("SELECT COUNT(.+)").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	err = authDB.CreateUser(ctx, user)
	if err == nil || err.Error() != "username already exists" {
		t.Errorf("expected 'username already exists' error, got: %v", err)
	}

	// Mocking error during insert
	mock.ExpectQuery("SELECT COUNT(.+)").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec("INSERT INTO users").WithArgs("username", "password").WillReturnError(errors.New("some error"))

	err = authDB.CreateUser(ctx, user)
	if err == nil || err.Error() != "some error" {
		t.Errorf("expected 'some error' error, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
