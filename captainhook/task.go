package captainhook

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"

	pb "github.com/belljustin/captainhook/proto/captainhook"
)

const (
	TypeFanoutMessage   = "message:fanout"
	TypeDeliveryMessage = "message:delivery"

	TypeCreateSubscription = "subscription:create"
)

type createSubscriptionPayload struct {
	TenantID      uuid.UUID
	ApplicationID uuid.UUID
	ID            uuid.UUID
	Name          string
	Types         []string
	Endpoint      *url.URL
}

func NewCreateSubscriptionTask(id, tenantID, appID uuid.UUID, name string, types []string, endpoint *url.URL) (*asynq.Task, error) {
	payload, err := json.Marshal(createSubscriptionPayload{
		TenantID:      tenantID,
		ApplicationID: appID,
		ID:            id,
		Name:          name,
		Types:         types,
		Endpoint:      endpoint,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeCreateSubscription, payload), nil
}

type CreateSubscriptionTaskHandler struct {
	Storage Storage
}

func (h *CreateSubscriptionTaskHandler) Handle(ctx context.Context, t *asynq.Task) error {
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
		Types:         SubscriptionTypes(p.Types),
		State:         pb.Subscription_PENDING.String(),
		Endpoint:      p.Endpoint.String(),

		TimeDetails: TimeDetails{
			CreateTime: now,
			UpdateTime: now,
		},
	}
	_, err := h.Storage.NewSubscription(ctx, &sub)
	if err != nil {
		log.Printf(" [ERROR] Failed to insert subscription: %v", err)
	}
	// TODO: queue subscription confirmation delivery
	return err
}

type fanoutPayload struct {
	Message Message
}

func NewFanoutTask(message Message) (*asynq.Task, error) {
	payload, err := json.Marshal(fanoutPayload{Message: message})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeFanoutMessage, payload), nil
}

type FanoutTaskHandler struct {
	Storage     Storage
	AsynqClient *asynq.Client
}

func (h *FanoutTaskHandler) Handle(ctx context.Context, t *asynq.Task) error {
	var p fanoutPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	msg := p.Message
	log.Printf(" [*] Fanout message %s", p.Message.ID)

	_, err := h.Storage.NewMessage(context.Background(), &msg)
	if err != nil {
		log.Printf(" [ERROR] could not create a new message: %v", err)
		return err
	}

	subCollection, err := h.Storage.GetSubscriptions(ctx, msg.TenantID, msg.ApplicationID, nil)
	if err != nil {
		log.Printf(" [ERROR] Failed to get subscriptions: %v", err)
		return err
	}

	for subCollection.Results != nil && len(subCollection.Results) > 0 {
		for _, sub := range subCollection.Results {
			log.Printf(" [*] Enqueuing message for subscription %s", sub.ID)
			deliveryTask, err := NewDeliveryTask(msg, sub.Endpoint)
			if err != nil {
				return err
			}

			_, err = h.AsynqClient.EnqueueContext(ctx, deliveryTask)
			if err != nil {
				return err
			}
		}

		pageOpt := &Pagination{Token: subCollection.NextPageToken}
		subCollection, err = h.Storage.GetSubscriptions(ctx, msg.TenantID, msg.ApplicationID, pageOpt)
		if err != nil {
			log.Printf(" [ERROR] Failed to get subscriptions: %v", err)
			return err
		}
	}

	log.Printf(" [*] Completed fanout of message %s", p.Message.ID)
	return nil
}

type deliveryPayload struct {
	Message  Message
	Endpoint string
}

func NewDeliveryTask(m Message, endpoint string) (*asynq.Task, error) {
	payload, err := json.Marshal(deliveryPayload{Message: m, Endpoint: endpoint})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeDeliveryMessage, payload), nil
}

type DeliveryTaskHandler struct{}

func (h *DeliveryTaskHandler) Handle(ctx context.Context, t *asynq.Task) error {
	var p deliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	msg := p.Message
	endpoint := p.Endpoint

	log.Printf(" [*] Delivering message %s to %s", msg.ID, endpoint)

	return nil
}
