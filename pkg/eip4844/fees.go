package eip4844

import (
	"math/big"
	"math/bits"
)

func CalcExcessBlobGas(parent Header) uint64 {
	sum, carry := bits.Add64(parent.ExcessBlobGas, parent.BlobGasUsed, 0)
	if carry != 0 {
		return ^uint64(0)
	}
	if sum < TargetBlobGasPerBlock {
		return 0
	}
	return sum - TargetBlobGasPerBlock
}

func GetTotalBlobGas(tx Transaction) uint64 {
	return GasPerBlob * uint64(len(tx.BlobVersionedHashes))
}

func GetBaseFeePerBlobGas(header Header) *big.Int {
	return FakeExponential(
		MinBaseFeePerBlobGas,
		header.ExcessBlobGas,
		BlobBaseFeeUpdateFraction,
	)
}

func CalcBlobFee(header Header, tx Transaction) *big.Int {
	totalBlobGas := new(big.Int).SetUint64(GetTotalBlobGas(tx))
	return new(big.Int).Mul(totalBlobGas, GetBaseFeePerBlobGas(header))
}
