package handler

import (
	"context"
	"fmt"

	"github.com/korayakpinar/p2pclient/src/message"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"google.golang.org/protobuf/proto"
)

type Handler struct {
	sub   *pubsub.Subscription
	topic *pubsub.Topic
}

func NewHandler(sub *pubsub.Subscription, topic *pubsub.Topic) *Handler {
	return &Handler{sub: sub, topic: topic}
}

func (h *Handler) Start(ctx context.Context, errChan chan error) {

	for {
		msg, err := h.sub.Next(ctx)
		if err != nil {
			errChan <- err
			return
		}

		fmt.Println(msg.ReceivedFrom, ": ", string(msg.Message.Data))

		// Just an Example
		tx := &message.EncryptedTransaction{
			Header: &message.TransactionHeader{
				Hash:    "hash",
				GammaG2: "gammaG2",
			},
			Body: &message.TransactionBody{
				DecryptParams: "decryptParams",
				PkIDs:         []uint32{1, 2, 3},
			},
		}
		newMessage := &message.Message{
			Message: &message.Message_EncryptedTransaction{
				EncryptedTransaction: tx,
			},
			MessageType: message.MessageType_ENCRYPTED_TRANSACTION,
		}

		msgBytes, err := proto.Marshal(newMessage)

		if err != nil {
			errChan <- err
			return
		}
		err = h.topic.Publish(ctx, msgBytes)
	}

}

func (h *Handler) Stop() {
	h.sub.Cancel()
}
