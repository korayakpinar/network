package types

import (
	"bytes"
	"encoding/binary"
)

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

	// EncTxs slice'ının uzunluğunu yaz
	binary.Write(buffer, binary.LittleEndian, uint32(len(b.EncTxs)))

	// Her bir EncryptedTxHeader'ı sırayla yaz
	for _, encTx := range b.EncTxs {
		// Hash'i yaz
		hashBytes := []byte(encTx.Hash)
		binary.Write(buffer, binary.LittleEndian, uint32(len(hashBytes)))
		buffer.Write(hashBytes)

		// GammaG2'yi yaz
		binary.Write(buffer, binary.LittleEndian, uint32(len(encTx.GammaG2)))
		buffer.Write(encTx.GammaG2)

		// PkIDs slice'ının uzunluğunu yaz
		binary.Write(buffer, binary.LittleEndian, uint32(len(encTx.PkIDs)))

		// PkIDs'leri yaz
		for _, pkID := range encTx.PkIDs {
			binary.Write(buffer, binary.LittleEndian, pkID)
		}
	}

	return buffer.Bytes()
}

type OrderSig struct {
	TxHeaders []*EncryptedTxHeader
	Signature string
}
