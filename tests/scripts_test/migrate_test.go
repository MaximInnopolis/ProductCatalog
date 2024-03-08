package scripts_test

import (
	"database/sql"
	"errors"
	"github.com/MaximInnopolis/ProductCatalog/scripts"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrate(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	// Run migration
	err = scripts.Migrate(db)
	require.NoError(t, err)

	// Check if tables were created
	tables := []string{"categories", "products", "product_categories", "users"}
	for _, table := range tables {
		exists, err := tableExists(db, table)
		require.NoError(t, err)
		assert.True(t, exists, "Table %s should exist", table)
	}
}

func tableExists(db *sql.DB, tableName string) (bool, error) {
	query := "SELECT name FROM sqlite_master WHERE type='table' AND name=?"
	var name string
	err := db.QueryRow(query, tableName).Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
