package eip4844

import (
	"errors"
	"math/big"
	"testing"
)

func TestValidateBlockSuccess(t *testing.T) {
	parent := Header{}
	header := Header{
		ExcessBlobGas: CalcExcessBlobGas(parent),
	}

	tx := sampleBlobTx(2)
	header.BlobGasUsed = GetTotalBlobGas(tx)

	block := Block{
		Parent:       parent,
		Header:       header,
		Transactions: []Transaction{tx},
	}

	if err := ValidateBlock(block); err != nil {
		t.Fatalf("expected block to be valid, got error: %v", err)
	}
}

func TestValidateBlockRejectsLowBlobFeeCap(t *testing.T) {
	parent := Header{}
	header := Header{
		ExcessBlobGas: CalcExcessBlobGas(parent),
	}

	tx := sampleBlobTx(1)
	tx.MaxFeePerBlobGas = big.NewInt(0)
	header.BlobGasUsed = GetTotalBlobGas(tx)

	block := Block{
		Parent:       parent,
		Header:       header,
		Transactions: []Transaction{tx},
	}

	err := ValidateBlock(block)
	if !errors.Is(err, ErrBlobFeeCapTooLow) {
		t.Fatalf("expected ErrBlobFeeCapTooLow, got: %v", err)
	}
}

func TestValidateBlockRejectsInvalidVersionedHashPrefix(t *testing.T) {
	parent := Header{}
	header := Header{
		ExcessBlobGas: CalcExcessBlobGas(parent),
	}

	tx := sampleBlobTx(1)
	tx.BlobVersionedHashes[0][0] = 0x02
	header.BlobGasUsed = GetTotalBlobGas(tx)

	block := Block{
		Parent:       parent,
		Header:       header,
		Transactions: []Transaction{tx},
	}

	err := ValidateBlock(block)
	if !errors.Is(err, ErrInvalidVersionedHashPrefix) {
		t.Fatalf("expected ErrInvalidVersionedHashPrefix, got: %v", err)
	}
}

func TestValidateBlockRejectsBlobGasAboveLimit(t *testing.T) {
	parent := Header{}
	header := Header{
		ExcessBlobGas: CalcExcessBlobGas(parent),
	}

	tx := sampleBlobTx(7) // 7 * GAS_PER_BLOB > MAX_BLOB_GAS_PER_BLOCK
	header.BlobGasUsed = GetTotalBlobGas(tx)

	block := Block{
		Parent:       parent,
		Header:       header,
		Transactions: []Transaction{tx},
	}

	err := ValidateBlock(block)
	if !errors.Is(err, ErrBlobGasAboveLimit) {
		t.Fatalf("expected ErrBlobGasAboveLimit, got: %v", err)
	}
}

func TestValidateBlockRejectsWrongHeaderExcessBlobGas(t *testing.T) {
	parent := Header{
		BlobGasUsed:   GasPerBlob,
		ExcessBlobGas: 0,
	}
	header := Header{
		BlobGasUsed:   0,
		ExcessBlobGas: 1, // expected 0 because parent usage is below target
	}

	block := Block{
		Parent: parent,
		Header: header,
	}

	err := ValidateBlock(block)
	if !errors.Is(err, ErrInvalidExcessBlobGas) {
		t.Fatalf("expected ErrInvalidExcessBlobGas, got: %v", err)
	}
}

func sampleBlobTx(blobCount int) Transaction {
	hashes := make([]VersionedHash, blobCount)
	for i := range hashes {
		hashes[i] = sampleVersionedHash(byte(i + 1))
	}

	return Transaction{
		Type:                BlobTxType,
		GasLimit:            21_000,
		MaxFeePerGas:        big.NewInt(100),
		SenderBalance:       new(big.Int).SetUint64(1_000_000_000_000),
		To:                  sampleAddress(0x11),
		MaxFeePerBlobGas:    big.NewInt(2),
		BlobVersionedHashes: hashes,
	}
}

func sampleVersionedHash(seed byte) VersionedHash {
	var h VersionedHash
	h[0] = VersionedHashVersionKZG
	for i := 1; i < len(h); i++ {
		h[i] = seed
	}
	return h
}

func sampleAddress(seed byte) []byte {
	out := make([]byte, 20)
	for i := range out {
		out[i] = seed
	}
	return out
}
