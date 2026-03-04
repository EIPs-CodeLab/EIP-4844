# EIP-4844 (Go Starter)

`EIP-4844: Shard Blob Transactions` adds a new transaction type (`0x03`) that carries blob data for cheaper rollup data availability.  
Blobs are not directly accessible by EVM execution, but their commitments (as versioned hashes) are.

Primary references:
- https://www.eip4844.com/
- https://eips.ethereum.org/EIPS/eip-4844

## What EIP-4844 Changes

EIP-4844 introduces:
- A new blob transaction format (`BLOB_TX_TYPE = 0x03`)
- Blob gas accounting with a separate base-fee market
- Header extensions: `blob_gas_used` and `excess_blob_gas`
- Blob-specific execution validity rules
- `BLOBHASH` opcode support (`0x49`)
- Point-evaluation precompile behavior for KZG proof verification
- Networking checks for blob side data (blobs/commitments/proofs)

This is the execution-layer foundation for proto-danksharding and a forward-compatible path to full sharding.

## About This Project

This repository is an EIPs CodeLab style implementation that translates EIP-4844 core logic into clean, testable Go code.

Current scope:
- Execution-layer validation logic
- Fee and gas math (`fake_exponential`, blob base fee, excess blob gas)
- Precompile input/verification flow (with pluggable KZG backend)
- Blob gossip wrapper consistency checks
- Example cases and unit tests

Out of scope:
- Full Ethereum client integration
- Full consensus layer implementation
- Real network propagation stack

## Project Structure

```text
.
├── cmd/
│   └── eip4844-demo/
│       └── main.go
├── pkg/
│   └── eip4844/
│       ├── constants.go
│       ├── doc.go
│       ├── errors.go
│       ├── fees.go
│       ├── helpers.go
│       ├── network.go
│       ├── precompile.go
│       ├── types.go
│       ├── validation.go
│       ├── helpers_test.go
│       ├── network_test.go
│       ├── precompile_test.go
│       └── validation_test.go
├── Makefile
└── go.mod
```

## Example Cases

Run example cases:

```bash
make examples
```

Expected output:

```text
valid block -> VALID
blob fee cap too low -> INVALID (...)
invalid versioned hash prefix -> INVALID (...)
```

Covered example/test scenarios include:
- Valid block with correct `blob_gas_used` accounting
- Blob tx rejected when `max_fee_per_blob_gas < base_fee_per_blob_gas`
- Blob tx rejected for wrong versioned hash prefix
- Precompile rejects bad length, non-canonical field values, commitment/hash mismatch
- Network wrapper rejects mismatched list lengths and invalid batch proofs

## How To Run

Prerequisite:
- Go 1.22+

Quick start:

```bash
make help
make test
make examples
```

Direct Go commands:

```bash
go test ./...
go run ./cmd/eip4844-demo
```

## Make Targets

- `make fmt`: Format Go code
- `make test`: Run all tests
- `make run`: Run demo app
- `make examples`: Run example cases
- `make check`: Run format + tests

## Security Checklist Before Production

Before production use, complete this checklist:

- [ ] Replace stub/mock KZG verifiers with an audited production KZG backend.
- [ ] Validate behavior against canonical fixtures from `ethereum/execution-spec-tests` (`eip4844_blobs`).
- [ ] Add fuzz/property tests for malformed transactions, wrapper payloads, and precompile inputs.
- [ ] Enforce strict overflow/underflow protections on all fee and gas arithmetic paths.
- [ ] Confirm canonical field-element checks for `z` and `y` in precompile paths (`< BLS_MODULUS`).
- [ ] Add mempool DoS controls (rate limits, size limits, replacement rules for blob txs).
- [ ] Gate rules by fork activation config to avoid pre-fork acceptance.
- [ ] Add observability: metrics for blob gas usage, base fee changes, and rejection reasons.
- [ ] Pin dependencies, run vulnerability scanning, and use reproducible CI builds.
- [ ] Run external security review/audit before integrating into a production node.
