package eip4844

import (
	"errors"
	"testing"
)

type staticKZGBatchVerifier struct {
	valid bool
}

func (v staticKZGBatchVerifier) VerifyBlobKZGProofBatch(blobs []Blob, commitments []KZGCommitment, proofs []KZGProof) bool {
	return v.valid
}

func TestValidateBlobGossipSuccess(t *testing.T) {
	var commitment KZGCommitment
	commitment[0] = 0xB1
	hash := KZGToVersionedHash(commitment)

	wrapper := BlobTxNetworkWrapper{
		BlobVersionedHashes: []VersionedHash{hash},
		Blobs:               make([]Blob, 1),
		Commitments:         []KZGCommitment{commitment},
		Proofs:              make([]KZGProof, 1),
	}

	if err := ValidateBlobGossip(wrapper, staticKZGBatchVerifier{valid: true}); err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestValidateBlobGossipRejectsMismatchedLengths(t *testing.T) {
	var commitment KZGCommitment
	hash := KZGToVersionedHash(commitment)

	wrapper := BlobTxNetworkWrapper{
		BlobVersionedHashes: []VersionedHash{hash},
		Blobs:               make([]Blob, 1),
		Commitments:         []KZGCommitment{commitment},
		Proofs:              nil,
	}

	err := ValidateBlobGossip(wrapper, staticKZGBatchVerifier{valid: true})
	if !errors.Is(err, ErrMismatchedBlobBundleLengths) {
		t.Fatalf("expected ErrMismatchedBlobBundleLengths, got: %v", err)
	}
}

func TestValidateBlobGossipRejectsCommitmentHashMismatch(t *testing.T) {
	var commitment KZGCommitment
	hash := KZGToVersionedHash(commitment)
	hash[2] ^= 0xFF

	wrapper := BlobTxNetworkWrapper{
		BlobVersionedHashes: []VersionedHash{hash},
		Blobs:               make([]Blob, 1),
		Commitments:         []KZGCommitment{commitment},
		Proofs:              make([]KZGProof, 1),
	}

	err := ValidateBlobGossip(wrapper, staticKZGBatchVerifier{valid: true})
	if !errors.Is(err, ErrCommitmentHashMismatch) {
		t.Fatalf("expected ErrCommitmentHashMismatch, got: %v", err)
	}
}

func TestValidateBlobGossipRejectsInvalidBatchProof(t *testing.T) {
	var commitment KZGCommitment
	hash := KZGToVersionedHash(commitment)

	wrapper := BlobTxNetworkWrapper{
		BlobVersionedHashes: []VersionedHash{hash},
		Blobs:               make([]Blob, 1),
		Commitments:         []KZGCommitment{commitment},
		Proofs:              make([]KZGProof, 1),
	}

	err := ValidateBlobGossip(wrapper, staticKZGBatchVerifier{valid: false})
	if !errors.Is(err, ErrInvalidBlobBatchProof) {
		t.Fatalf("expected ErrInvalidBlobBatchProof, got: %v", err)
	}
}
