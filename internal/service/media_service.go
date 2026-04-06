package service

import (
	"mime/multipart"
	"nineshop-be/internal/utils"
)

type mediaService struct {
	uploadDir string
	publicURL string // Ví dụ: http://localhost:8080/uploads/
}

func NewMediaService(uploadDir, publicURL string) MediaService {
	return &mediaService{uploadDir: uploadDir, publicURL: publicURL}
}

func (s *mediaService) UploadMultiple(files []*multipart.FileHeader) ([]string, []error) {
	var urls []string
	var errs []error

	for _, file := range files {
		// Dùng hàm ValidateAndSaveFile bạn đã viết rất tốt trước đó
		filename, err := utils.ValidateAndSaveFile(file, s.uploadDir)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		urls = append(urls, s.publicURL+filename)
	}
	return urls, errs
}
