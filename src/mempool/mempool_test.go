package mempool_test

import (
	"bytes"
	"testing"

	"github.com/korayakpinar/network/src/mempool"
	"github.com/korayakpinar/network/src/types"
)

func TestMempool(t *testing.T) {
	// Test the mempool
	memp := mempool.NewMempool()
	if memp == nil {
		t.Errorf("Mempool is nil")
	}

	// Create dummy encrypted transaction
	encTx := &types.EncryptedTransaction{
		Header: &types.EncryptedTxHeader{
			Hash:    "hash",
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

	memp.AddEncryptedTx(encTx)

	memps := memp.GetTransactions()
	if len(memps) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(memps))
	}

	// Test the GetTransaction method
	tx := memp.GetTransaction("hash")
	if tx == nil {
		t.Errorf("Transaction not found")
	}
	beforeGammaG2 := tx.Header.GammaG2
	tx.Header.GammaG2 = []byte{4, 5, 6}
	afterGammaG2 := memp.GetTransaction("hash").Header.GammaG2
	if beforeGammaG2[0] == afterGammaG2[0] {
		t.Errorf("GammaG2 not updated")
	}

	th := memp.GetThreshold("hash")
	if th != encTx.Body.Threshold {
		t.Errorf("Threshold not zero")
	}

	memp.AddPartialDecryption("hash", &[]byte{1, 2, 3})
	preLength := len(*memp.GetPartialDecryptions("hash"))
	memp.AddPartialDecryption("hash", &[]byte{1, 2, 3})
	postLength := len(*memp.GetPartialDecryptions("hash"))
	if preLength != postLength {
		t.Errorf("Partial decryption added twice")
	}

}

func TestAddPartialDecryption(t *testing.T) {
	t.Run("Add new partial decryption", func(t *testing.T) {
		memp := mempool.NewMempool()
		hash := "hash1"
		partDec := []byte{1, 2, 3}

		memp.AddPartialDecryption(hash, &partDec)

		partDecs := memp.GetPartialDecryptions(hash)
		if partDecs == nil || len(*partDecs) != 1 {
			t.Errorf("Expected 1 partial decryption, got %d", len(*partDecs))
		}
	})

	t.Run("Add multiple unique partial decryptions", func(t *testing.T) {
		memp := mempool.NewMempool()
		hash := "hash1"
		partDec1 := []byte{1, 2, 3}
		partDec2 := []byte{4, 5, 6}

		memp.AddPartialDecryption(hash, &partDec1)
		memp.AddPartialDecryption(hash, &partDec2)

		partDecs := memp.GetPartialDecryptions(hash)
		if partDecs == nil || len(*partDecs) != 2 {
			t.Errorf("Expected 2 partial decryptions, got %d", len(*partDecs))
		}
	})

	t.Run("Attempt to add duplicate partial decryption", func(t *testing.T) {
		memp := mempool.NewMempool()
		hash := "hash1"
		partDec := []byte{1, 2, 3}

		memp.AddPartialDecryption(hash, &partDec)
		memp.AddPartialDecryption(hash, &partDec)

		partDecs := memp.GetPartialDecryptions(hash)
		if partDecs == nil || len(*partDecs) != 1 {
			t.Errorf("Expected 1 partial decryption after attempting to add duplicate, got %d", len(*partDecs))
		}

		// Verify the content
		if !bytes.Equal((*partDecs)[0], partDec) {
			t.Errorf("Partial decryption content does not match expected value")
		}
	})

	t.Run("Add partial decryptions for different hashes", func(t *testing.T) {
		memp := mempool.NewMempool()
		hash1 := "hash1"
		hash2 := "hash2"
		partDec1 := []byte{1, 2, 3}
		partDec2 := []byte{4, 5, 6}

		memp.AddPartialDecryption(hash1, &partDec1)
		memp.AddPartialDecryption(hash2, &partDec2)

		partDecs1 := memp.GetPartialDecryptions(hash1)
		partDecs2 := memp.GetPartialDecryptions(hash2)

		if partDecs1 == nil || len(*partDecs1) != 1 {
			t.Errorf("Expected 1 partial decryption for hash1, got %d", len(*partDecs1))
		}
		if partDecs2 == nil || len(*partDecs2) != 1 {
			t.Errorf("Expected 1 partial decryption for hash2, got %d", len(*partDecs2))
		}
	})

	t.Run("Check partial decryption count", func(t *testing.T) {
		memp := mempool.NewMempool()
		hash := "hash1"
		partDec1 := []byte{1, 2, 3}
		partDec2 := []byte{4, 5, 6}

		memp.AddPartialDecryption(hash, &partDec1)
		memp.AddPartialDecryption(hash, &partDec2)
		memp.AddPartialDecryption(hash, &partDec1) // Attempt to add duplicate

		count := memp.GetPartialDecryptionCount(hash)
		if count != 2 {
			t.Errorf("Expected partial decryption count of 2, got %d", count)
		}
	})
}
