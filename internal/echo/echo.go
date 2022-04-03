package echo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/belljustin/captainhook/captainhook"
)

type EchoServer struct {
	port int
}

func New(port int) *EchoServer {
	return &EchoServer{port: port}
}

func (s *EchoServer) Run() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// TODO print message received
		decoder := json.NewDecoder(r.Body)
		var p captainhook.DeliveryPayload
		err := decoder.Decode(&p)
		if err != nil {
			fmt.Printf(" [ERROR] could not decode body: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Printf("%+v", p)
		fmt.Fprint(w, "Success")
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}
