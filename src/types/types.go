package types

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

type OrderSig struct {
	TxHeaders []*EncryptedTxHeader
	Signature string
}
