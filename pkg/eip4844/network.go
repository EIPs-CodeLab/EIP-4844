package eip4844

import "fmt"

type KZGBatchVerifier interface {
	VerifyBlobKZGProofBatch(blobs []Blob, commitments []KZGCommitment, proofs []KZGProof) bool
}

func ValidateBlobGossip(wrapper BlobTxNetworkWrapper, verifier KZGBatchVerifier) error {
	if len(wrapper.BlobVersionedHashes) == 0 {
		return ErrMissingBlobHashes
	}

	n := len(wrapper.BlobVersionedHashes)
	if len(wrapper.Blobs) != n || len(wrapper.Commitments) != n || len(wrapper.Proofs) != n {
		return fmt.Errorf(
			"%w: hashes=%d blobs=%d commitments=%d proofs=%d",
			ErrMismatchedBlobBundleLengths,
			len(wrapper.BlobVersionedHashes),
			len(wrapper.Blobs),
			len(wrapper.Commitments),
			len(wrapper.Proofs),
		)
	}

	for i, commitment := range wrapper.Commitments {
		expected := KZGToVersionedHash(commitment)
		if expected != wrapper.BlobVersionedHashes[i] {
			return fmt.Errorf("%w: index=%d", ErrCommitmentHashMismatch, i)
		}
	}

	if verifier == nil {
		return ErrMissingKZGVerifier
	}

	if !verifier.VerifyBlobKZGProofBatch(wrapper.Blobs, wrapper.Commitments, wrapper.Proofs) {
		return ErrInvalidBlobBatchProof
	}

	return nil
}
