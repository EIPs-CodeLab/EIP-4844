package main

import (
	"fmt"
	"math/big"

	"github.com/EIPs-CodeLab/EIP-4844/pkg/eip4844"
)

func main() {
	runCase("valid block", validBlockCase())
	runCase("blob fee cap too low", lowBlobFeeCase())
	runCase("invalid versioned hash prefix", invalidPrefixCase())
}

func runCase(name string, block eip4844.Block) {
	err := eip4844.ValidateBlock(block)
	if err != nil {
		fmt.Printf("%s -> INVALID (%v)\n", name, err)
		return
	}
	fmt.Printf("%s -> VALID\n", name)
}

func validBlockCase() eip4844.Block {
	parent := eip4844.Header{}
	header := eip4844.Header{
		ExcessBlobGas: eip4844.CalcExcessBlobGas(parent),
	}

	tx := makeBlobTx(2, big.NewInt(2))
	header.BlobGasUsed = eip4844.GetTotalBlobGas(tx)

	return eip4844.Block{
		Parent:       parent,
		Header:       header,
		Transactions: []eip4844.Transaction{tx},
	}
}

func lowBlobFeeCase() eip4844.Block {
	block := validBlockCase()
	block.Transactions[0].MaxFeePerBlobGas = big.NewInt(0)
	return block
}

func invalidPrefixCase() eip4844.Block {
	block := validBlockCase()
	block.Transactions[0].BlobVersionedHashes[0][0] = 0xFF
	return block
}

func makeBlobTx(blobCount int, maxFeePerBlobGas *big.Int) eip4844.Transaction {
	hashes := make([]eip4844.VersionedHash, blobCount)
	for i := range hashes {
		hashes[i] = makeVersionedHash(byte(i + 1))
	}

	return eip4844.Transaction{
		Type:                eip4844.BlobTxType,
		GasLimit:            21_000,
		MaxFeePerGas:        big.NewInt(100),
		SenderBalance:       new(big.Int).SetUint64(1_000_000_000_000),
		To:                  makeAddress(0x77),
		MaxFeePerBlobGas:    maxFeePerBlobGas,
		BlobVersionedHashes: hashes,
	}
}

func makeVersionedHash(seed byte) eip4844.VersionedHash {
	var h eip4844.VersionedHash
	h[0] = eip4844.VersionedHashVersionKZG
	for i := 1; i < len(h); i++ {
		h[i] = seed
	}
	return h
}

func makeAddress(seed byte) []byte {
	out := make([]byte, 20)
	for i := range out {
		out[i] = seed
	}
	return out
}
