package handler

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/korayakpinar/network/src/crypto"
	"github.com/korayakpinar/network/src/mempool"
	"github.com/korayakpinar/network/src/message"
	"github.com/korayakpinar/network/src/proxy"
	"github.com/korayakpinar/network/src/types"
	"github.com/korayakpinar/network/src/utils"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	libp2pCrypto "github.com/libp2p/go-libp2p/core/crypto"
	"google.golang.org/protobuf/proto"
)

type Handler struct {
	sub     *pubsub.Subscription
	topic   *pubsub.Topic
	mempool *mempool.Mempool
	proxy   *proxy.Proxy
	signers map[uint64]*types.Signer
	crypto  *crypto.Crypto
	privKey libp2pCrypto.PrivKey

	networkIndex  uint64
	networkSize   uint64
	rpcUrl        string
	committeeSize uint64
	threshold     uint64
}

func NewHandler(sub *pubsub.Subscription, topic *pubsub.Topic, proxy *proxy.Proxy, mempool *mempool.Mempool, signers map[uint64]*types.Signer, privKey libp2pCrypto.PrivKey, rpcUrl, apiPort string, commmitteeSize, threshold, networkIndex, networkSize uint64) *Handler {
	crypto := crypto.NewCrypto(apiPort)
	return &Handler{
		sub:           sub,
		topic:         topic,
		mempool:       mempool,
		proxy:         proxy,
		signers:       signers,
		crypto:        crypto,
		privKey:       privKey,
		rpcUrl:        rpcUrl,
		committeeSize: commmitteeSize,
		threshold:     threshold,
		networkIndex:  networkIndex,
		networkSize:   networkSize,
	}

}

func (h *Handler) Start(ctx context.Context, errChan chan error) {
	log.Println("Starting the handler")
	var leaderIndex uint64 = 0

	log.Printf("Initial state: ourIndex=%d, leaderIndex=%d, CommitteeSize=%d", h.networkIndex, leaderIndex, h.committeeSize)

	go h.leaderRoutine(ctx, errChan, &leaderIndex)
	go h.messageHandlingRoutine(ctx, errChan, &leaderIndex, h.committeeSize)

	log.Println("Handler started successfully")
}

func (h *Handler) leaderRoutine(ctx context.Context, errChan chan error, leaderIndex *uint64) {
	log.Println("Starting leader routine")
	for {
		select {
		case <-ctx.Done():
			log.Println("Context cancelled, stopping leader routine")
			return
		default:

			if h.networkIndex == *leaderIndex {
				if err := h.performLeaderDuties(ctx, leaderIndex); err != nil {
					log.Printf("Error performing leader duties: %v", err)
					errChan <- err
					return
				}
			}
			time.Sleep(time.Second / 8) // Prevent tight loop
		}
	}
}

func (h *Handler) performLeaderDuties(ctx context.Context, leaderIndex *uint64) error {
	log.Println("Performing leader duties")
	encTxs := h.mempool.GetPendingTransactions()
	if len(encTxs) == 0 {
		log.Println("No transactions to submit")
		return nil
	}
	log.Printf("Preparing to submit %d transactions", len(encTxs))

	encBatch, err := h.prepareEncryptedBatch(encTxs)
	if err != nil {
		return fmt.Errorf("failed to prepare encrypted batch: %w", err)
	}

	if err := h.publishEncryptedBatch(ctx, encBatch); err != nil {
		return fmt.Errorf("failed to publish encrypted batch: %w", err)
	}
	//time.Sleep(3 * time.Second) // Wait for other nodes to prepare partial decryptions

	log.Println("Leader duties completed successfully")

	*leaderIndex++
	*leaderIndex %= h.networkSize

	return nil
}

func (h *Handler) prepareEncryptedBatch(encTxs []*types.Transaction) (*types.EncryptedBatch, error) {
	log.Println("Preparing encrypted batch")
	txHeaders := []*types.EncryptedTxHeader{}
	for _, tx := range encTxs {
		txHeaders = append(txHeaders, tx.EncryptedTransaction.Header)
	}

	encBatchBody := &types.BatchBody{
		EncTxs: txHeaders,
	}

	hashBytes := sha256.Sum256(encBatchBody.Bytes())
	hashDigest := hashBytes[:] // Bu, 32 byte uzunluğunda bir []byte olacak

	sig, err := h.privKey.Sign(hashDigest)
	if err != nil {
		return nil, fmt.Errorf("failed to sign batch: %w", err)
	}

	encBatchHeader := &types.BatchHeader{
		LeaderID:  h.networkIndex,
		BlockNum:  0,
		Hash:      hex.EncodeToString(hashDigest),
		Signature: hex.EncodeToString(sig),
	}

	encBatch := &types.EncryptedBatch{
		Header: encBatchHeader,
		Body:   encBatchBody,
	}

	log.Println("Encrypted batch prepared successfully")
	return encBatch, nil
}

func (h *Handler) publishEncryptedBatch(ctx context.Context, encBatch *types.EncryptedBatch) error {
	log.Println("Publishing encrypted batch")
	msg := &message.Message{
		MessageType: message.MessageType_ENCRYPTED_BATCH,
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

	var txHashes []string
	for _, tx := range encBatch.Body.EncTxs {
		msg.Message.(*message.Message_EncryptedBatch).EncryptedBatch.Body.Transactions = append(
			msg.Message.(*message.Message_EncryptedBatch).EncryptedBatch.Body.Transactions,
			&message.TransactionHeader{
				Hash:    string(tx.Hash),
				GammaG2: tx.GammaG2,
				PkIDs:   tx.PkIDs,
			},
		)
		txHashes = append(txHashes, string(tx.Hash))
	}

	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	h.mempool.SetMultipleTransactionsProposed(txHashes)
	for _, txHash := range txHashes {
		incTx := h.mempool.GetTransaction(txHash).EncryptedTransaction
		if incTx != nil {
			log.Printf("Transaction included: %s", txHash)
			for _, id := range incTx.Header.PkIDs {
				if h.signers[id].IsLocal {
					signer := h.signers[id]
					partDec, err := h.crypto.PartialDecrypt([]byte(signer.KeyPair.BLSPrivateKey), incTx.Header.GammaG2)
					if err != nil {
						return fmt.Errorf("failed to partially decrypt transaction: %w", err)
					}
					h.mempool.AddPartialDecryption(txHash, id, partDec)
				}
			}

		}
	}

	err = h.topic.Publish(ctx, msgBytes)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Println("Encrypted batch published successfully")
	return nil
}

func (h *Handler) messageHandlingRoutine(ctx context.Context, errChan chan error, leaderIndex *uint64, CommitteeSize uint64) {
	log.Println("Starting message handling routine")
	for {
		select {
		case <-ctx.Done():
			log.Println("Context cancelled, stopping message handling routine")
			return
		default:
			time.Sleep(time.Second / 8) // Prevent tight loop
			if err := h.handleNextMessage(ctx, leaderIndex, CommitteeSize); err != nil {
				log.Printf("Error handling message: %v", err)
				errChan <- err
				return
			}
		}
	}
}

func (h *Handler) handleNextMessage(ctx context.Context, leaderIndex *uint64, CommitteeSize uint64) error {
	log.Println("Waiting for next message")
	msg, err := h.sub.Next(ctx)
	if err != nil {
		return fmt.Errorf("failed to get next message: %w", err)
	}

	newMsg := &message.Message{}
	if err := proto.Unmarshal(msg.Data, newMsg); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	log.Printf("Received message of type: %s", newMsg.MessageType)
	switch newMsg.MessageType {
	case message.MessageType_ENCRYPTED_TRANSACTION:
		return h.handleEncryptedTransaction(newMsg)
	case message.MessageType_PARTIAL_DECRYPTION_BATCH:
		return h.handlePartialDecryptionBatch(newMsg, *leaderIndex, CommitteeSize)
	case message.MessageType_ENCRYPTED_BATCH:
		return h.handleEncryptedBatch(ctx, newMsg, leaderIndex)
	case message.MessageType_ORDER_SIGNATURE:
		return h.handleOrderSignature(newMsg)
	default:
		log.Printf("Unknown message type: %s", newMsg.MessageType)
	}
	return nil
}

func (h *Handler) handleEncryptedTransaction(msg *message.Message) error {
	log.Println("Handling encrypted transaction")
	encTxMsg := msg.Message.(*message.Message_EncryptedTransaction).EncryptedTransaction
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
	log.Printf("Encrypted transaction received: %s", encTx.Header.Hash)
	newTx := types.NewTransaction(nil, encTx.Header.Hash)
	newTx.SetEncrypted(encTx)
	h.mempool.AddTransaction(newTx)
	return nil
}

func (h *Handler) handlePartialDecryptionBatch(msg *message.Message, leaderIndex uint64, CommitteeSize uint64) error {

	partDecMsg := msg.Message.(*message.Message_PartialDecryptionBatch).PartialDecryptionBatch
	decryptions := partDecMsg.Decryptions
	sender := partDecMsg.Sender
	txHash := partDecMsg.TxHash
	log.Println("Handling partial decryption from: ", sender)

	// Skip if we are not the first node in the committee or partial decryptions from us
	if sender == h.networkIndex || h.networkIndex != 0 {
		log.Printf("Ignoring partial decryption from us or we are not a leader: %d", sender)
		return nil
	}

	for _, dec := range decryptions {
		h.mempool.AddPartialDecryption(txHash, dec.Signer, dec.PartDec)
	}

	if int(h.mempool.GetThreshold(txHash)) < h.mempool.GetPartialDecryptionCount(txHash) && !h.mempool.CheckTransactionDecrypted(txHash) {
		log.Printf("All partial decryptions received for: %s", txHash)
		encTx := h.mempool.GetTransaction(txHash).EncryptedTransaction
		encryptedContent := encTx.Body.EncText

		/*
			pks := make([][]byte, CommitteeSize)
			for i := 0; i < int(CommitteeSize); i++ {
				pks[i] = h.GetSignerByIndex(uint64(i)).BLSPublicKey
			}
		*/

		partDecs := h.mempool.GetPartialDecryptions(txHash)

		// TODO: Seperate the committee size and the nodes count, for now we are using the same number
		content, err := h.crypto.DecryptTransaction(encryptedContent, partDecs, encTx.Header.GammaG2, encTx.Body.Sa1, encTx.Body.Sa2, encTx.Body.Iv, encTx.Body.Threshold, CommitteeSize+1)

		if err != nil {
			return fmt.Errorf("failed to decrypt transaction: %w", err)
		}
		h.mempool.SetTransactionDecrypted(txHash, &types.DecryptedTransaction{
			Header: &types.DecryptedTxHeader{
				Hash:  txHash,
				PkIDs: encTx.Header.PkIDs,
			},
			Body: &types.DecryptedTxBody{
				Content: string(content),
			},
		})
		log.Printf("Transaction decrypted: %s", txHash)
		if h.networkIndex == 0 {
			tx, err := sendRawTransaction(h.rpcUrl, string(content))
			if err != nil {
				return fmt.Errorf("failed to send transaction to blockchain: %w", err)
			}
			h.mempool.SetTransactionIncluded(txHash)

			log.Printf("Transaction sent to the blockchain: %s", tx)
		}
	}
	return nil
}

func (h *Handler) handleEncryptedBatch(ctx context.Context, msg *message.Message, leaderIndex *uint64) error {

	log.Println("Handling encrypted batch")
	encBatchMsg := msg.Message.(*message.Message_EncryptedBatch).EncryptedBatch

	var txHashes []string
	for _, encTx := range encBatchMsg.Body.Transactions {
		txHashes = append(txHashes, string(encTx.Hash))
		partDecs := make([]*message.PartialDecryption, 0)

		for _, id := range encTx.PkIDs {
			if h.signers[id].IsLocal {
				signer := h.signers[id]
				partDec, err := h.crypto.PartialDecrypt([]byte(signer.KeyPair.BLSPrivateKey), encTx.GammaG2)
				if err != nil {
					return fmt.Errorf("failed to partially decrypt transaction: %w", err)
				}
				h.mempool.AddPartialDecryption(string(encTx.Hash), id, partDec)

				fmt.Printf("Partial decryption added for transaction: %s\n", encTx.Hash)

				partDecMsg := &message.PartialDecryption{
					Signer:  id,
					PartDec: partDec,
				}

				partDecs = append(partDecs, partDecMsg)

			}
		}

		if len(partDecs) > 0 {
			batchMessage := &message.Message{
				MessageType: message.MessageType_PARTIAL_DECRYPTION_BATCH,
				Message: &message.Message_PartialDecryptionBatch{
					PartialDecryptionBatch: &message.PartialDecryptionBatch{
						Decryptions: partDecs,
						Sender:      h.networkIndex,
						TxHash:      string(encTx.Hash),
					},
				},
			}
			msgBytes, err := proto.Marshal(batchMessage)
			if err != nil {
				return fmt.Errorf("failed to marshal partial decryption message: %w", err)
			}
			err = h.topic.Publish(ctx, msgBytes)
			if err != nil {
				return fmt.Errorf("failed to publish partial decryption message: %w", err)
			}
		}

	}

	if encBatchMsg.Header.LeaderID == h.networkIndex {
		log.Println("Ignoring own batch")
		return nil
	}

	if encBatchMsg.Header.LeaderID != *leaderIndex {
		log.Printf("Ignoring batch from non-leader (received: %d, expected: %d)", encBatchMsg.Header.LeaderID, *leaderIndex)
		return nil
	}

	h.mempool.SetMultipleTransactionsProposed(txHashes)

	*leaderIndex++
	*leaderIndex %= h.networkSize
	log.Printf("New leader index: %d", *leaderIndex)
	return nil
}

func (h *Handler) handleOrderSignature(msg *message.Message) error {
	log.Println("Handling order signature")
	orderSigMsg := msg.Message.(*message.Message_OrderSignature).OrderSignature
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

	log.Printf("Order signature added for block number: %d", orderSigMsg.BlockNum)
	return nil
}

func (h *Handler) HandleTransaction(tx string) error {
	randomIndexes := make([]uint64, h.committeeSize)
	// TODO: Uncomment this section and randomly select the committee members
	for i := 0; i < int(h.committeeSize); i++ {
		/* randNum := uint64(rand.Intn(int(h.committeeSize)))
		if !slices.Contains(randomIndexes, randNum) {
			randomIndexes[i] = randNum
		} else {
			i--
		} */
		randomIndexes[i] = uint64(i)
	}

	randThreshold := uint64(rand.Intn(int(h.committeeSize)-3)) + 1

	log.Println("Random indexes: ", randomIndexes)

	/*
		pks := make([][]byte, h.committeeSize)

		for index, signer := range h.GetSigners() {
			pks[index] = signer.BLSPublicKey
		}
	*/

	encResponse, err := h.crypto.EncryptTransaction([]byte(tx), randThreshold, h.committeeSize)
	if err != nil {
		log.Println("Error while encrypting the transaction: ", err)
		return err
	}

	txHash, err := utils.CalculateTxHash(tx)
	if err != nil {
		log.Println("Error while calculating the transaction hash: ", err)
		return err
	}

	log.Printf("Transaction encrypted successfully, plaintext tx: %s, encrypted tx: %s, tx hash: %s", tx, encResponse.Enc, txHash)

	// Construct the encrypted transaction
	encTxHeader := &types.EncryptedTxHeader{
		Hash:    txHash,
		GammaG2: encResponse.GammaG2,
		PkIDs:   randomIndexes,
	}

	encTxBody := &types.EncryptedTxBody{
		Sa1:       encResponse.Sa1,
		Sa2:       encResponse.Sa2,
		Iv:        encResponse.Iv,
		EncText:   encResponse.Enc,
		Threshold: randThreshold,
	}

	encTx := &types.EncryptedTransaction{
		Header: encTxHeader,
		Body:   encTxBody,
	}

	newTx := types.NewTransaction(&tx, txHash)
	newTx.SetEncrypted(encTx)
	h.mempool.AddTransaction(newTx)

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
		MessageType: message.MessageType_ENCRYPTED_TRANSACTION,
	}
	log.Println("Threshold number of message: ", encTx.Body.Threshold)
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

func (h *Handler) HandleEncryptedTransaction(encTx *types.EncryptedTransaction) error {
	log.Printf("Handling encrypted transaction with hash %s and threshold %d", encTx.Header.Hash, encTx.Body.Threshold)

	newTx := types.NewTransaction(nil, encTx.Header.Hash)
	newTx.SetEncrypted(encTx)
	h.mempool.AddTransaction(newTx)

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
		MessageType: message.MessageType_ENCRYPTED_TRANSACTION,
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

func (h *Handler) GetPartialDecryptionCount(hash string) int {
	tx := h.mempool.GetTransaction(hash)
	if tx != nil {
		return tx.PartialDecryptionCount
	}
	return 0
}

func (h *Handler) GetMetadataOfTx(hash string) (committeeSize, threshold int) {
	tx := h.mempool.GetTransaction(hash)
	if tx != nil {
		return tx.CommitteeSize, tx.Threshold
	}
	return 0, 0
}

func NewLocalSigner(index uint64, address common.Address, blsKey []byte) types.Signer {
	return types.Signer{Index: index, Address: address, BLSPublicKey: blsKey, IsLocal: false, KeyPair: nil}
}

func (h *Handler) AddSigner(p types.Signer, index uint64) {
	h.signers[index] = &p
}

func (h *Handler) GetSigners() map[uint64]*types.Signer {
	return h.signers
}

func (h *Handler) GetSignerByIndex(index uint64) *types.Signer {
	return h.signers[index]
}

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
	// Create the JSON-RPC request
	request := JSONRPCRequest{
		JsonRPC: "2.0",
		Method:  "eth_sendRawTransaction",
		Params:  []interface{}{content},
		ID:      1,
	}

	// Convert the request to JSON
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("JSON encoding error: %v", err)
	}

	// Send the request
	resp, err := http.Post(rpcUrl, "application/json", bytes.NewBuffer(jsonRequest))
	if err != nil {
		return "", fmt.Errorf("HTTP POST error: %v", err)
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Response reading error: %v", err)
	}

	// Decode the JSON response
	var response JSONRPCResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("JSON decoding error: %v", err)
	}

	fmt.Println("Response: ", response)

	// Check for errors
	if response.Error != nil && response.Error.Message != "already known" {
		return "", fmt.Errorf("RPC error: %v", response.Error.Message)
	}

	// Return the transaction hash
	return response.Result, nil
}

func (h *Handler) BroadcastNewTransaction(tx *types.Transaction) {
	h.proxy.BroadcastNewTransaction(tx)
}

func (h *Handler) BroadcastTransactionUpdate(tx *types.Transaction) {
	h.proxy.BroadcastTransactionUpdate(tx)
}
