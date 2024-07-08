package mempool_test

import (
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

	memp.AddPartialDecryption("hash", &[]byte{1, 2, 3})
	preLength := len(*memp.GetPartialDecryptions("hash"))
	memp.AddPartialDecryption("hash", &[]byte{1, 2, 3})
	postLength := len(*memp.GetPartialDecryptions("hash"))
	if preLength != postLength {
		t.Errorf("Partial decryption added twice")
	}

}
