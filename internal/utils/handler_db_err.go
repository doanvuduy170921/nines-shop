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

	// Kiểm tra xem lỗi có phải từ Postgres không
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // Unique Violation (Trùng lặp dữ liệu UNIQUE)
			detail := pgErr.Detail // Ví dụ: "Key (sku)=(ROG-G15) already exists."
			return NewError(ErrorCodeBadRequest, formatUniqueError(detail))

		case "23503": // Foreign Key Violation (Lỗi khóa ngoại)
			return NewError(ErrorCodeBadRequest, "Dữ liệu liên quan (Brand/Category) không tồn tại hoặc đã bị xóa.")

		case "23502": // Not Null Violation
			return NewError(ErrorCodeBadRequest, fmt.Sprintf("Trường dữ liệu %s không được để trống.", pgErr.ColumnName))

		case "22P02": // Invalid Text Representation (Lỗi JSON syntax bạn vừa gặp)
			return NewError(ErrorCodeBadRequest, "Định dạng dữ liệu (JSON) không hợp lệ.")
		}
	}

	// Các lỗi không xác định khác trả về 500
	return WrapError(err, "Lỗi cơ sở dữ liệu hệ thống", ErrorCodeInternalError)
}

// Hàm bổ trợ để lấy tên cột bị trùng từ Detail của Postgres
func formatUniqueError(detail string) string {
	// Detail format: Key (sku)=(ROG-G15) already exists.
	if strings.Contains(detail, "sku") {
		return "Mã SKU này đã tồn tại trong hệ thống."
	}
	if strings.Contains(detail, "slug") {
		return "Đường dẫn (Slug) này đã được sử dụng cho sản phẩm khác."
	}
	if strings.Contains(detail, "email") {
		return "Email này đã được đăng ký."
	}
	return "Dữ liệu bị trùng lặp: " + detail
}
