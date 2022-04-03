package postgres

import (
	"context"
	"database/sql"
	"github.com/belljustin/captainhook/captainhook"
	"github.com/belljustin/captainhook/internal"
	pb "github.com/belljustin/captainhook/proto/captainhook"
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

func (s *PostgresTestSuite) TestSubscription() {
	tenantID, err := uuid.NewRandom()
	require.NoError(s.T(), err)
	appID, err := uuid.NewRandom()
	require.NoError(s.T(), err)

	n := 4
	subs := make([]captainhook.Subscription, n)
	for i := 0; i < n-1; i++ {
		subs[i], err = randomSub(tenantID, appID, []string{"ch/test", "ch/test2"})
		require.NoError(s.T(), err)

		_, err := s.storage.NewSubscription(context.Background(), &subs[i])
		require.NoError(s.T(), err)
	}
	subs[n-1], err = randomSub(tenantID, appID, nil)
	require.NoError(s.T(), err)
	_, err = s.storage.NewSubscription(context.Background(), &subs[n-1])
	require.NoError(s.T(), err)

	subCollection, err := s.storage.GetSubscriptions(context.Background(), tenantID, appID, &pgPageOpt{Size: 20})
	require.NoError(s.T(), err)

	for i := 0; i < n; i++ {
		assert.Equal(s.T(), subs[i], subCollection.Results[n-1-i])
	}

	subCollection, err = s.storage.GetSubscriptions(context.Background(), tenantID, appID, &pgPageOpt{Size: 1})
	require.NoError(s.T(), err)
	assert.Len(s.T(), subCollection.Results, 1)
	assert.Equal(s.T(), subs[n-1], subCollection.Results[0])

	nextPageToken := pgPageTokenString(subCollection.NextPageToken, 1)
	subCollection, err = s.storage.GetSubscriptions(context.Background(), tenantID, appID, &nextPageToken)
	require.NoError(s.T(), err)
	assert.Len(s.T(), subCollection.Results, 1)
	assert.Equal(s.T(), subs[n-2], subCollection.Results[0])
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

func randomSub(tenantID, appID uuid.UUID, types []string) (captainhook.Subscription, error) {
	subID, err := uuid.NewRandom()
	if err != nil {
		return captainhook.Subscription{}, err
	}
	now := time.Now().UTC()

	return captainhook.Subscription{
		TenantID: tenantID,
		ID:       subID,

		ApplicationID: appID,
		Name:          internal.RandString(16),
		Types:         types,
		State:         pb.Subscription_PENDING.String(),
		Endpoint:      internal.RandString(16),
		TimeDetails: captainhook.TimeDetails{
			CreateTime: now,
			UpdateTime: now,
		},
	}, nil
}
