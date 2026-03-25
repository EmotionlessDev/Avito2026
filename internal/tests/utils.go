package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	Container  testcontainers.Container
	DB         *sql.DB
	ConnString string
	Ctx        context.Context
}

func SetupPostgres(t *testing.T) *PostgresContainer {
	t.Helper()

	ctx := context.Background()

	container, err := testpostgres.Run(ctx,
		"postgres:15-alpine",
		testpostgres.WithDatabase("booking_test"),
		testpostgres.WithUsername("postgres"),
		testpostgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	require.NoError(t, err)

	connString, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := sql.Open("postgres", connString)
	require.NoError(t, err)

	require.NoError(t, db.Ping())

	return &PostgresContainer{
		Container:  container,
		DB:         db,
		ConnString: connString,
		Ctx:        ctx,
	}
}

func TeardownPostgres(t *testing.T, pg *PostgresContainer) {
	t.Helper()

	if pg.DB != nil {
		pg.DB.Close()
	}

	if pg.Container != nil {
		err := pg.Container.Terminate(pg.Ctx)
		require.NoError(t, err)
	}
}

func RunMigrations(t *testing.T, db *sql.DB, migrationsPath string) {
	t.Helper()

	driver, err := migratepostgres.WithInstance(db, &migratepostgres.Config{})
	require.NoError(t, err)

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	require.NoError(t, err)

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err)
	}

	t.Logf("Migrations applied successfully from %s", migrationsPath)
}

func TruncateTables(t *testing.T, db *sql.DB) {
	t.Helper()

	_, err := db.Exec(`
        TRUNCATE TABLE bookings, slots, schedule_days, schedules, rooms, users 
        RESTART IDENTITY CASCADE
    `)
	require.NoError(t, err)
}

func NewTestTx(t *testing.T, db *sql.DB) *sql.Tx {
	t.Helper()

	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = tx.Rollback()
	})

	return tx
}

