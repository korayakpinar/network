package mempool_test

import (
	"testing"
	"time"

	"github.com/korayakpinar/network/src/mempool"
	"github.com/korayakpinar/network/src/types"
	"github.com/stretchr/testify/assert"
)

func TestMempool(t *testing.T) {
	t.Run("NewMempool", func(t *testing.T) {
		mp := mempool.NewMempool()
		assert.NotNil(t, mp, "Mempool should not be nil")
	})

	t.Run("Transactions", func(t *testing.T) {
		mp := mempool.NewMempool()
		tx1 := createDummyTransaction("hash1")
		tx2 := createDummyTransaction("hash2")

		t.Run("AddTransaction", func(t *testing.T) {
			mp.AddTransaction(tx1)
			mp.AddTransaction(tx2)
			assert.NotNil(t, mp.GetTransaction("hash1"), "Should find transaction with hash1")
			assert.NotNil(t, mp.GetTransaction("hash2"), "Should find transaction with hash2")
		})

		t.Run("GetPendingTransactions", func(t *testing.T) {
			txs := mp.GetPendingTransactions()
			assert.Len(t, txs, 2, "Should return 2 pending transactions")
		})

		t.Run("SetTransactionProposed", func(t *testing.T) {
			mp.SetTransactionProposed("hash1")
			tx := mp.GetTransaction("hash1")
			assert.Equal(t, types.StatusProposed, tx.Status, "Transaction status should be Proposed")
			assert.False(t, tx.ProposedAt.IsZero(), "ProposedAt should be set")
		})

		t.Run("SetTransactionDecrypted", func(t *testing.T) {
			decryptedTx := &types.DecryptedTransaction{
				Header: &types.DecryptedTxHeader{Hash: "hash1"},
				Body:   &types.DecryptedTxBody{Content: "Decrypted content"},
			}
			mp.SetTransactionDecrypted("hash1", decryptedTx)
			tx := mp.GetTransaction("hash1")
			assert.Equal(t, types.StatusDecrypted, tx.Status, "Transaction status should be Decrypted")
			assert.False(t, tx.DecryptedAt.IsZero(), "DecryptedAt should be set")
			assert.Equal(t, "Decrypted content", tx.DecryptedTransaction.Body.Content)
		})

		t.Run("SetTransactionIncluded", func(t *testing.T) {
			mp.SetTransactionIncluded("hash1")
			tx := mp.GetTransaction("hash1")
			assert.Equal(t, types.StatusIncluded, tx.Status, "Transaction status should be Included")
			assert.False(t, tx.IncludedAt.IsZero(), "IncludedAt should be set")
		})
	})

	t.Run("PartialDecryptions", func(t *testing.T) {
		mp := mempool.NewMempool()
		tx := createDummyTransaction("hash1")
		mp.AddTransaction(tx)

		partDec1 := []byte{1, 2, 3}
		partDec2 := []byte{4, 5, 6}

		t.Run("AddPartialDecryption", func(t *testing.T) {
			mp.AddPartialDecryption("hash1", 1, partDec1)
			mp.AddPartialDecryption("hash1", 2, partDec2)
			mp.AddPartialDecryption("hash1", 1, partDec1) // Attempt to add duplicate

			tx := mp.GetTransaction("hash1")
			assert.Equal(t, 2, tx.PartialDecryptionCount, "Should have 2 unique partial decryptions")
		})

		t.Run("GetPartialDecryptions", func(t *testing.T) {
			partDecs := mp.GetPartialDecryptions("hash1")
			assert.NotNil(t, partDecs, "Partial decryptions should not be nil")
			assert.Len(t, partDecs, 2, "Should have 2 partial decryptions")
			assert.Equal(t, partDec1, partDecs[1], "Index 1 should contain partDec1")
			assert.Equal(t, partDec2, partDecs[2], "Index 2 should contain partDec2")
		})
	})

}

func createDummyTransaction(hash string) *types.Transaction {
	return &types.Transaction{
		Hash:           hash,
		RawTransaction: nil,
		Status:         types.StatusPending,
		ReceivedAt:     time.Now(),
		EncryptedTransaction: &types.EncryptedTransaction{
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
				Threshold: 2,
			},
		},
		CommitteeSize: 3,
		Threshold:     2,
	}
}
