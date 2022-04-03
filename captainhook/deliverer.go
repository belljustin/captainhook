package captainhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Deliverer interface {
	Deliver(m Message, endpoint string) error
}

type httpDeliverer struct {
	client http.Client
}

func NewHttpDeliverer() *httpDeliverer {
	return &httpDeliverer{
		client: http.Client{
			Timeout: time.Duration(1) * time.Second,
		},
	}
}

type DeliveryPayload struct {
	Type string
	Data []byte
}

func (d httpDeliverer) Deliver(m Message, endpoint string) error {
	p, err := json.Marshal(DeliveryPayload{Type: m.Type, Data: m.Data})
	if err != nil {
		return err
	}

	resp, err := d.client.Post(endpoint, "application/json", bytes.NewBuffer(p))
	if err != nil {
		return err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("subscription returned with non 2xx response: %d", resp.StatusCode)
	}

	return nil
}
