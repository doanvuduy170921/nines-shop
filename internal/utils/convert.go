package utils

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
)

func Float64ToPgTypeNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(fmt.Sprintf("%f", f))
	return n
}
