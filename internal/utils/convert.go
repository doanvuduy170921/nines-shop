package utils

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
)

func Float64ToPgTypeNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(fmt.Sprintf("%f", f))
	return n
}

func PgNumericToInt64(n pgtype.Numeric) (int64, error) {
	// 1. Check NULL
	if !n.Valid {
		return 0, errors.New("numeric is NULL")
	}

	// 2. Không cho phép số thập phân
	if n.Exp != 0 {
		return 0, errors.New("numeric has decimal part, cannot convert to int64")
	}

	// 3. Convert *big.Int -> int64
	if !n.Int.IsInt64() {
		return 0, errors.New("numeric overflows int64")
	}

	return n.Int.Int64(), nil
}

func NumericToFloat64(n pgtype.Numeric) (float64, error) {
	if !n.Valid {
		return 0, errors.New("numeric is NULL or invalid")
	}

	v, err := n.Float64Value()
	if err != nil {
		return 0, err
	}

	return v.Float64, nil
}

func IntToPInt32(val int) *int32 {
	n := int32(val)
	return &n
}
