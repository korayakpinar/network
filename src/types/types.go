package types

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type TransactionStatus int

const (
	StatusPending TransactionStatus = iota
	StatusProposed
	StatusDecrypted
	StatusIncluded
)

type Transaction struct {
	Hash                   string
	RawTransaction         *RawTransaction
	EncryptedTransaction   *EncryptedTransaction
	DecryptedTransaction   *DecryptedTransaction
	Status                 TransactionStatus
	ReceivedAt             time.Time
	ProposedAt             time.Time
	DecryptedAt            time.Time
	IncludedAt             time.Time
	PartialDecryptionCount int
	CommitteeSize          int
	Threshold              int
	AlreadyEncrypted       bool
}

type RawTransaction struct {
	Nonce    uint64
	GasPrice *big.Int
	GasLimit uint64
	From     *common.Address
	To       *common.Address
	Value    *big.Int
	Data     []byte
	V, R, S  *big.Int // Signature values
}

func NewTransaction(rawTx *string, hash string) *Transaction {

	if rawTx == nil {
		return &Transaction{
			Hash:             hash,
			RawTransaction:   nil,
			Status:           StatusPending,
			ReceivedAt:       time.Now(),
			AlreadyEncrypted: false,
		}
	} else {
		decodedTx, err := DecodeRawTransaction(*rawTx)
		if err != nil {
			decodedTx = nil
		}

		return &Transaction{
			Hash:             hash,
			RawTransaction:   decodedTx,
			Status:           StatusPending,
			ReceivedAt:       time.Now(),
			AlreadyEncrypted: false,
		}
	}

}

func NewEncryptedTransaction(hash string, encTx *EncryptedTransaction) *Transaction {
	t := &Transaction{
		Hash:                 hash,
		RawTransaction:       nil,
		EncryptedTransaction: encTx,
		Status:               StatusPending,
		ReceivedAt:           time.Now(),
		CommitteeSize:        len(encTx.Header.PkIDs),
		Threshold:            int(encTx.Body.Threshold),
		AlreadyEncrypted:     true,
	}

	return t

}

func (t *Transaction) SetEncrypted(encTx *EncryptedTransaction) {
	t.EncryptedTransaction = encTx
	t.CommitteeSize = len(encTx.Header.PkIDs)
	t.Threshold = int(encTx.Body.Threshold)
}

func (t *Transaction) SetProposed() {
	t.Status = StatusProposed
	t.ProposedAt = time.Now()
}

func (t *Transaction) SetDecrypted(decTx *DecryptedTransaction) error {
	t.DecryptedTransaction = decTx
	t.Status = StatusDecrypted
	t.DecryptedAt = time.Now()

	rawTx, err := DecodeRawTransaction(decTx.Body.Content[2:])
	if err != nil {
		return fmt.Errorf("failed to decode raw transaction: %w", err)
	}
	t.RawTransaction = rawTx

	log.Println("Decoded raw transaction: ", t.RawTransaction)

	return nil
}

func (t *Transaction) SetIncluded() {
	t.Status = StatusIncluded
	t.IncludedAt = time.Now()

	if t.RawTransaction == nil {
		rawTx, err := DecodeRawTransaction(t.DecryptedTransaction.Body.Content[2:])
		if err != nil {
			log.Printf("Failed to decode raw transaction: %v", err)
		}
		t.RawTransaction = rawTx
	}
}

func (t *Transaction) IncrementPartialDecryptionCount() {
	t.PartialDecryptionCount++
}

// EncryptedTransaction represents the encrypted transaction
type EncryptedTransaction struct {
	Header *EncryptedTxHeader
	Body   *EncryptedTxBody
}

type EncryptedTxHeader struct {
	Hash    string
	GammaG2 []byte
	PkIDs   []uint64
}

// TransactionBody represents the body of the transaction
type EncryptedTxBody struct {
	Sa1       []byte
	Sa2       []byte
	Iv        []byte
	EncText   []byte
	Threshold uint64
}

type DecryptedTransaction struct {
	Header *DecryptedTxHeader
	Body   *DecryptedTxBody
}

type DecryptedTxHeader struct {
	Hash  string
	PkIDs []uint64
}

type DecryptedTxBody struct {
	Content string
}

type EncryptedBatch struct {
	Header *BatchHeader
	Body   *BatchBody
}

type BatchHeader struct {
	LeaderID  uint64
	BlockNum  uint64
	Hash      string
	Signature string
}

type BatchBody struct {
	EncTxs []*EncryptedTxHeader
}

func (b *BatchBody) Bytes() []byte {
	buffer := new(bytes.Buffer)

	binary.Write(buffer, binary.LittleEndian, uint32(len(b.EncTxs)))

	for _, encTx := range b.EncTxs {
		hashBytes := []byte(encTx.Hash)
		binary.Write(buffer, binary.LittleEndian, uint32(len(hashBytes)))
		buffer.Write(hashBytes)

		binary.Write(buffer, binary.LittleEndian, uint32(len(encTx.GammaG2)))
		buffer.Write(encTx.GammaG2)

		binary.Write(buffer, binary.LittleEndian, uint32(len(encTx.PkIDs)))

		for _, pkID := range encTx.PkIDs {
			binary.Write(buffer, binary.LittleEndian, pkID)
		}
	}

	return buffer.Bytes()
}

func DecodeRawTransaction(hexTx string) (*RawTransaction, error) {
	// Remove "0x" prefix if present
	hexTx = strings.TrimPrefix(hexTx, "0x")

	txBytes, err := hex.DecodeString(hexTx)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string: %w", err)
	}

	var ethTx types.Transaction
	err = ethTx.UnmarshalBinary(txBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}

	chainID := big.NewInt(17000) // Holesky
	signer := types.LatestSignerForChainID(chainID)
	from, err := types.Sender(signer, &ethTx)
	if err != nil {
		return nil, fmt.Errorf("failed to derive sender address: %w", err)
	}

	rawTx := &RawTransaction{
		Nonce:    ethTx.Nonce(),
		GasPrice: ethTx.GasPrice(),
		GasLimit: ethTx.Gas(),
		From:     &from,
		To:       ethTx.To(),
		Value:    ethTx.Value(),
		Data:     ethTx.Data(),
	}

	// Handle different transaction types
	switch ethTx.Type() {
	case types.LegacyTxType:
		rawTx.GasPrice = ethTx.GasPrice()
	case types.AccessListTxType:
		rawTx.GasPrice = ethTx.GasPrice()
		// You might want to handle AccessList here if needed
	case types.DynamicFeeTxType:
		rawTx.GasPrice = ethTx.GasFeeCap()
		// You might want to handle MaxPriorityFeePerGas and MaxFeePerGas here if needed
	}

	v, r, s := ethTx.RawSignatureValues()
	rawTx.V = v
	rawTx.R = r
	rawTx.S = s

	return rawTx, nil
}

type OrderSig struct {
	TxHeaders []*EncryptedTxHeader
	Signature string
}

type Signer struct {
	Index        uint64
	Address      common.Address
	BLSPublicKey []byte
	IsLocal      bool
	KeyPair      *KeyPair
}

func (s *Signer) GetAddress() common.Address {
	return s.Address
}

type Keystore struct {
	Keys map[uint64]KeyPair
}

type KeyPair struct {
	ECDSAPrivateKey string
	BLSPrivateKey   string
	BLSPublicKey    string
}

type Handler interface {
	BroadcastNewTransaction(tx *Transaction)
	BroadcastTransactionUpdate(tx *Transaction)
	HandleTransaction(tx string) error
	HandleEncryptedTransaction(*EncryptedTransaction) error
}
