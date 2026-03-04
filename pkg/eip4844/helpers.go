package eip4844

import (
	"crypto/sha256"
	"math/big"
)

// KZGToVersionedHash implements:
// VERSIONED_HASH_VERSION_KZG + sha256(commitment)[1:].
func KZGToVersionedHash(commitment KZGCommitment) VersionedHash {
	digest := sha256.Sum256(commitment[:])
	var out VersionedHash
	out[0] = VersionedHashVersionKZG
	copy(out[1:], digest[1:])
	return out
}

// FakeExponential approximates factor * e**(numerator/denominator) using
// integer math and a Taylor series expansion.
func FakeExponential(factor, numerator, denominator uint64) *big.Int {
	if denominator == 0 {
		panic("denominator must be non-zero")
	}

	i := uint64(1)
	output := new(big.Int)

	numeratorAccum := new(big.Int).Mul(
		new(big.Int).SetUint64(factor),
		new(big.Int).SetUint64(denominator),
	)

	numeratorBI := new(big.Int).SetUint64(numerator)
	denominatorBI := new(big.Int).SetUint64(denominator)

	for numeratorAccum.Sign() > 0 {
		output.Add(output, numeratorAccum)

		term := new(big.Int).Mul(numeratorAccum, numeratorBI)
		divisor := new(big.Int).Mul(denominatorBI, new(big.Int).SetUint64(i))
		numeratorAccum = term.Div(term, divisor)
		i++
	}

	return output.Div(output, denominatorBI)
}

func IsCanonicalFieldElement(input [32]byte) bool {
	value := new(big.Int).SetBytes(input[:])
	return value.Cmp(blsModulus) < 0
}

func Uint64ToBytes32(value uint64) [32]byte {
	var out [32]byte
	b := new(big.Int).SetUint64(value).Bytes()
	copy(out[32-len(b):], b)
	return out
}

func BigIntToBytes32(value *big.Int) ([32]byte, error) {
	var out [32]byte
	if value == nil || value.Sign() < 0 {
		return out, ErrInvalidBigIntValue
	}

	b := value.Bytes()
	if len(b) > 32 {
		return out, ErrValueDoesNotFitIn32Bytes
	}

	copy(out[32-len(b):], b)
	return out, nil
}
