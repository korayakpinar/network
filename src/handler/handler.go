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
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
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
	log.Println("Starting the handler")
	var leaderIndex uint64 = 0
	var ourIndex uint64 = h.ourIndex
	var CommitteSize uint64 = h.committeeSize

	log.Printf("Initial state: ourIndex=%d, leaderIndex=%d, CommitteSize=%d", ourIndex, leaderIndex, CommitteSize)

	go h.leaderRoutine(ctx, errChan, &leaderIndex, ourIndex)
	go h.messageHandlingRoutine(ctx, errChan, &leaderIndex, ourIndex, CommitteSize)

	log.Println("Handler started successfully")
}

func (h *Handler) leaderRoutine(ctx context.Context, errChan chan error, leaderIndex *uint64, ourIndex uint64) {
	log.Println("Starting leader routine")
	for {
		select {
		case <-ctx.Done():
			log.Println("Context cancelled, stopping leader routine")
			return
		default:
			log.Printf("Leader check: ourIndex=%d, leaderIndex=%d", ourIndex, *leaderIndex)
			if ourIndex == *leaderIndex {
				if err := h.performLeaderDuties(ctx, leaderIndex); err != nil {
					log.Printf("Error performing leader duties: %v", err)
					errChan <- err
					return
				}
			}
			time.Sleep(time.Second) // Prevent tight loop
		}
	}
}

func (h *Handler) performLeaderDuties(ctx context.Context, leaderIndex *uint64) error {
	log.Println("Performing leader duties")
	encTxs := h.mempool.GetEncryptedTransactions()
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
	time.Sleep(3 * time.Second) // Wait for other nodes to prepare partial decryptions

	log.Println("Leader duties completed successfully")

	*leaderIndex++
	*leaderIndex %= h.committeeSize

	return nil
}

func (h *Handler) prepareEncryptedBatch(encTxs []*types.EncryptedTransaction) (*types.EncryptedBatch, error) {
	log.Println("Preparing encrypted batch")
	txHeaders := []*types.EncryptedTxHeader{}
	for _, tx := range encTxs {
		txHeaders = append(txHeaders, tx.Header)
	}

	encBatchBody := &types.BatchBody{
		EncTxs: txHeaders,
	}

	hashBytes := sha256.Sum256(encBatchBody.Bytes())
	hashDigest := hashBytes[:] // Bu, 32 byte uzunluğunda bir []byte olacak

	sig, err := utils.SignTheHash(h.privKey, hashDigest)
	if err != nil {
		return nil, fmt.Errorf("failed to sign batch: %w", err)
	}

	encBatchHeader := &types.BatchHeader{
		LeaderID:  h.ourIndex,
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

	h.mempool.IncludeEncryptedTxs(txHashes)
	for _, txHash := range txHashes {
		incTx := h.mempool.GetIncludedTransaction(txHash)
		if incTx != nil {
			log.Printf("Transaction included: %s", txHash)
			if slices.Contains(incTx.Header.PkIDs, h.ourIndex) {
				partDec, err := h.crypto.PartialDecrypt(incTx.Header.GammaG2)
				if err != nil {
					return fmt.Errorf("failed to partially decrypt transaction: %w", err)
				}
				h.mempool.AddPartialDecryption(txHash, h.ourIndex, partDec)

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

func (h *Handler) messageHandlingRoutine(ctx context.Context, errChan chan error, leaderIndex *uint64, ourIndex, CommitteSize uint64) {
	log.Println("Starting message handling routine")
	for {
		select {
		case <-ctx.Done():
			log.Println("Context cancelled, stopping message handling routine")
			return
		default:
			if err := h.handleNextMessage(ctx, leaderIndex, ourIndex, CommitteSize); err != nil {
				log.Printf("Error handling message: %v", err)
				errChan <- err
				return
			}
		}
	}
}

func (h *Handler) handleNextMessage(ctx context.Context, leaderIndex *uint64, ourIndex, CommitteSize uint64) error {
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
	case message.MessageType_PARTIAL_DECRYPTION:
		return h.handlePartialDecryption(newMsg, *leaderIndex, CommitteSize)
	case message.MessageType_ENCRYPTED_BATCH:
		return h.handleEncryptedBatch(ctx, newMsg, leaderIndex, ourIndex, CommitteSize)
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
	h.mempool.AddEncryptedTx(encTx)
	return nil
}

func (h *Handler) handlePartialDecryption(msg *message.Message, leaderIndex uint64, CommitteSize uint64) error {
	log.Println("Handling partial decryption")
	partDecMsg := msg.Message.(*message.Message_PartialDecryption).PartialDecryption
	partDec := partDecMsg.PartDec
	sender := partDecMsg.Sender
	txHash := partDecMsg.TxHash
	h.mempool.AddPartialDecryption(txHash, sender, partDec)

	if int(h.mempool.GetThreshold(txHash)) < h.mempool.GetPartialDecryptionCount(txHash) {
		log.Printf("All partial decryptions received for: %s", txHash)
		encTx := h.mempool.GetIncludedTransaction(txHash)
		encryptedContent := encTx.Body.EncText

		pks := make([][]byte, CommitteSize)
		for i := 0; i < int(CommitteSize); i++ {
			pks[i] = h.GetSignerByIndex(uint64(i)).blsKey
		}

		partDecs := h.mempool.GetPartialDecryptions(txHash)

		// TODO: Seperate the committee size and the nodes count, for now we are using the same number
		content, err := h.crypto.DecryptTransaction(encryptedContent, pks, partDecs, encTx.Header.GammaG2, encTx.Body.Sa1, encTx.Body.Sa2, encTx.Body.Iv, encTx.Body.Threshold, CommitteSize+1)

		if err != nil {
			return fmt.Errorf("failed to decrypt transaction: %w", err)
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
		log.Printf("Transaction decrypted: %s", txHash)
		log.Printf("Our index: %d, Leader index: %d", h.ourIndex, leaderIndex)
		if h.ourIndex == leaderIndex-1 || leaderIndex == h.ourIndex {
			tx, err := sendRawTransaction(h.rpcUrl, string(content))
			if err != nil {
				return fmt.Errorf("failed to send transaction to blockchain: %w", err)
			}
			log.Printf("Transaction sent to the blockchain: %s", tx)
		}
	}
	return nil
}

func (h *Handler) handleEncryptedBatch(ctx context.Context, msg *message.Message, leaderIndex *uint64, ourIndex, CommitteSize uint64) error {

	log.Println("Handling encrypted batch")
	encBatchMsg := msg.Message.(*message.Message_EncryptedBatch).EncryptedBatch

	if encBatchMsg.Header.LeaderID == ourIndex {
		log.Println("Ignoring own batch")
		return nil
	}

	if encBatchMsg.Header.LeaderID != *leaderIndex {
		log.Printf("Ignoring batch from non-leader (received: %d, expected: %d)", encBatchMsg.Header.LeaderID, *leaderIndex)
		return nil
	}

	var txHashes []string
	for _, encTx := range encBatchMsg.Body.Transactions {
		txHashes = append(txHashes, string(encTx.Hash))
		if slices.Contains(encTx.PkIDs, ourIndex) {
			partDec, err := h.crypto.PartialDecrypt(encTx.GammaG2)
			if err != nil {
				return fmt.Errorf("failed to partially decrypt transaction: %w", err)
			}

			newMessage := &message.Message{
				Message: &message.Message_PartialDecryption{
					PartialDecryption: &message.PartialDecryption{
						TxHash:  encTx.Hash,
						Sender:  ourIndex,
						PartDec: partDec,
					},
				},
				MessageType: message.MessageType_PARTIAL_DECRYPTION,
			}

			msgBytes, err := proto.Marshal(newMessage)
			if err != nil {
				return fmt.Errorf("failed to marshal partial decryption message: %w", err)
			}
			err = h.topic.Publish(ctx, msgBytes)
			if err != nil {
				return fmt.Errorf("failed to publish partial decryption message: %w", err)
			}
		}
	}
	h.mempool.IncludeEncryptedTxs(txHashes)

	*leaderIndex++
	*leaderIndex %= CommitteSize
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
	h.mempool.AddOrderSig(orderSigMsg.BlockNum, *orderSig)
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

	log.Println("Random indexes: ", randomIndexes)

	pks := make([][]byte, h.committeeSize)
	ourSigners := *h.GetSigners()
	for i := 0; i < int(h.committeeSize); i++ {
		pks[i] = ourSigners[i].blsKey
	}

	encResponse, err := h.crypto.EncryptTransaction([]byte(tx), pks, h.threshold, h.committeeSize)
	if err != nil {
		log.Println("Error while encrypting the transaction: ", err)
		return err
	}

	txHash, err := h.CalculateTxHash(tx)
	if err != nil {
		log.Println("Error while calculating the transaction hash: ", err)
		return err
	}

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
		MessageType: message.MessageType_ENCRYPTED_TRANSACTION,
	}
	log.Println("Threshold number of message: ", h.threshold)
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

func (h *Handler) CheckTransactionDuplicate(txHash string) bool {
	encTx := h.mempool.GetEncryptedTransaction(txHash)
	decTx := h.mempool.GetDecryptedTransaction(txHash)
	incTx := h.mempool.GetIncludedTransaction(txHash)

	if encTx != nil || decTx != nil || incTx != nil {
		return true
	}

	return false
}

func (h *Handler) CheckTransactionDecrypted(hash string) bool {
	decTx := h.mempool.GetDecryptedTransaction(hash)
	return decTx != nil
}

func (h *Handler) CheckTransactionIncluded(hash string) bool {
	incTx := h.mempool.GetIncludedTransaction(hash)
	return incTx != nil
}

func (h *Handler) GetPartialDecryptionCount(hash string) int {
	return h.mempool.GetPartialDecryptionCount(hash)
}

func (h *Handler) GetMetadataOfTx(hash string) (committeeSize, threshold int) {
	encTx := h.mempool.GetEncryptedTransaction(hash)
	if encTx != nil {
		return len(encTx.Header.PkIDs), int(encTx.Body.Threshold)
	}

	incTx := h.mempool.GetIncludedTransaction(hash)
	if incTx != nil {
		return len(incTx.Header.PkIDs), int(incTx.Body.Threshold)
	}

	return 0, 0
}

func (h *Handler) CalculateTxHash(rawTx string) (string, error) {
	// Remove the '0x' prefix if present
	rawTx = strings.TrimPrefix(rawTx, "0x")

	// Decode the hex string to bytes
	txBytes, err := hex.DecodeString(rawTx)
	if err != nil {
		return "", fmt.Errorf("failed to decode hex string: %v", err)
	}

	// Create a new transaction object
	tx := new(ethTypes.Transaction)

	// Decode the RLP-encoded transaction
	err = rlp.DecodeBytes(txBytes, tx)
	if err != nil {
		return "", fmt.Errorf("failed to decode RLP: %v", err)
	}

	// Calculate the hash
	hash := tx.Hash()

	return hash.Hex(), nil
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
	body, err := io.ReadAll(resp.Body)
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
