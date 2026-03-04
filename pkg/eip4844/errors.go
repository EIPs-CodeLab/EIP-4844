package eip4844

import "errors"

var (
	ErrInvalidBigIntValue           = errors.New("invalid big.Int value")
	ErrValueDoesNotFitIn32Bytes     = errors.New("value does not fit in 32 bytes")
	ErrInvalidExcessBlobGas         = errors.New("invalid excess_blob_gas in block header")
	ErrMissingSenderBalance         = errors.New("sender balance is required")
	ErrInsufficientBalance          = errors.New("insufficient sender balance for max_total_fee")
	ErrMissingMaxFeePerGas          = errors.New("max_fee_per_gas is required")
	ErrMissingMaxFeePerBlobGas      = errors.New("max_fee_per_blob_gas is required for blob transactions")
	ErrNegativeFeeValue             = errors.New("fee values must be non-negative")
	ErrMissingBlobHashes            = errors.New("blob transaction must include at least one blob versioned hash")
	ErrInvalidToAddressLength       = errors.New("blob transaction to field must be a 20-byte address")
	ErrInvalidVersionedHashPrefix   = errors.New("blob versioned hash must use VERSIONED_HASH_VERSION_KZG prefix")
	ErrBlobFeeCapTooLow             = errors.New("max_fee_per_blob_gas is below current blob base fee")
	ErrBlobGasOverflow              = errors.New("blob gas accumulator overflow")
	ErrBlobGasAboveLimit            = errors.New("blob gas used exceeds MAX_BLOB_GAS_PER_BLOCK")
	ErrBlobGasMismatch              = errors.New("header blob_gas_used does not match transaction-derived blob gas")
	ErrInvalidPrecompileInputLength = errors.New("point evaluation precompile input must be 192 bytes")
	ErrNonCanonicalFieldElement     = errors.New("point evaluation input includes non-canonical field element")
	ErrVersionedHashMismatch        = errors.New("commitment does not match provided versioned hash")
	ErrMissingKZGVerifier           = errors.New("kzg verifier is required")
	ErrInvalidKZGProof              = errors.New("invalid kzg proof")
	ErrMismatchedBlobBundleLengths  = errors.New("blob tx wrapper contains mismatched blob/commitment/proof lengths")
	ErrCommitmentHashMismatch       = errors.New("commitment hash does not match corresponding versioned hash")
	ErrInvalidBlobBatchProof        = errors.New("invalid blob batch proof")
)
