package postgres

import (
	"context"
	"database/sql"
	"github.com/belljustin/captainhook/internal"
	"github.com/stretchr/testify/assert"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/belljustin/captainhook"
)

const databaseTestName = "postgres"

type PostgresTestSuite struct {
	suite.Suite
	db      *sql.DB
	m       *migrate.Migrate
	storage captainhook.Storage
}

func (s *PostgresTestSuite) SetupSuite() {
	var err error

	s.db, err = sql.Open("pgx", *database)
	require.NoError(s.T(), err)

	_, filename, _, _ := runtime.Caller(0)
	migrationPath := "file://" + path.Join(path.Dir(filename), "migrations")
	driver, err := pgx.WithInstance(s.db, &pgx.Config{DatabaseName: databaseTestName})
	require.NoError(s.T(), err)
	s.m, err = migrate.NewWithDatabaseInstance(migrationPath, databaseTestName, driver)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), s.m)
	require.NoError(s.T(), s.m.Up())

	s.storage = NewStorage()
}

func (s *PostgresTestSuite) TestApplication() {
	tenantID, err := uuid.NewRandom()
	require.NoError(s.T(), err)
	appID, err := uuid.NewRandom()
	require.NoError(s.T(), err)
	now := time.Now().UTC()

	app := &captainhook.Application{
		TenantID: tenantID,
		ID:       appID,
		Name:     internal.RandString(16),
		TimeDetails: captainhook.TimeDetails{
			CreateTime: now,
			UpdateTime: now,
		},
	}

	_, err = s.storage.NewApplication(context.Background(), app)
	require.NoError(s.T(), err)

	recvApp, err := s.storage.GetApplication(context.Background(), tenantID, appID)
	assert.Equal(s.T(), app, recvApp)
}

func (s *PostgresTestSuite) TearDownSuite() {
	require.NoError(s.T(), s.m.Down())
}

func TestPostgresTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long-running test.")
	}

	suite.Run(t, new(PostgresTestSuite))
}
