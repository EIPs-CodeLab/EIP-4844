package eip4844

import "math/big"

type Blob [BlobSize]byte
type VersionedHash [32]byte
type KZGCommitment [48]byte
type KZGProof [48]byte

type Header struct {
	BlobGasUsed   uint64
	ExcessBlobGas uint64
}

type Transaction struct {
	Type          byte
	GasLimit      uint64
	MaxFeePerGas  *big.Int
	SenderBalance *big.Int

	// Blob-specific fields (required when Type == BlobTxType).
	To                  []byte
	MaxFeePerBlobGas    *big.Int
	BlobVersionedHashes []VersionedHash
}

type Block struct {
	Parent       Header
	Header       Header
	Transactions []Transaction
}

type BlobTxNetworkWrapper struct {
	BlobVersionedHashes []VersionedHash
	Blobs               []Blob
	Commitments         []KZGCommitment
	Proofs              []KZGProof
}
