package eip4844

import "fmt"

type KZGProofVerifier interface {
	VerifyKZGProof(commitment KZGCommitment, z [32]byte, y [32]byte, proof KZGProof) bool
}

// PointEvaluationPrecompile verifies the EIP-4844 point-evaluation input.
//
// Input format (192 bytes):
// versioned_hash (32) | z (32) | y (32) | commitment (48) | proof (48)
func PointEvaluationPrecompile(input []byte, verifier KZGProofVerifier) ([]byte, error) {
	if len(input) != 192 {
		return nil, fmt.Errorf(
			"%w: got=%d",
			ErrInvalidPrecompileInputLength,
			len(input),
		)
	}
	if verifier == nil {
		return nil, ErrMissingKZGVerifier
	}

	var versionedHash VersionedHash
	var z [32]byte
	var y [32]byte
	var commitment KZGCommitment
	var proof KZGProof

	copy(versionedHash[:], input[0:32])
	copy(z[:], input[32:64])
	copy(y[:], input[64:96])
	copy(commitment[:], input[96:144])
	copy(proof[:], input[144:192])

	if !IsCanonicalFieldElement(z) || !IsCanonicalFieldElement(y) {
		return nil, ErrNonCanonicalFieldElement
	}

	if KZGToVersionedHash(commitment) != versionedHash {
		return nil, ErrVersionedHashMismatch
	}

	if !verifier.VerifyKZGProof(commitment, z, y, proof) {
		return nil, ErrInvalidKZGProof
	}

	fieldElements := Uint64ToBytes32(FieldElementsPerBlob)
	modulus, err := BigIntToBytes32(blsModulus)
	if err != nil {
		return nil, err
	}

	out := make([]byte, 64)
	copy(out[0:32], fieldElements[:])
	copy(out[32:64], modulus[:])

	return out, nil
}
