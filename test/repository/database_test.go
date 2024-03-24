package repository_test

import (
	"testing"

	"github.com/MaximInnopolis/ProductCatalog/internal/repository"
)

func TestNewDB(t *testing.T) {
	db, err := repository.NewDB(":memory:")
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatalf("failed to ping database: %s", err)
	}
}

func TestClose(t *testing.T) {
	db, err := repository.NewDB(":memory:")
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatalf("failed to ping database: %s", err)
	}
}
