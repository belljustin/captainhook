package captainhook

import (
	"context"
	"database/sql/driver"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/belljustin/captainhook/proto/captainhook"
)

type PaginationOpt interface {
	GetPageToken() string
	GetPageSize() int32
}

type Storage interface {
	NewApplication(ctx context.Context, app *Application) (*Application, error)
	GetApplication(ctx context.Context, tenantID, id uuid.UUID) (*Application, error)

	NewMessage(ctx context.Context, msg *Message) (*Message, error)

	NewSubscription(ctx context.Context, sub *Subscription) (*Subscription, error)
	GetSubscriptions(ctx context.Context, tenantID, appID uuid.UUID, pageOpt PaginationOpt) (*SubscriptionCollection, error)
}

type TimeDetails struct {
	CreateTime time.Time `db:"create_time"`
	UpdateTime time.Time `db:"update_time"`
}

type Application struct {
	TenantID uuid.UUID `db:"tenant_id"`
	ID       uuid.UUID

	Name string

	TimeDetails
}

func NewApplication(storage Storage, tenantID uuid.UUID, name string) (*Application, error) {
	id, _ := uuid.NewRandom()
	now := time.Now()

	app := &Application{
		TenantID: tenantID,
		ID:       id,
		Name:     name,
		TimeDetails: TimeDetails{
			CreateTime: now,
			UpdateTime: now,
		},
	}
	return storage.NewApplication(context.Background(), app)
}

func (app Application) ToProtobuf() *pb.Application {
	return &pb.Application{
		TenantId:   app.TenantID.String(),
		Id:         app.ID.String(),
		Name:       app.Name,
		CreateTime: timestamppb.New(app.CreateTime),
		UpdateTime: timestamppb.New(app.UpdateTime),
	}
}

func GetApplication(storage Storage, tenantID, id uuid.UUID) (*Application, error) {
	return storage.GetApplication(context.Background(), tenantID, id)
}

type Message struct {
	TenantID uuid.UUID `db:"tenant_id"`
	ID       uuid.UUID

	ApplicationID uuid.UUID `db:"application_id"`
	Type          string
	Data          []byte
	State         string
	Signature     []byte

	TimeDetails
}

func CreateMessage(asynqClient *asynq.Client, tenantID, appID uuid.UUID, msgType string, data []byte) (uuid.UUID, error) {
	id, _ := uuid.NewRandom()
	t1, err := NewSignMessageTask(id, tenantID, appID, msgType, data)
	if err != nil {
		return uuid.UUID{}, err
	}

	if _, err := asynqClient.Enqueue(t1); err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

type SubscriptionTypes []string

func (p *SubscriptionTypes) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	stypes := fmt.Sprintf("%v", src)
	stypes = stypes[1 : len(stypes)-1]
	if len(stypes) == 0 {
		*p = make([]string, 0)
		return nil
	}

	value := SubscriptionTypes(strings.Split(stypes, ","))
	*p = value
	return nil
}

func (p SubscriptionTypes) Value() (driver.Value, error) {
	if p == nil || len(p) == 0 {
		return nil, nil
	}
	value := strings.Join(p, ",")
	return "{" + value + "}", nil
}

type Subscription struct {
	TenantID uuid.UUID `db:"tenant_id"`
	ID       uuid.UUID

	ApplicationID uuid.UUID `db:"application_id"`
	Name          string
	Types         SubscriptionTypes
	State         string
	Endpoint      string

	TimeDetails
}

func (s *Subscription) ToProtobuf() *pb.Subscription {
	return &pb.Subscription{
		TenantId: s.TenantID.String(),
		Id:       s.ID.String(),

		ApplicationId: s.ApplicationID.String(),
		Name:          s.Name,
		Types:         s.Types,
		State:         pb.Subscription_State(pb.Subscription_State_value[s.State]),
		Endpoint:      s.Endpoint,

		CreateTime: timestamppb.New(s.CreateTime),
		UpdateTime: timestamppb.New(s.UpdateTime),
	}
}

func CreateSubscription(asynqClient *asynq.Client, tenantID, appID uuid.UUID, name string, types []string, endpoint *url.URL) (uuid.UUID, error) {
	id, _ := uuid.NewRandom()
	t1, err := NewCreateSubscriptionTask(id, tenantID, appID, name, types, endpoint)
	if err != nil {
		return uuid.UUID{}, err
	}

	if _, err := asynqClient.Enqueue(t1); err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

type SubscriptionCollection struct {
	Results       []Subscription
	NextPageToken string
	PrevPageToken string
}

func (col *SubscriptionCollection) ToProtobuf() *pb.SubscriptionCollection {
	results := make([]*pb.Subscription, len(col.Results))
	for i, sub := range col.Results {
		results[i] = sub.ToProtobuf()
	}

	return &pb.SubscriptionCollection{
		Results: results,
		Next:    col.NextPageToken,
		Prev:    col.PrevPageToken,
	}
}

func GetSubscriptions(storage Storage, tenantID, appID uuid.UUID, opt PaginationOpt) (*SubscriptionCollection, error) {
	return storage.GetSubscriptions(context.Background(), tenantID, appID, opt)
}
