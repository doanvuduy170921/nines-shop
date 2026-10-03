package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"math"
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

type ProductPaginationResponse struct {
	Data      []sqlc.GetAllProductByFilterRow `json:"Data"`
	Total     int64                           `json:"Total"`
	TotalPage int                             `json:"TotalPage"`
	Page      int                             `json:"Page"`
}

func (ps *productService) GetAllProductByFilter(ctx context.Context, search string, categoryID, brandID *int64, minPrice, maxPrice *float64, page, limit int) (ProductPaginationResponse, error) {
	// 1. Lấy dữ liệu từ Repo như bình thường
	var searchArg *string
	if search != "" {
		searchArg = &search
	}

	var minPriceArg, maxPriceArg pgtype.Numeric
	if minPrice != nil {
		minPriceArg.Scan(fmt.Sprintf("%f", *minPrice))
	}
	if maxPrice != nil {
		maxPriceArg.Scan(fmt.Sprintf("%f", *maxPrice))
	}

	products, err := ps.repo.GetAllProductByFilter(ctx, sqlc.GetAllProductByFilterParams{
		SearchName: searchArg,
		CategoryID: categoryID,
		BrandID:    brandID,
		MinPrice:   minPriceArg,
		MaxPrice:   maxPriceArg,
	})
	if err != nil {
		return ProductPaginationResponse{}, utils.HandleDbError(err)
	}

	// 2. Tính toán phân trang trên mảng kết quả
	total := int64(len(products))
	totalPage := int(math.Ceil(float64(total) / float64(limit)))
	if totalPage == 0 {
		totalPage = 1
	}

	start := (page - 1) * limit
	end := start + limit
	if start > int(total) {
		start = int(total)
	}
	if end > int(total) {
		end = int(total)
	}

	paginatedData := []sqlc.GetAllProductByFilterRow{}
	if total > 0 && start < int(total) {
		paginatedData = products[start:end]
	}

	// 3. Trả về đúng object có chứa trường Data, Total, TotalPage, Page
	return ProductPaginationResponse{
		Data:      paginatedData,
		Total:     total,
		TotalPage: totalPage,
		Page:      page,
	}, nil
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
	if err = tx.Commit(ctx); err != nil {
		return sqlc.Product{}, utils.WrapError(err, "Commit Transaction fail", http.StatusInternalServerError)
	}

	return product, nil
}

func (ps *productService) GetTop3Thumbnail(ctx context.Context) (dto.GetTop3TrendingRes, error) {

	images, err := ps.repo.GetTop3Trending(ctx, "Earphone")
	if err != nil {
		return dto.GetTop3TrendingRes{}, utils.HandleDbError(err)
	}

	laptops, err := ps.repo.GetTop3Trending(ctx, "Laptop")
	if err != nil {
		return dto.GetTop3TrendingRes{}, utils.HandleDbError(err)
	}
	keyboards, err := ps.repo.GetTop3Trending(ctx, "Keyboard")
	if err != nil {
		return dto.GetTop3TrendingRes{}, utils.HandleDbError(err)
	}
	screens, err := ps.repo.GetTop3Trending(ctx, "Screen")
	if err != nil {
		return dto.GetTop3TrendingRes{}, utils.HandleDbError(err)
	}

	mouses, err := ps.repo.GetTop3Trending(ctx, "Mouse")
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
