package captainhook

import (
	"context"
	"encoding/json"
	pb "github.com/belljustin/captainhook/proto/captainhook"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"log"
	"time"
)

const (
	TypeSignMessage = "message:sign"

	TypeCreateSubscription = "subscription:create"
)

type signMessagePayload struct {
	TenantID      uuid.UUID
	ID            uuid.UUID
	ApplicationID uuid.UUID
	Type          string
	Data          []byte
}

func NewSignMessageTask(id, tenantID, appID uuid.UUID, msgType string, data []byte) (*asynq.Task, error) {
	payload, err := json.Marshal(signMessagePayload{
		TenantID:      tenantID,
		ID:            id,
		ApplicationID: appID,
		Type:          msgType,
		Data:          data,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeSignMessage, payload), nil
}

type SignMessageTaskHandler struct {
	Storage Storage
}

func (h *SignMessageTaskHandler) Handle(ctx context.Context, t *asynq.Task) error {
	now := time.Now()

	var p signMessagePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	log.Printf(" [*] Sign Message %q", p.Data)

	msg := Message{
		TenantID: p.TenantID,
		ID:       p.ID,

		ApplicationID: p.ApplicationID,
		Type:          p.Type,
		Data:          p.Data,
		State:         pb.Message_PENDING.String(),
		Signature:     []byte{},

		TimeDetails: TimeDetails{
			CreateTime: now,
			UpdateTime: now,
		},
	}
	_, err := h.Storage.NewMessage(ctx, &msg)
	return err
}

type createSubscriptionPayload struct {
	TenantID      uuid.UUID
	ApplicationID uuid.UUID
	ID            uuid.UUID
	Name          string
	Types         []string
}

func NewCreateSubscriptionTask(id, tenantID, appID uuid.UUID, name string, types []string) (*asynq.Task, error) {
	payload, err := json.Marshal(createSubscriptionPayload{
		TenantID:      tenantID,
		ApplicationID: appID,
		ID:            id,
		Name:          name,
		Types:         types,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeCreateSubscription, payload), nil
}

type CreateSubscriptionTaskHandler struct {
	Storage Storage
}

func (h *CreateSubscriptionTaskHandler) HandleCreateSubscription(ctx context.Context, t *asynq.Task) error {
	now := time.Now()

	var p createSubscriptionPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	log.Printf(" [*] Create subscription %s", p.Name)

	sub := Subscription{
		TenantID: p.TenantID,
		ID:       p.ID,

		ApplicationID: p.ApplicationID,
		Name:          p.Name,
		Types:         p.Types,
		State:         pb.Subscription_PENDING.String(),

		TimeDetails: TimeDetails{
			CreateTime: now,
			UpdateTime: now,
		},
	}
	_, err := h.Storage.NewSubscription(ctx, &sub)
	// TODO: queue subscription confirmation delivery
	return err
}
