package captainhook

import (
	"context"
	"github.com/hibiken/asynq"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/belljustin/captainhook/proto/captainhook"
)

type Storage interface {
	NewApplication(ctx context.Context, app *Application) (*Application, error)
	GetApplication(ctx context.Context, tenantID, id uuid.UUID) (*Application, error)

	NewMessage(ctx context.Context, msg *Message) (*Message, error)

	NewSubscription(ctx context.Context, sub *Subscription) (*Subscription, error)
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

type Subscription struct {
	TenantID uuid.UUID `db:"tenant_id"`
	ID       uuid.UUID

	ApplicationID uuid.UUID `db:"application_id"`
	Name          string
	Types         []string
	State         string

	TimeDetails
}

func (s *Subscription) ToProtobuf() *pb.Subscription {
	return nil
}

func CreateSubscription(asynqClient *asynq.Client, tenantID, appID uuid.UUID, name string, types []string) (uuid.UUID, error) {
	id, _ := uuid.NewRandom()
	t1, err := NewCreateSubscriptionTask(id, tenantID, appID, name, types)
	if err != nil {
		return uuid.UUID{}, err
	}

	if _, err := asynqClient.Enqueue(t1); err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}
