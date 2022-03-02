package captainhook

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/belljustin/captainhook/proto/captainhook"
)

type Storage interface {
	NewApplication(ctx context.Context, app *Application) (*Application, error)
	GetApplication(ctx context.Context, tenantID, id uuid.UUID) (*Application, error)
}

type Application struct {
	TenantID   uuid.UUID `db:"tenant_id"`
	ID         uuid.UUID
	Name       string
	CreateTime time.Time `db:"create_time"`
	UpdateTime time.Time `db:"update_time"`
}

func NewApplication(storage Storage, tenantID uuid.UUID, name string) (*Application, error) {
	id, _ := uuid.NewRandom()
	now := time.Now()

	app := &Application{
		TenantID:   tenantID,
		ID:         id,
		Name:       name,
		CreateTime: now,
		UpdateTime: now,
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
