package service

import (
	"github.com/gin-gonic/gin"
	"nineshop-be/internal/db/sqlc"
	"nineshop-be/internal/repository"
)

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{
		repo: repo}

}

func (cs *categoryService) Create(ctx *gin.Context, name string) (sqlc.Category, error) {
	context := ctx.Request.Context()
	category, err := cs.repo.Create(context, name)
	if err != nil {
		return sqlc.Category{}, err
	}
	return category, err
}

func (cs *categoryService) GetAll(ctx *gin.Context) ([]sqlc.Category, error) {
	context := ctx.Request.Context()
	categories, err := cs.repo.GetAll(context)
	if err != nil {
		return []sqlc.Category{}, err
	}
	return categories, err
}
