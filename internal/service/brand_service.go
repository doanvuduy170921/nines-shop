package service

import (
	"github.com/gin-gonic/gin"
	"nineshop-be/internal/db/sqlc"
	"nineshop-be/internal/repository"
)

type brandService struct {
	repo repository.BrandRepository
}

func NewBrandService(repo repository.BrandRepository) BrandService {
	return &brandService{
		repo: repo}

}

func (cs *brandService) GetAll(ctx *gin.Context) ([]sqlc.Brand, error) {
	context := ctx.Request.Context()
	categories, err := cs.repo.GetAll(context)
	if err != nil {
		return []sqlc.Brand{}, err
	}
	return categories, err
}
