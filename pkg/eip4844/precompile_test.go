package eip4844

import (
	"errors"
	"testing"
)

type staticKZGProofVerifier struct {
	valid bool
}

func (v staticKZGProofVerifier) VerifyKZGProof(commitment KZGCommitment, z [32]byte, y [32]byte, proof KZGProof) bool {
	return v.valid
}

func TestPointEvaluationPrecompileSuccess(t *testing.T) {
	var commitment KZGCommitment
	var proof KZGProof
	commitment[0] = 0xAB
	proof[0] = 0xCD

	versionedHash := KZGToVersionedHash(commitment)
	z := Uint64ToBytes32(5)
	y := Uint64ToBytes32(9)
	input := makePrecompileInput(versionedHash, z, y, commitment, proof)

	out, err := PointEvaluationPrecompile(input, staticKZGProofVerifier{valid: true})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(out) != 64 {
		t.Fatalf("unexpected output length: got=%d want=64", len(out))
	}

	expectedFieldElements := Uint64ToBytes32(FieldElementsPerBlob)
	if string(out[0:32]) != string(expectedFieldElements[:]) {
		t.Fatal("unexpected FIELD_ELEMENTS_PER_BLOB output")
	}

	expectedModulus, err := BigIntToBytes32(BLSModulus())
	if err != nil {
		t.Fatalf("failed to encode modulus: %v", err)
	}
	if string(out[32:64]) != string(expectedModulus[:]) {
		t.Fatal("unexpected BLS_MODULUS output")
	}
}

func TestPointEvaluationPrecompileRejectsVersionedHashMismatch(t *testing.T) {
	var commitment KZGCommitment
	var proof KZGProof
	commitment[0] = 0x01

	versionedHash := KZGToVersionedHash(commitment)
	versionedHash[1] ^= 0xFF

	z := Uint64ToBytes32(7)
	y := Uint64ToBytes32(11)
	input := makePrecompileInput(versionedHash, z, y, commitment, proof)

	_, err := PointEvaluationPrecompile(input, staticKZGProofVerifier{valid: true})
	if !errors.Is(err, ErrVersionedHashMismatch) {
		t.Fatalf("expected ErrVersionedHashMismatch, got: %v", err)
	}
}

func TestPointEvaluationPrecompileRejectsNonCanonicalFieldElement(t *testing.T) {
	var commitment KZGCommitment
	var proof KZGProof
	commitment[0] = 0x01

	versionedHash := KZGToVersionedHash(commitment)
	modulusBytes, err := BigIntToBytes32(BLSModulus())
	if err != nil {
		t.Fatalf("failed to encode modulus: %v", err)
	}

	// z == modulus is non-canonical by spec (must be strictly less than modulus).
	z := modulusBytes
	y := Uint64ToBytes32(1)
	input := makePrecompileInput(versionedHash, z, y, commitment, proof)

	_, err = PointEvaluationPrecompile(input, staticKZGProofVerifier{valid: true})
	if !errors.Is(err, ErrNonCanonicalFieldElement) {
		t.Fatalf("expected ErrNonCanonicalFieldElement, got: %v", err)
	}
}

func TestPointEvaluationPrecompileRequiresVerifier(t *testing.T) {
	var commitment KZGCommitment
	var proof KZGProof
	versionedHash := KZGToVersionedHash(commitment)
	z := Uint64ToBytes32(1)
	y := Uint64ToBytes32(2)
	input := makePrecompileInput(versionedHash, z, y, commitment, proof)

	_, err := PointEvaluationPrecompile(input, nil)
	if !errors.Is(err, ErrMissingKZGVerifier) {
		t.Fatalf("expected ErrMissingKZGVerifier, got: %v", err)
	}
}

func makePrecompileInput(
	versionedHash VersionedHash,
	z [32]byte,
	y [32]byte,
	commitment KZGCommitment,
	proof KZGProof,
) []byte {
	input := make([]byte, 192)
	copy(input[0:32], versionedHash[:])
	copy(input[32:64], z[:])
	copy(input[64:96], y[:])
	copy(input[96:144], commitment[:])
	copy(input[144:192], proof[:])
	return input
}
