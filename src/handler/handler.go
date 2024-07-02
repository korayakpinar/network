package handler

import (
	"context"
	"slices"

	"github.com/korayakpinar/p2pclient/src/api"
	"github.com/korayakpinar/p2pclient/src/mempool"
	"github.com/korayakpinar/p2pclient/src/message"
	"github.com/korayakpinar/p2pclient/src/types"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"google.golang.org/protobuf/proto"
)

type Handler struct {
	sub     *pubsub.Subscription
	topic   *pubsub.Topic
	mempool *mempool.Mempool
}

func NewHandler(sub *pubsub.Subscription, topic *pubsub.Topic) *Handler {
	mempool := mempool.NewMempool()
	return &Handler{sub: sub, topic: topic, mempool: mempool}
}

func (h *Handler) Start(ctx context.Context, errChan chan error) {
	//TODO: Get the array that contains the order of the leaders/validators from the smart contract
	//TODO: Get the node's public key index from smart contract
	var ourIndex uint32 = 0

	for {
		msg, err := h.sub.Next(ctx)
		if err != nil {
			errChan <- err
			return
		}

		// Deserialize the proto message
		newMsg := &message.Message{}
		err = proto.Unmarshal(msg.Data, newMsg)
		if err != nil {
			errChan <- err
			return
		}

		switch newMsg.MessageType {
		case message.MessageType_ENCRYPTED_TRANSACTION:
			encTxMsg := newMsg.Message.(*message.Message_EncryptedTransaction).EncryptedTransaction
			// Maybe functionize these conversions?
			encTx := &types.EncryptedTransaction{
				Header: &types.TransactionHeader{
					Hash:    types.TxHash(encTxMsg.Header.Hash),
					GammaG2: encTxMsg.Header.GammaG2,
					PkIDs:   encTxMsg.Header.PkIDs,
				},
				Body: &types.TransactionBody{
					Sa1:       encTxMsg.Body.Sa1,
					Sa2:       encTxMsg.Body.Sa2,
					Iv:        encTxMsg.Body.Iv,
					EncText:   encTxMsg.Body.EncText,
					Threshold: encTxMsg.Body.T,
				},
			}
			h.mempool.AddTransaction(encTx)
		case message.MessageType_PARTIAL_DECRYPTION:
			partDecMsg := newMsg.Message.(*message.Message_PartialDecryption).PartialDecryption
			partDec := types.PartialDecryption(partDecMsg.PartDec)
			txHash := types.TxHash(partDecMsg.TxHash)
			h.mempool.AddPartialDecryption(txHash, partDec)
		case message.MessageType_ENCRYPTED_BATCH:
			encBatchMsg := newMsg.Message.(*message.Message_EncryptedBatch).EncryptedBatch
			encBatch := &types.EncryptedBatch{
				Header: &types.BatchHeader{
					LeaderID:  encBatchMsg.Header.LeaderID,
					BlockNum:  types.BlockNum(encBatchMsg.Header.BlockNum),
					Hash:      encBatchMsg.Header.Hash,
					Signature: encBatchMsg.Header.Signature,
				},
				Body: &types.BatchBody{
					EncTxs: []*types.TransactionHeader{},
				},
			}
			var txHashes []types.TxHash
			for _, encTx := range encBatchMsg.Body.Transactions {
				txHashes = append(txHashes, types.TxHash(encTx.Hash))
				if slices.Contains(encTx.PkIDs, ourIndex) {
					partDec, err := api.PartialDecrypt(encTx.GammaG2)
					if err != nil {
						errChan <- err
						return
					}

					newMessage := &message.Message{
						Message: &message.Message_PartialDecryption{
							PartialDecryption: &message.PartialDecryption{
								TxHash:  encTx.Hash,
								PartDec: partDec,
							},
						},
					}

					msgBytes, err := proto.Marshal(newMessage)
					if err != nil {
						errChan <- err
						return
					}
					err = h.topic.Publish(ctx, msgBytes)
					if err != nil {
						errChan <- err
						return
					}
				}
				newTx := &types.TransactionHeader{
					Hash:    types.TxHash(encTx.Hash),
					GammaG2: encTx.GammaG2,
				}
				encBatch.Body.EncTxs = append(encBatch.Body.EncTxs, newTx)

			}
			h.mempool.RemoveTransactions(txHashes)
		case message.MessageType_ORDER_SIGNATURE:
			orderSigMsg := newMsg.Message.(*message.Message_OrderSignature).OrderSignature
			orderSig := &types.OrderSig{
				Signature: orderSigMsg.Signature,
				TxHeaders: []*types.TransactionHeader{},
			}
			for _, tx := range orderSigMsg.Order {
				newTx := &types.TransactionHeader{
					Hash:    types.TxHash(tx.Hash),
					GammaG2: tx.GammaG2,
				}
				orderSig.TxHeaders = append(orderSig.TxHeaders, newTx)
			}
			h.mempool.AddOrderSig(types.BlockNum(orderSigMsg.BlockNum), *orderSig)

		}

		/*
			// Just an Example
			fmt.Println(msg.ReceivedFrom, ": ", string(msg.Message.Data))


			tx := &message.EncryptedTransaction{
				Header: &message.TransactionHeader{
					Hash:    "hash",
					GammaG2: "gammaG2",
				},
				Body: &message.TransactionBody{
					PkIDs: []uint32{1, 2, 3},
					Sa1:   []string{"sa1"},
					Sa2:   []string{"sa2"},
					Iv:    []byte{1, 2, 3},
					T:     5,
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
			if err != nil {
				errChan <- err
				return
			} */
	}

}

func (h *Handler) Stop() {
	h.sub.Cancel()
}
