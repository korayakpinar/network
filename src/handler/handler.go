package handler

import (
	"context"
	"slices"

	"github.com/korayakpinar/p2pclient/src/crypto"
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
	var CommitteSize uint32 = 512

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
				Header: &types.EncryptedTxHeader{
					Hash:    string(encTxMsg.Header.Hash),
					GammaG2: encTxMsg.Header.GammaG2,
					PkIDs:   encTxMsg.Header.PkIDs,
				},
				Body: &types.EncryptedTxBody{
					Sa1:       encTxMsg.Body.Sa1,
					Sa2:       encTxMsg.Body.Sa2,
					Iv:        encTxMsg.Body.Iv,
					EncText:   encTxMsg.Body.EncText,
					Threshold: encTxMsg.Body.T,
				},
			}
			h.mempool.AddEncryptedTx(encTx)
		case message.MessageType_PARTIAL_DECRYPTION:
			partDecMsg := newMsg.Message.(*message.Message_PartialDecryption).PartialDecryption
			partDec := partDecMsg.PartDec
			txHash := partDecMsg.TxHash
			h.mempool.AddPartialDecryption(txHash, &partDec)
			if h.mempool.GetThreshold(txHash) == uint32(h.mempool.GetPartialDecryptionCount(txHash)) {
				encTx := h.mempool.GetTransaction(txHash)
				encryptedContent := encTx.Body.EncText

				content, err := crypto.DecryptTransaction(encryptedContent, [][]byte{}, [][]byte{{}, {}}, []byte{}, []byte{}, []byte{}, 0, CommitteSize)
				if err != nil {
					errChan <- err
					return
				}
				h.mempool.AddDecryptedTx(&types.DecryptedTransaction{
					Header: &types.DecryptedTxHeader{
						Hash:  txHash,
						PkIDs: encTx.Header.PkIDs,
					},
					Body: &types.DecryptedTxBody{
						Content: content,
					},
				})

				/* newMessage := &message.Message{
					Message: &message.Me{
						DecryptedTransaction: &message.DecryptedTransaction{
							TxHash:  string(txHash),
							DecText: decTx,
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
				} */

			}
		case message.MessageType_ENCRYPTED_BATCH:
			encBatchMsg := newMsg.Message.(*message.Message_EncryptedBatch).EncryptedBatch
			encBatch := &types.EncryptedBatch{
				Header: &types.BatchHeader{
					LeaderID:  encBatchMsg.Header.LeaderID,
					BlockNum:  uint32(encBatchMsg.Header.BlockNum),
					Hash:      encBatchMsg.Header.Hash,
					Signature: encBatchMsg.Header.Signature,
				},
				Body: &types.BatchBody{
					EncTxs: []*types.EncryptedTxHeader{},
				},
			}
			var txHashes []string
			for _, encTx := range encBatchMsg.Body.Transactions {
				txHashes = append(txHashes, string(encTx.Hash))
				if slices.Contains(encTx.PkIDs, ourIndex) {
					partDec, err := crypto.PartialDecrypt(encTx.GammaG2)
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
				newTx := &types.EncryptedTxHeader{
					Hash:    string(encTx.Hash),
					GammaG2: encTx.GammaG2,
				}
				encBatch.Body.EncTxs = append(encBatch.Body.EncTxs, newTx)

			}
			h.mempool.RemoveTransactions(txHashes)
		case message.MessageType_ORDER_SIGNATURE:
			orderSigMsg := newMsg.Message.(*message.Message_OrderSignature).OrderSignature
			orderSig := &types.OrderSig{
				Signature: orderSigMsg.Signature,
				TxHeaders: []*types.EncryptedTxHeader{},
			}
			for _, tx := range orderSigMsg.Order {
				newTx := &types.EncryptedTxHeader{
					Hash:    string(tx.Hash),
					GammaG2: tx.GammaG2,
				}
				orderSig.TxHeaders = append(orderSig.TxHeaders, newTx)
			}
			h.mempool.AddOrderSig(uint32(orderSigMsg.BlockNum), *orderSig)

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

func (h *Handler) HandleTransaction(encTx *types.EncryptedTransaction) {
	h.mempool.AddEncryptedTx(encTx)

	msg := &message.Message{
		Message: &message.Message_EncryptedTransaction{
			EncryptedTransaction: &message.EncryptedTransaction{
				Header: &message.TransactionHeader{
					Hash:    encTx.Header.Hash,
					GammaG2: encTx.Header.GammaG2,
					PkIDs:   encTx.Header.PkIDs,
				},
				Body: &message.TransactionBody{
					Sa1:     encTx.Body.Sa1,
					Sa2:     encTx.Body.Sa2,
					Iv:      encTx.Body.Iv,
					EncText: encTx.Body.EncText,
					T:       encTx.Body.Threshold,
				},
			},
		},
	}
	bytesMsg, err := proto.Marshal(msg)
	if err != nil {
		return
	}

	h.topic.Publish(context.Background(), bytesMsg)
}

func (h *Handler) Stop() {
	h.sub.Cancel()
}
