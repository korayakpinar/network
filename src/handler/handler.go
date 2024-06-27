package handler

import (
	"context"
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type Handler struct {
	sub *pubsub.Subscription
}

func NewHandler(sub *pubsub.Subscription) *Handler {
	return &Handler{sub: sub}
}

func (h *Handler) Start(ctx context.Context, errChan chan error) {

	for {
		msg, err := h.sub.Next(ctx)
		if err != nil {
			errChan <- err
			return
		}
		// TODO: Unmarshal the message with protobuf and update the price cache

		fmt.Println(msg.ReceivedFrom, ": ", string(msg.Message.Data))
	}

}

func (h *Handler) Stop() {
	h.sub.Cancel()
}
