package captainhook

import (
	"net/http"
	"time"
)

type Deliverer struct {
	client http.Client
}

func newDeliverer() *Deliverer {
	return &Deliverer{
		client: http.Client{
			Timeout: time.Duration(1) * time.Second,
		},
	}
}

func (d Deliverer) Deliver(m Message, endpoint string) error {
	return nil
}
