package eip4844

import (
	"math/big"
	"testing"
)

func TestKZGToVersionedHashPrefix(t *testing.T) {
	var commitment KZGCommitment
	commitment[0] = 0xAA

	hash := KZGToVersionedHash(commitment)
	if hash[0] != VersionedHashVersionKZG {
		t.Fatalf("unexpected versioned hash prefix: got=0x%02x", hash[0])
	}
}

func TestFakeExponentialSimple(t *testing.T) {
	got := FakeExponential(1, 1, 1)
	if got.Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("unexpected fake exponential result: got=%s want=2", got.String())
	}
}

func TestIsCanonicalFieldElement(t *testing.T) {
	modulusBytes, err := BigIntToBytes32(BLSModulus())
	if err != nil {
		t.Fatalf("failed to convert modulus to bytes32: %v", err)
	}
	if IsCanonicalFieldElement(modulusBytes) {
		t.Fatal("modulus must be rejected as non-canonical")
	}

	lessThanModulus := new(big.Int).Sub(BLSModulus(), big.NewInt(1))
	lessThanModulusBytes, err := BigIntToBytes32(lessThanModulus)
	if err != nil {
		t.Fatalf("failed to convert value to bytes32: %v", err)
	}
	if !IsCanonicalFieldElement(lessThanModulusBytes) {
		t.Fatal("value below modulus must be canonical")
	}
}
