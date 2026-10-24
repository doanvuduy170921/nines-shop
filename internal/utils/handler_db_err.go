package utils

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"strings"
)

func HandleDbError(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			detail := pgErr.Detail
			return NewError(ErrorCodeBadRequest, formatUniqueError(detail))

		case "23503":
			return NewError(ErrorCodeBadRequest, "The relevant data does not exist or has been deleted")

		case "23502":
			return NewError(ErrorCodeBadRequest, fmt.Sprintf("The data field %s cannot be left empty", pgErr.ColumnName))

		case "22P02":
			return NewError(ErrorCodeBadRequest, "Invalid data format (JSON)")
		}
	}
	return WrapError(err, "System database error", ErrorCodeInternalError)
}

func formatUniqueError(detail string) string {
	if strings.Contains(detail, "sku") {
		return "This SKU already exists in the system"
	}
	if strings.Contains(detail, "slug") {
		return "This slug is already in use for another product"
	}
	if strings.Contains(detail, "email") {
		return "This email has already been registered"
	}
	return "Duplicate data:" + detail
}
