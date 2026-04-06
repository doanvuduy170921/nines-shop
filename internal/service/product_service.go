package service

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"nineshop-be/internal/db"
	"nineshop-be/internal/db/sqlc"
	"nineshop-be/internal/dto"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/utils"
)

type productService struct {
	repo repository.ProductRepository
	iu   ImagesUpdater
}

type ImagesUpdater interface {
	GetImagesByProductId(ctx context.Context, id int32) ([]string, error)
}

func NewProductService(repo repository.ProductRepository, iu ImagesUpdater) ProductService {
	return &productService{
		repo: repo,
		iu:   iu,
	}

}

func (ps *productService) GetAllProductByFilter(ctx context.Context) ([]sqlc.GetAllProductByFilterRow, error) {
	products, err := ps.repo.GetAllProductByFilter(ctx)
	if err != nil {
		return nil, utils.HandleDbError(err)
	}
	return products, nil
}

func (ps *productService) GetListVariantByPid(ctx context.Context, productID int64) ([]sqlc.GetListVariantByPidRow, error) {
	return ps.repo.GetListVariantByPid(ctx, productID)
}

func (ps *productService) AddProduct(ctx context.Context, input dto.AddProductRequestDto) (product sqlc.Product, err error) {
	// tạo transaction
	tx, err := db.DBPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return sqlc.Product{}, utils.WrapError(err, "Create Transaction fail", http.StatusInternalServerError)
	}
	qTx := db.DB.WithTx(tx)
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	product, err = qTx.AddProduct(ctx, sqlc.AddProductParams{
		Name:             input.Name,
		Slug:             utils.GenProductSlug(input.Name),
		BrandID:          utils.IntToPInt32(input.BrandID),
		CategoryID:       utils.IntToPInt32(input.CategoryID),
		Description:      &input.Description,
		ShortDescription: &input.ShortDescription,
		Status:           &input.Status,
		Thumbnail:        input.Thumbnail,
		HasVariants:      &input.HasVariant,
	})
	if err != nil {
		log.Printf("Error adding product: %v", err)
		return sqlc.Product{}, utils.HandleDbError(err)
	}

	for _, imgUrl := range input.Images {
		_, err = qTx.SaveAndUploadImg(ctx, sqlc.SaveAndUploadImgParams{
			ProductID: product.ID,
			ImageUrl:  imgUrl,
		})
		if err != nil {
			return sqlc.Product{}, utils.HandleDbError(err)
		}
	}

	// add spec product
	for _, spec := range input.Specifications {
		_, err = qTx.AddProductSpec(ctx, sqlc.AddProductSpecParams{
			ProductID:    product.ID,
			SpecKey:      spec.SpecKey,
			SpecValue:    spec.SpecValue,
			DisplayOrder: utils.IntToPInt32(spec.DisplayOrder),
		})
		if err != nil {
			return sqlc.Product{}, utils.HandleDbError(err)
		}
	}

	// add variant product
	for _, v := range input.Variants {
		attrBytes, err := json.Marshal(v.Attributes)
		if err != nil {
			return sqlc.Product{}, utils.WrapError(err, "Convert Attributes to JSON fail", http.StatusInternalServerError)
		}

		// Marshal images
		imageJSON, err := json.Marshal([]string{input.Thumbnail})
		if err != nil {
			return sqlc.Product{}, utils.WrapError(err, "Convert Image to JSON fail", http.StatusInternalServerError)
		}
		SKUVariant := utils.GenSKU(input.Name, v.Attributes)
		_, err = qTx.AddProductVariant(ctx, sqlc.AddProductVariantParams{
			ProductID:     product.ID,
			Sku:           &SKUVariant,
			Attributes:    attrBytes,
			Price:         utils.Float64ToPgTypeNumeric(v.Price),
			StockQuantity: v.StockQuantity,
			Images:        imageJSON,
			IsActive:      &v.IsActive,
		})
		if err != nil {
			log.Printf("Debug err:%v ", err)
			return sqlc.Product{}, utils.HandleDbError(err)
		}

	}
	// ✅ COMMIT TRANSACTION - QUAN TRỌNG!
	if err = tx.Commit(ctx); err != nil {
		return sqlc.Product{}, utils.WrapError(err, "Commit Transaction fail", http.StatusInternalServerError)
	}

	return product, nil
}

func (ps *productService) GetTop3Thumbnail(ctx context.Context) (dto.GetTop3TrendingRes, error) {

	images, err := ps.repo.GetTop3Trending(ctx, utils.IntToPInt32(8))
	if err != nil {
		return dto.GetTop3TrendingRes{}, utils.HandleDbError(err)
	}

	laptops, err := ps.repo.GetTop3Trending(ctx, utils.IntToPInt32(1))
	if err != nil {
		return dto.GetTop3TrendingRes{}, utils.HandleDbError(err)
	}
	keyboards, err := ps.repo.GetTop3Trending(ctx, utils.IntToPInt32(9))
	if err != nil {
		return dto.GetTop3TrendingRes{}, utils.HandleDbError(err)
	}
	screens, err := ps.repo.GetTop3Trending(ctx, utils.IntToPInt32(10))
	if err != nil {
		return dto.GetTop3TrendingRes{}, utils.HandleDbError(err)
	}

	mouses, err := ps.repo.GetTop3Trending(ctx, utils.IntToPInt32(7))
	res := dto.GetTop3TrendingRes{
		Images:    images,
		Mouses:    mouses,
		Laptops:   laptops,
		Screens:   screens,
		Keyboards: keyboards,
	}
	return res, nil
}

func (ps *productService) GetProductBySlug(ctx context.Context, slug string) (sqlc.GetProductBySlugRow, error) {
	return ps.repo.GetProductBySlug(ctx, slug)
}

func (ps *productService) GetListProducts(ctx context.Context, arg dto.GetProductsRequest) ([]sqlc.GetListProductsRow, error) {
	offset := (arg.Page - 1) * arg.Limit
	params := sqlc.GetListProductsParams{
		Limit:      arg.Limit,
		Offset:     offset,
		SearchName: &arg.SearchName,
		CateName:   &arg.CateName,
		MaxPrice:   utils.Float64ToPgTypeNumeric(arg.MaxPrice),
		MinPrice:   utils.Float64ToPgTypeNumeric(arg.MinPrice),
		SortDesc:   arg.SortBy,
	}

	return ps.repo.GetListProducts(ctx, params)
}
