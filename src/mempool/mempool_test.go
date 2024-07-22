package mempool_test

import (
	"bytes"
	"testing"

	"github.com/korayakpinar/network/src/mempool"
	"github.com/korayakpinar/network/src/types"
	"github.com/stretchr/testify/assert"
)

func TestMempool(t *testing.T) {
	t.Run("NewMempool", func(t *testing.T) {
		mp := mempool.NewMempool()

		assert.NotNil(t, mp, "Mempool should not be nil")
	})

	t.Run("EncryptedTransactions", func(t *testing.T) {
		mp := mempool.NewMempool()
		tx1 := createDummyEncryptedTx("hash1")
		tx2 := createDummyEncryptedTx("hash2")

		t.Run("AddEncryptedTx", func(t *testing.T) {
			mp.AddEncryptedTx(tx1)
			mp.AddEncryptedTx(tx2)
			assert.Equal(t, 2, mp.GetEncryptedTransactionCount(), "Should have 2 encrypted transactions")
		})

		t.Run("GetEncryptedTransactions", func(t *testing.T) {
			txs := mp.GetEncryptedTransactions()
			assert.Len(t, txs, 2, "Should return 2 encrypted transactions")
		})

		t.Run("GetEncryptedTransaction", func(t *testing.T) {
			tx := mp.GetEncryptedTransaction("hash1")
			assert.NotNil(t, tx, "Should find transaction with hash1")
			assert.Equal(t, "hash1", tx.Header.Hash, "Retrieved transaction should have correct hash")
		})

		t.Run("RemoveEncryptedTransactions", func(t *testing.T) {
			mp.RemoveEncryptedTransactions([]string{"hash1"})
			assert.Equal(t, 1, mp.GetEncryptedTransactionCount(), "Should have 1 encrypted transaction after removal")
			assert.Nil(t, mp.GetEncryptedTransaction("hash1"), "Removed transaction should not be found")
		})

		t.Run("GetThreshold", func(t *testing.T) {
			threshold := mp.GetThreshold("hash2")
			assert.Equal(t, uint64(1), threshold, "Threshold should match the one set in dummy transaction")
		})
	})

	t.Run("DecryptedTransactions", func(t *testing.T) {
		mp := mempool.NewMempool()
		tx1 := createDummyDecryptedTx("hash1")
		tx2 := createDummyDecryptedTx("hash2")

		t.Run("AddDecryptedTx", func(t *testing.T) {
			mp.AddDecryptedTx(tx1)
			mp.AddDecryptedTx(tx2)
			assert.Equal(t, 2, mp.GetDecryptedTransactionCount(), "Should have 2 decrypted transactions")
		})

		t.Run("GetDecryptedTransactions", func(t *testing.T) {
			txs := mp.GetDecryptedTransactions()
			assert.Len(t, txs, 2, "Should return 2 decrypted transactions")
		})

		t.Run("GetDecryptedTransaction", func(t *testing.T) {
			tx := mp.GetDecryptedTransaction("hash1")
			assert.NotNil(t, tx, "Should find transaction with hash1")
			assert.Equal(t, "hash1", tx.Header.Hash, "Retrieved transaction should have correct hash")
		})

		t.Run("RemoveDecryptedTransactions", func(t *testing.T) {
			mp.RemoveDecryptedTransactions([]string{"hash1"})
			assert.Equal(t, 1, mp.GetDecryptedTransactionCount(), "Should have 1 decrypted transaction after removal")
			assert.Nil(t, mp.GetDecryptedTransaction("hash1"), "Removed transaction should not be found")
		})
	})

	t.Run("PartialDecryptions", func(t *testing.T) {
		mp := mempool.NewMempool()
		hash := "hash1"
		partDec1 := []byte{1, 2, 3}
		partDec2 := []byte{4, 5, 6}

		t.Run("AddPartialDecryption", func(t *testing.T) {
			mp.AddPartialDecryption(hash, &partDec1)
			mp.AddPartialDecryption(hash, &partDec2)
			mp.AddPartialDecryption(hash, &partDec1) // Attempt to add duplicate

			count := mp.GetPartialDecryptionCount(hash)
			assert.Equal(t, uint64(2), count, "Should have 2 unique partial decryptions")
		})

		t.Run("GetPartialDecryptions", func(t *testing.T) {
			partDecs := mp.GetPartialDecryptions(hash)
			assert.NotNil(t, partDecs, "Partial decryptions should not be nil")
			assert.Len(t, *partDecs, 2, "Should have 2 partial decryptions")
			assert.True(t, bytes.Equal((*partDecs)[0], partDec1) || bytes.Equal((*partDecs)[1], partDec1), "Should contain partDec1")
			assert.True(t, bytes.Equal((*partDecs)[0], partDec2) || bytes.Equal((*partDecs)[1], partDec2), "Should contain partDec2")
		})
	})

	t.Run("OrderSignatures", func(t *testing.T) {
		mp := mempool.NewMempool()
		blockNum := uint64(1)
		orderSig1 := types.OrderSig{Signature: string([]byte{1, 2, 3})}
		orderSig2 := types.OrderSig{Signature: string([]byte{4, 5, 6})}

		t.Run("AddOrderSig", func(t *testing.T) {
			mp.AddOrderSig(blockNum, orderSig1)
			mp.AddOrderSig(blockNum, orderSig2)
		})

		t.Run("GetOrderSigs", func(t *testing.T) {
			orderSigs := mp.GetOrderSigs(blockNum)
			assert.NotNil(t, orderSigs, "Order signatures should not be nil")
			assert.Len(t, *orderSigs, 2, "Should have 2 order signatures")
			assert.Equal(t, orderSig1, (*orderSigs)[0], "First order signature should match")
			assert.Equal(t, orderSig2, (*orderSigs)[1], "Second order signature should match")
		})
	})
}

func createDummyEncryptedTx(hash string) *types.EncryptedTransaction {
	return &types.EncryptedTransaction{
		Header: &types.EncryptedTxHeader{
			Hash:    hash,
			GammaG2: []byte{1, 2, 3},
			PkIDs:   []uint64{1, 2, 3},
		},
		Body: &types.EncryptedTxBody{
			Sa1:       []byte{1, 2, 3},
			Sa2:       []byte{1, 2, 3},
			Iv:        []byte{1, 2, 3},
			EncText:   []byte{1, 2, 3},
			Threshold: 1,
		},
	}
}

func createDummyDecryptedTx(hash string) *types.DecryptedTransaction {
	return &types.DecryptedTransaction{
		Header: &types.DecryptedTxHeader{
			Hash:  hash,
			PkIDs: []uint64{1, 2, 3},
		},
		Body: &types.DecryptedTxBody{
			Content: "Decrypted content",
		},
	}
}
