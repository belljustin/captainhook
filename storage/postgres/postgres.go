package postgres

import (
	"context"
	"flag"
	"fmt"
	"github.com/belljustin/captainhook/captainhook"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
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
		"VALUES (:id, :tenant_id, :application_id, :type, :data, :state, :signature, :create_time, :update_time) "+
		"ON CONFLICT DO NOTHING", msg)
	return msg, err
}

func (s *Storage) NewSubscription(ctx context.Context, sub *captainhook.Subscription) (*captainhook.Subscription, error) {
	var retSub captainhook.Subscription
	stmt, err := s.db.PrepareNamed("INSERT INTO subscriptions (id, tenant_id, application_id, name, endpoint, types, state, create_time, update_time) " +
		"VALUES (:id, :tenant_id, :application_id, :name, :endpoint, :types, :state, :create_time, :update_time) " +
		"ON CONFLICT(id) DO UPDATE SET id=EXCLUDED.id " +
		"RETURNING *")
	if err != nil {
		return nil, err
	}
	err = stmt.QueryRowx(sub).StructScan(&retSub)
	return &retSub, err
}

const pgPageTokenTimeFormat = time.RFC3339Nano

type pgPageOpt struct {
	ID            uuid.UUID
	CreatedAfter  time.Time
	CreatedBefore time.Time
	Size          int32
}

const (
	pgPageTokenIdKey            = "id"
	pgPageTokenCreatedAfterKey  = "gte"
	pgPageTokenCreatedBeforeKey = "lte"
)

func (t *pgPageOpt) GetPageToken() string {
	if t == nil || t.ID == uuid.Nil {
		return ""
	}

	if !t.CreatedBefore.IsZero() {
		createdBefore := t.CreatedBefore.Format(pgPageTokenTimeFormat)
		return fmt.Sprintf("%s=%s;%s=%s", pgPageTokenIdKey, t.ID, pgPageTokenCreatedBeforeKey, createdBefore)
	} else if !t.CreatedAfter.IsZero() {
		createdAfter := t.CreatedAfter.Format(pgPageTokenTimeFormat)
		return fmt.Sprintf("%s=%s;%s=%s", pgPageTokenIdKey, t.ID, pgPageTokenCreatedAfterKey, createdAfter)
	}

	return ""
}

func (t *pgPageOpt) GetPageSize() int32 {
	return t.Size
}

func pgPageTokenString(pageToken string, pageSize int32) pgPageOpt {
	ret := pgPageOpt{Size: pageSize}
	tokens := strings.Split(pageToken, ";")
	for _, token := range tokens {
		splitToken := strings.Split(token, "=")
		if len(splitToken) != 2 {
			continue
		}
		key, value := splitToken[0], splitToken[1]
		switch key {
		case pgPageTokenIdKey:
			id, err := uuid.Parse(value)
			if err == nil {
				ret.ID = id
			}
		case pgPageTokenCreatedBeforeKey:
			createdBefore, err := time.Parse(pgPageTokenTimeFormat, value)
			if err == nil {
				ret.CreatedBefore = createdBefore
			} else {
				log.Printf("[DEBUG] could not parse createdAfter: %s", err)
			}
		case pgPageTokenCreatedAfterKey:
			createdAfter, err := time.Parse(pgPageTokenTimeFormat, value)
			if err == nil {
				ret.CreatedAfter = createdAfter
			} else {
				log.Printf("[DEBUG] could not parse createdBefore: %s", err)
			}
		}
	}

	return ret
}

func newPrevPageToken(subscriptions []captainhook.Subscription) string {
	if len(subscriptions) == 0 {
		return ""
	}
	firstSub := subscriptions[0]
	t := &pgPageOpt{
		ID:           firstSub.ID,
		CreatedAfter: firstSub.CreateTime,
	}

	return t.GetPageToken()
}

func newNextPageToken(subscriptions []captainhook.Subscription) string {
	if len(subscriptions) == 0 {
		return ""
	}
	lastSub := subscriptions[len(subscriptions)-1]
	t := &pgPageOpt{
		ID:            lastSub.ID,
		CreatedBefore: lastSub.CreateTime,
	}
	return t.GetPageToken()
}

func (s *Storage) GetSubscriptions(ctx context.Context, tenantID, applicationID uuid.UUID, pageOpt captainhook.PaginationOpt) (*captainhook.SubscriptionCollection, error) {
	var pageToken pgPageOpt
	if pageOpt != nil {
		pageToken = pgPageTokenString(pageOpt.GetPageToken(), pageOpt.GetPageSize())
	} else {
		pageToken = pgPageOpt{Size: 20}
	}

	var err error
	var subscriptions []captainhook.Subscription
	if !pageToken.CreatedAfter.IsZero() {
		subscriptions, err = s.getPrevSubscriptions(ctx, tenantID, applicationID, pageToken, pageToken.GetPageSize())
	} else {
		subscriptions, err = s.getNextSubscriptions(ctx, tenantID, applicationID, pageToken, pageToken.GetPageSize())
	}

	if err != nil {
		return nil, err
	}

	return &captainhook.SubscriptionCollection{
		Results:       subscriptions,
		NextPageToken: newNextPageToken(subscriptions),
		PrevPageToken: newPrevPageToken(subscriptions),
	}, nil
}

func (s *Storage) getNextSubscriptions(ctx context.Context, tenantID, applicationID uuid.UUID, pageToken pgPageOpt, pageSize int32) ([]captainhook.Subscription, error) {
	createdBefore := pageToken.CreatedBefore
	if createdBefore.IsZero() {
		createdBefore = time.Now()
	}

	var subscriptions []captainhook.Subscription
	err := s.db.Select(&subscriptions, "SELECT * FROM subscriptions WHERE tenant_id = $1 AND application_id = $2 "+
		"AND (create_time, id) < ($3, $4) "+
		"ORDER BY create_time DESC, id DESC "+
		"LIMIT $5",
		tenantID,
		applicationID,
		createdBefore.UTC(),
		pageToken.ID,
		pageSize,
	)
	return subscriptions, err
}

func (s *Storage) getPrevSubscriptions(ctx context.Context, tenantID, applicationID uuid.UUID, pageToken pgPageOpt, pageSize int32) ([]captainhook.Subscription, error) {
	var subscriptions []captainhook.Subscription
	err := s.db.Select(&subscriptions, "SELECT * FROM subscriptions "+
		"WHERE tenant_id = $1 AND application_id = $2 "+
		"AND (create_time, id) > ($3, $4) "+
		"ORDER BY create_time DESC, id DESC "+
		"LIMIT $5",
		tenantID,
		applicationID,
		pageToken.CreatedAfter.UTC(),
		pageToken.ID,
		pageSize,
	)
	return subscriptions, err
}
