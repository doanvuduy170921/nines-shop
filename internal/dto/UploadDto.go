package dto

import "nineshop-be/internal/db/sqlc"

type UploadResult struct {
	SavedImages  []sqlc.ProductImage
	FailedFiles  []UploadError
	SuccessCount int
	FailedCount  int
}

type UploadError struct {
	Filename string
	Error    string
}
