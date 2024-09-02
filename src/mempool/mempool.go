package mempool

import (
	"sync"

	"github.com/korayakpinar/network/src/types"
)

type Mempool struct {
	Transactions    sync.Map // string (hash) -> *types.Transaction
	PartDecRegistry sync.Map // string (hash) -> map[uint64][]byte
}

func NewMempool() *Mempool {
	return &Mempool{
		Transactions:    sync.Map{},
		PartDecRegistry: sync.Map{},
	}
}

func (m *Mempool) AddTransaction(tx *types.Transaction) {
	m.Transactions.Store(tx.Hash, tx)
}

func (m *Mempool) GetTransaction(hash string) *types.Transaction {
	if tx, ok := m.Transactions.Load(hash); ok {
		return tx.(*types.Transaction)
	}
	return nil
}

func (m *Mempool) GetPendingTransactions() []*types.Transaction {
	var pendingTxs []*types.Transaction
	m.Transactions.Range(func(_, value interface{}) bool {
		tx := value.(*types.Transaction)
		if tx.Status == types.StatusPending {
			pendingTxs = append(pendingTxs, tx)
		}
		return true
	})
	return pendingTxs
}

func (m *Mempool) GetThreshold(hash string) uint64 {
	if tx, ok := m.Transactions.Load(hash); ok {
		return uint64(tx.(*types.Transaction).Threshold)
	}
	return 0
}

func (m *Mempool) SetTransactionProposed(hash string) {
	if tx, ok := m.Transactions.Load(hash); ok {
		tx.(*types.Transaction).SetProposed()
	}
}

func (m *Mempool) SetMultipleTransactionsProposed(hashes []string) {
	for _, hash := range hashes {
		m.SetTransactionProposed(hash)
	}
}

func (m *Mempool) SetTransactionDecrypted(hash string, decTx *types.DecryptedTransaction) {
	if tx, ok := m.Transactions.Load(hash); ok {
		tx.(*types.Transaction).SetDecrypted(decTx)
	}
}

func (m *Mempool) SetTransactionIncluded(hash string) {
	if tx, ok := m.Transactions.Load(hash); ok {
		tx.(*types.Transaction).SetIncluded()
	}
}

func (m *Mempool) SetMultipleTransactionsIncluded(hashes []string) {
	for _, hash := range hashes {
		m.SetTransactionIncluded(hash)
	}
}

func (m *Mempool) AddPartialDecryption(hash string, index uint64, partDec []byte) {
	tx := m.GetTransaction(hash)
	if tx == nil {
		return
	}

	partMap, _ := m.PartDecRegistry.LoadOrStore(hash, &sync.Map{})
	partMapPtr := partMap.(*sync.Map)
	if _, loaded := partMapPtr.LoadOrStore(index, partDec); !loaded {
		tx.IncrementPartialDecryptionCount()
	}
}

func (m *Mempool) GetPartialDecryptions(hash string) map[uint64][]byte {
	if partMap, ok := m.PartDecRegistry.Load(hash); ok {
		result := make(map[uint64][]byte)
		partMap.(*sync.Map).Range(func(key, value interface{}) bool {
			result[key.(uint64)] = value.([]byte)
			return true
		})
		return result
	}
	return nil
}

func (m *Mempool) GetPartialDecryptionCount(hash string) int {
	if tx, ok := m.Transactions.Load(hash); ok {
		return tx.(*types.Transaction).PartialDecryptionCount
	}
	return 0
}
