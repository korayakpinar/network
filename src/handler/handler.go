package handler

import (
	"context"
	"fmt"

	"github.com/korayakpinar/p2pclient/src/price"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type Handler struct {
	sub        *pubsub.Subscription
	priceCache *price.Cache
}

func NewHandler(sub *pubsub.Subscription, priceCache *price.Cache) *Handler {
	return &Handler{sub: sub, priceCache: priceCache}
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
