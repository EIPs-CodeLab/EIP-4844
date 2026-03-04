package eip4844

import (
	"fmt"
	"math/big"
)

func ValidateBlock(block Block) error {
	expectedExcess := CalcExcessBlobGas(block.Parent)
	if block.Header.ExcessBlobGas != expectedExcess {
		return fmt.Errorf(
			"%w: got=%d want=%d",
			ErrInvalidExcessBlobGas,
			block.Header.ExcessBlobGas,
			expectedExcess,
		)
	}

	baseFeePerBlobGas := GetBaseFeePerBlobGas(block.Header)
	var blobGasUsed uint64

	for i, tx := range block.Transactions {
		maxTotalFee, err := maxTotalFee(tx)
		if err != nil {
			return fmt.Errorf("tx %d: %w", i, err)
		}

		if tx.SenderBalance == nil {
			return fmt.Errorf("tx %d: %w", i, ErrMissingSenderBalance)
		}
		if tx.SenderBalance.Cmp(maxTotalFee) < 0 {
			return fmt.Errorf("tx %d: %w", i, ErrInsufficientBalance)
		}

		if tx.Type != BlobTxType {
			continue
		}

		if err := validateBlobTransaction(tx, baseFeePerBlobGas); err != nil {
			return fmt.Errorf("tx %d: %w", i, err)
		}

		txBlobGas := GetTotalBlobGas(tx)
		if blobGasUsed > ^uint64(0)-txBlobGas {
			return fmt.Errorf("tx %d: %w", i, ErrBlobGasOverflow)
		}
		blobGasUsed += txBlobGas
	}

	if blobGasUsed > MaxBlobGasPerBlock {
		return fmt.Errorf(
			"%w: used=%d max=%d",
			ErrBlobGasAboveLimit,
			blobGasUsed,
			MaxBlobGasPerBlock,
		)
	}

	if block.Header.BlobGasUsed != blobGasUsed {
		return fmt.Errorf(
			"%w: got=%d want=%d",
			ErrBlobGasMismatch,
			block.Header.BlobGasUsed,
			blobGasUsed,
		)
	}

	return nil
}

func validateBlobTransaction(tx Transaction, baseFeePerBlobGas *big.Int) error {
	if len(tx.To) != 20 {
		return ErrInvalidToAddressLength
	}

	if len(tx.BlobVersionedHashes) == 0 {
		return ErrMissingBlobHashes
	}

	if tx.MaxFeePerBlobGas == nil {
		return ErrMissingMaxFeePerBlobGas
	}

	for i, h := range tx.BlobVersionedHashes {
		if h[0] != VersionedHashVersionKZG {
			return fmt.Errorf(
				"%w: index=%d found=0x%02x expected=0x%02x",
				ErrInvalidVersionedHashPrefix,
				i,
				h[0],
				VersionedHashVersionKZG,
			)
		}
	}

	if tx.MaxFeePerBlobGas.Cmp(baseFeePerBlobGas) < 0 {
		return fmt.Errorf(
			"%w: tx_cap=%s current_base_fee=%s",
			ErrBlobFeeCapTooLow,
			tx.MaxFeePerBlobGas.String(),
			baseFeePerBlobGas.String(),
		)
	}

	return nil
}

func maxTotalFee(tx Transaction) (*big.Int, error) {
	if tx.MaxFeePerGas == nil {
		return nil, ErrMissingMaxFeePerGas
	}
	if tx.MaxFeePerGas.Sign() < 0 {
		return nil, ErrNegativeFeeValue
	}

	gasFeeCap := new(big.Int).Mul(
		new(big.Int).SetUint64(tx.GasLimit),
		tx.MaxFeePerGas,
	)
	maxTotalFee := new(big.Int).Set(gasFeeCap)

	if tx.Type != BlobTxType {
		return maxTotalFee, nil
	}

	if tx.MaxFeePerBlobGas == nil {
		return nil, ErrMissingMaxFeePerBlobGas
	}
	if tx.MaxFeePerBlobGas.Sign() < 0 {
		return nil, ErrNegativeFeeValue
	}

	blobFeeCap := new(big.Int).Mul(
		new(big.Int).SetUint64(GetTotalBlobGas(tx)),
		tx.MaxFeePerBlobGas,
	)
	maxTotalFee.Add(maxTotalFee, blobFeeCap)

	return maxTotalFee, nil
}
