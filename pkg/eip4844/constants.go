package eip4844

import "math/big"

const (
	BlobTxType byte = 0x03

	BytesPerFieldElement = 32
	FieldElementsPerBlob = 4096
	BlobSize             = BytesPerFieldElement * FieldElementsPerBlob

	VersionedHashVersionKZG byte = 0x01

	PointEvaluationPrecompileAddress byte   = 0x0A
	PointEvaluationPrecompileGas     uint64 = 50_000

	MaxBlobGasPerBlock    uint64 = 786_432
	TargetBlobGasPerBlock uint64 = 393_216

	MinBaseFeePerBlobGas      uint64 = 1
	BlobBaseFeeUpdateFraction uint64 = 3_338_477
	GasPerBlob                uint64 = 1 << 17

	HashOpcodeByte byte   = 0x49
	HashOpcodeGas  uint64 = 3

	MinEpochsForBlobSidecarsRequests uint64 = 4096
)

const blsModulusDecimal = "52435875175126190479447740508185965837690552500527637822603658699938581184513"

var blsModulus = mustBigInt(blsModulusDecimal)

// BLSModulus returns a defensive copy of the BLS modulus used by EIP-4844.
func BLSModulus() *big.Int {
	return new(big.Int).Set(blsModulus)
}

func mustBigInt(value string) *big.Int {
	bi, ok := new(big.Int).SetString(value, 10)
	if !ok {
		panic("invalid big integer literal: " + value)
	}
	return bi
}
