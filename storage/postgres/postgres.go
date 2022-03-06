package postgres

import (
	"context"
	"flag"
	"log"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/belljustin/captainhook"
)

var (
	database = flag.String("database", "postgres://postgres:password@localhost:5432/postgres?sslmode=disable", "The database url")
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage() *Storage {
	db, err := sqlx.Connect("pgx", *database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	return &Storage{db}
}

func (s *Storage) NewApplication(ctx context.Context, app *captainhook.Application) (*captainhook.Application, error) {
	_, err := s.db.NamedExec("INSERT INTO applications (id, tenant_id, name, create_time, update_time) "+
		"VALUES (:id, :tenant_id, :name, :create_time, :update_time)", app)
	return app, err
}

func (s *Storage) GetApplication(ctx context.Context, tenantID, id uuid.UUID) (*captainhook.Application, error) {
	var app captainhook.Application
	err := s.db.Get(&app, "SELECT * FROM applications WHERE tenant_id = $1 AND id = $2", tenantID, id)
	return &app, err
}

func (s *Storage) NewMessage(ctx context.Context, msg *captainhook.Message) (*captainhook.Message, error) {
	_, err := s.db.NamedExec("INSERT INTO messages (id, tenant_id, application_id, type, data, state, signature, create_time, update_time) "+
		"VALUES (:id, :tenant_id, :application_id, :type, :data, :state, :signature, :create_time, :update_time)", msg)
	return msg, err
}

func (s *Storage) NewSubscription(ctx context.Context, sub *captainhook.Subscription) (*captainhook.Subscription, error) {
	return nil, nil
}
