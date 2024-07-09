package handler

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"slices"

	"github.com/korayakpinar/network/src/crypto"
	"github.com/korayakpinar/network/src/mempool"
	"github.com/korayakpinar/network/src/message"
	"github.com/korayakpinar/network/src/types"
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
	var ourIndex uint64 = 0
	var CommitteSize uint64 = 512

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
			if h.mempool.GetThreshold(txHash) == h.mempool.GetPartialDecryptionCount(txHash) {
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

			}
		case message.MessageType_ENCRYPTED_BATCH:
			encBatchMsg := newMsg.Message.(*message.Message_EncryptedBatch).EncryptedBatch
			encBatch := &types.EncryptedBatch{
				Header: &types.BatchHeader{
					LeaderID:  encBatchMsg.Header.LeaderID,
					BlockNum:  encBatchMsg.Header.BlockNum,
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
			h.mempool.AddOrderSig(orderSigMsg.BlockNum, *orderSig)

		}
	}

}

func (h *Handler) HandleTransaction(tx string) error {
	// TODO: Get the current committee size from the smart contract
	var n uint64 = 256
	var t uint64 = 64

	randomIndexes := make([]uint64, t)
	for i := 0; i < int(t); i++ {
		randomIndexes[i] = uint64(rand.Intn(int(n))) // 0 dahil, 256 hariç
	}

	// TODO: Replace with GetPublicKeys function that gets the public keys of the given indexes
	pks := make([][]byte, n)
	for i := 0; i < int(n); i++ {
		pks[i] = []byte("pk" + string(randomIndexes[i]))
	}

	rawResp, err := crypto.EncryptTransaction([]byte(tx), pks, t, n)
	if err != nil {
		return err
	}

	var encResponse crypto.EncryptResponse
	err = proto.Unmarshal(rawResp, &encResponse)
	if err != nil {
		return err
	}

	// Construct the encrypted transaction
	encTxHeader := &types.EncryptedTxHeader{
		Hash:    fmt.Sprintf("%x", sha256.Sum256([]byte(tx))),
		GammaG2: encResponse.GammaG2,
		PkIDs:   randomIndexes,
	}

	encTxBody := &types.EncryptedTxBody{
		Sa1:       encResponse.Sa1,
		Sa2:       encResponse.Sa2,
		Iv:        encResponse.Iv,
		EncText:   encResponse.Enc,
		Threshold: t,
	}

	encTx := &types.EncryptedTransaction{
		Header: encTxHeader,
		Body:   encTxBody,
	}

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
		return err
	}

	err = h.topic.Publish(context.Background(), bytesMsg)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) Stop() {
	h.sub.Cancel()
}
