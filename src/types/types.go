package types

type PartialDecryption string
type TxHash string
type BlockNum uint32

// EncryptedTransaction represents the encrypted transaction
type EncryptedTransaction struct {
	Header *TransactionHeader
	Body   *TransactionBody
}

type TransactionHeader struct {
	Hash    TxHash
	GammaG2 string
}

// TransactionBody represents the body of the transaction
type TransactionBody struct {
	PkIDs     []uint32
	Sa1       []string
	Sa2       []string
	Iv        []byte
	Threshold uint32
}

type EncryptedBatch struct {
	Header *BatchHeader
	Body   *BatchBody
}

type BatchHeader struct {
	LeaderID  uint32
	BlockNum  BlockNum
	Hash      string
	Signature string
}

type BatchBody struct {
	EncTxs []*TransactionHeader
}

type OrderSig struct {
	TxHeaders []*TransactionHeader
	Signature string
}
