package handler

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"slices"

	"github.com/ethereum/go-ethereum/common"
	"github.com/korayakpinar/network/src/crypto"
	"github.com/korayakpinar/network/src/mempool"
	"github.com/korayakpinar/network/src/message"
	"github.com/korayakpinar/network/src/types"
	"github.com/korayakpinar/network/src/utils"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"google.golang.org/protobuf/proto"
)

type Signer struct {
	address common.Address
	blsKey  []byte
}

type Handler struct {
	sub     *pubsub.Subscription
	topic   *pubsub.Topic
	mempool *mempool.Mempool
	signers *[]Signer
	crypto  *crypto.Crypto

	privKey       string
	rpcUrl        string
	committeeSize uint64
	ourIndex      uint64
	threshold     uint64
}

func NewHandler(sub *pubsub.Subscription, topic *pubsub.Topic, signers *[]Signer, privKey, rpcUrl, apiPort string, commmitteeSize, ourIndex, threshold uint64) *Handler {
	mempool := mempool.NewMempool()
	crypto := crypto.NewCrypto(apiPort)
	return &Handler{sub: sub, topic: topic, mempool: mempool, signers: signers, crypto: crypto, privKey: privKey, rpcUrl: rpcUrl, committeeSize: commmitteeSize, ourIndex: ourIndex, threshold: threshold}
}

func (h *Handler) Start(ctx context.Context, errChan chan error) {
	var leaderIndex uint64 = 0
	var ourIndex uint64 = h.ourIndex
	var CommitteSize uint64 = h.committeeSize

	go func() {
		if h.ourIndex == leaderIndex {
			// Leader
			// Submit a encrypted batch
			encTxs := h.mempool.GetTransactions()
			txHeaders := []*types.EncryptedTxHeader{}
			for _, tx := range encTxs {
				txHeaders = append(txHeaders, tx.Header)
			}

			encBatchBody := &types.BatchBody{
				EncTxs: txHeaders,
			}
			hashDigest := fmt.Sprintf("%x", sha256.Sum256(encBatchBody.Bytes()))
			sig, err := utils.SignTheHash(h.privKey, []byte(hashDigest))
			if err != nil {
				errChan <- err
				return
			}

			encBatchHeader := &types.BatchHeader{
				LeaderID:  leaderIndex,
				BlockNum:  0,
				Hash:      hashDigest,
				Signature: string(sig),
			}

			encBatch := &types.EncryptedBatch{
				Header: encBatchHeader,
				Body:   encBatchBody,
			}

			msg := &message.Message{
				Message: &message.Message_EncryptedBatch{
					EncryptedBatch: &message.EncryptedBatch{
						Header: &message.BatchHeader{
							LeaderID:  encBatch.Header.LeaderID,
							BlockNum:  encBatch.Header.BlockNum,
							Hash:      encBatch.Header.Hash,
							Signature: encBatch.Header.Signature,
						},
						Body: &message.BatchBody{
							Transactions: []*message.TransactionHeader{},
						},
					},
				},
			}

			msgBytes, err := proto.Marshal(msg)
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
	}()
	go func() {
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

					content, err := h.crypto.DecryptTransaction(encryptedContent, [][]byte{}, [][]byte{{}, {}}, []byte{}, []byte{}, []byte{}, 0, CommitteSize)
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
							Content: string(content),
						},
					})
					if h.ourIndex == leaderIndex {
						tx, err := sendRawTransaction(h.rpcUrl, string(content))
						if err != nil {
							errChan <- err
							return
						}
						fmt.Println("Transaction sent to the blockchain: ", tx)
					}
				}
			case message.MessageType_ENCRYPTED_BATCH:
				encBatchMsg := newMsg.Message.(*message.Message_EncryptedBatch).EncryptedBatch

				if encBatchMsg.Header.LeaderID != leaderIndex {
					continue // Ignore the message
				}

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
						partDec, err := h.crypto.PartialDecrypt(encTx.GammaG2)
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
				leaderIndex++
				leaderIndex %= CommitteSize
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
	}()

}

func (h *Handler) HandleTransaction(tx string) error {
	randomIndexes := make([]uint64, h.threshold)
	for i := 0; i < int(h.threshold); i++ {
		randomIndexes[i] = uint64(rand.Intn(int(h.committeeSize)))
	}

	pks := make([][]byte, h.committeeSize)
	ourSigners := *h.GetSigners()
	for i := 0; i < int(h.committeeSize); i++ {
		pks[i] = ourSigners[i].blsKey
	}

	rawResp, err := h.crypto.EncryptTransaction([]byte(tx), pks, h.threshold, h.committeeSize)
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
		Threshold: h.threshold,
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

func NewSigner(address common.Address, blsKey []byte) Signer {
	return Signer{address: address, blsKey: blsKey}
}

func (s *Signer) GetAddress() common.Address {
	return s.address
}

func (h *Handler) AddSigner(p Signer) {
	*h.signers = append(*h.signers, p)
}

func (h *Handler) GetSigners() *[]Signer {
	return h.signers
}

func (h *Handler) GetSignerByIndex(index uint64) *Signer {
	return &(*h.signers)[index]
}

/*
func (h *Handler) IsSigner(p peer.ID) bool {
	for _, pub := range *h.signers {
		if pub.peerID == p {
			return true
		}
	}
	return false
}
*/

func (h *Handler) Stop() {
	h.sub.Cancel()
}

type JSONRPCRequest struct {
	JsonRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type JSONRPCResponse struct {
	JsonRPC string        `json:"jsonrpc"`
	Result  string        `json:"result"`
	Error   *JSONRPCError `json:"error,omitempty"`
	ID      int           `json:"id"`
}

type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func sendRawTransaction(rpcUrl string, content string) (string, error) {
	// JSON-RPC isteği oluştur
	request := JSONRPCRequest{
		JsonRPC: "2.0",
		Method:  "eth_sendRawTransaction",
		Params:  []interface{}{content},
		ID:      1,
	}

	// İsteği JSON'a dönüştür
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("JSON encoding error: %v", err)
	}

	// HTTP POST isteği gönder
	resp, err := http.Post(rpcUrl, "application/json", bytes.NewBuffer(jsonRequest))
	if err != nil {
		return "", fmt.Errorf("HTTP POST error: %v", err)
	}
	defer resp.Body.Close()

	// Yanıtı oku
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Response reading error: %v", err)
	}

	// Yanıtı JSON'dan çöz
	var response JSONRPCResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("JSON decoding error: %v", err)
	}

	// Hata kontrolü
	if response.Error != nil {
		return "", fmt.Errorf("RPC error: %v", response.Error.Message)
	}

	// İşlem hash'ini döndür
	return response.Result, nil
}
