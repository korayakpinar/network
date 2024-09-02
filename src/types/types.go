package types

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
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
			Hash:           hash,
			RawTransaction: nil,
			Status:         StatusPending,
			ReceivedAt:     time.Now(),
		}
	} else {
		decodedTx, err := DecodeRawTransaction(*rawTx)
		if err != nil {
			decodedTx = nil
		}

		return &Transaction{
			Hash:           hash,
			RawTransaction: decodedTx,
			Status:         StatusPending,
			ReceivedAt:     time.Now(),
		}
	}

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

	txBytes, err := hex.DecodeString(hexTx)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string: %w", err)
	}

	var ethTx types.Transaction
	err = rlp.DecodeBytes(txBytes, &ethTx)
	if err != nil {
		return nil, fmt.Errorf("failed to decode RLP: %w", err)
	}

	chainID := big.NewInt(17000) //Holesky
	signer := types.NewEIP155Signer(chainID)
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
}
