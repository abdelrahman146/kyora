package inventory

import (
	"context"
	"fmt"
	"sort"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/abdelrahman146/kyora/internal/platform/types/atomic"
	"github.com/abdelrahman146/kyora/internal/platform/types/list"
	"github.com/shopspring/decimal"
)

type Service struct {
	storage         *Storage
	atomicProcessor atomic.AtomicProcessor
	bus             *bus.Bus
}

func NewService(storage *Storage, atomicProcessor atomic.AtomicProcessor, bus *bus.Bus) *Service {
	return &Service{
		storage:         storage,
		atomicProcessor: atomicProcessor,
		bus:             bus,
	}
}

func (s *Service) GetProductByID(ctx context.Context, actor *account.User, biz *business.Business, id string) (*Product, error) {
	return s.storage.products.FindOne(ctx,
		s.storage.products.ScopeBusinessID(biz.ID),
		s.storage.products.ScopeID(id),
	)
}

func (s *Service) GetVariantByID(ctx context.Context, actor *account.User, biz *business.Business, id string) (*Variant, error) {
	return s.storage.variants.FindOne(ctx,
		s.storage.variants.ScopeBusinessID(biz.ID),
		s.storage.variants.ScopeID(id),
	)
}

func (s *Service) GetCategoryByID(ctx context.Context, actor *account.User, biz *business.Business, id string) (*Category, error) {
	return s.storage.categories.FindOne(ctx,
		s.storage.categories.ScopeBusinessID(biz.ID),
		s.storage.categories.ScopeID(id),
	)
}

func (s *Service) ListProducts(ctx context.Context, actor *account.User, biz *business.Business, req *list.ListRequest) ([]*Product, error) {
	return s.storage.products.FindMany(ctx,
		s.storage.products.ScopeBusinessID(biz.ID),
		s.storage.products.WithPagination(req.Offset(), req.Limit()),
		s.storage.ScopeSearchTermByName(req.SearchTerm()),
		s.storage.products.WithOrderBy(req.ParsedOrderBy(ProductSchema)),
	)
}

func (s *Service) ListVariants(ctx context.Context, actor *account.User, biz *business.Business, req *list.ListRequest) ([]*Variant, error) {
	return s.storage.variants.FindMany(ctx,
		s.storage.variants.ScopeBusinessID(biz.ID),
		s.storage.variants.WithPagination(req.Offset(), req.Limit()),
		s.storage.ScopeSearchTermByName(req.SearchTerm()),
		s.storage.variants.WithOrderBy(req.ParsedOrderBy(VariantSchema)),
	)
}

func (s *Service) GetProductVariants(ctx context.Context, actor *account.User, biz *business.Business, productID string) ([]*Variant, error) {
	return s.storage.variants.FindMany(ctx,
		s.storage.variants.ScopeBusinessID(biz.ID),
		s.storage.variants.ScopeEquals(VariantSchema.ProductID, productID),
	)
}

func (s *Service) ListCategories(ctx context.Context, actor *account.User, biz *business.Business) ([]*Category, error) {
	return s.storage.categories.FindMany(ctx,
		s.storage.categories.ScopeBusinessID(biz.ID),
	)
}

func (s *Service) CountProducts(ctx context.Context, actor *account.User, biz *business.Business) (int64, error) {
	return s.storage.products.Count(ctx,
		s.storage.products.ScopeBusinessID(biz.ID),
	)
}

func (s *Service) CountVariants(ctx context.Context, actor *account.User, biz *business.Business) (int64, error) {
	return s.storage.variants.Count(ctx,
		s.storage.variants.ScopeBusinessID(biz.ID),
	)
}

func (s *Service) CountCategories(ctx context.Context, actor *account.User, biz *business.Business) (int64, error) {
	return s.storage.categories.Count(ctx,
		s.storage.categories.ScopeBusinessID(biz.ID),
	)
}

type CreateProductWithVariantsRequest struct {
	Product  CreateProductRequest   `json:"product" binding:"required"`
	Variants []CreateVariantRequest `json:"variants" binding:"required,dive,required"`
}

func (s *Service) CreateProductWithVariants(ctx context.Context, actor *account.User, biz *business.Business, req *CreateProductWithVariantsRequest) (*Product, error) {
	var product *Product
	err := s.atomicProcessor.Exec(ctx, func(txCtx context.Context) error {
		var err error
		product = &Product{
			BusinessID:  biz.ID,
			Name:        req.Product.Name,
			Description: req.Product.Description,
			CategoryID:  req.Product.CategoryID,
		}
		err = s.storage.products.CreateOne(txCtx, product)
		if err != nil {
			return err
		}
		variants := make([]*Variant, 0, len(req.Variants))
		for i, variantReq := range req.Variants {
			variants[i] = &Variant{
				BusinessID:         biz.ID,
				ProductID:          product.ID,
				Code:               variantReq.Code,
				Name:               fmt.Sprintf("%s - %s", product.Name, variantReq.Code),
				SKU:                variantReq.SKU,
				SalePrice:          variantReq.SalePrice,
				CostPrice:          variantReq.CostPrice,
				Currency:           biz.Currency,
				StockQuantity:      variantReq.StockQuantity,
				StockQuantityAlert: variantReq.StockQuantityAlert,
			}
			if variantReq.SKU == "" {
				variantReq.SKU = CreateProductSKU(biz.Descriptor, product.Name, variantReq.Code)
			}
		}
		err = s.storage.variants.CreateMany(txCtx, variants)
		if err != nil {
			return err
		}
		product.Variants = variants
		return nil
	})
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *Service) CreateCategory(ctx context.Context, actor *account.User, biz *business.Business, req *CreateCategoryRequest) (*Category, error) {
	category := &Category{
		BusinessID: biz.ID,
		Name:       req.Name,
	}
	err := s.storage.categories.CreateOne(ctx, category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *Service) CreateProduct(ctx context.Context, actor *account.User, biz *business.Business, req *CreateProductRequest) (*Product, error) {
	product := &Product{
		BusinessID:  biz.ID,
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryID,
	}
	err := s.storage.products.CreateOne(ctx, product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *Service) CreateVariant(ctx context.Context, actor *account.User, biz *business.Business, req *CreateVariantRequest) (*Variant, error) {
	variant := &Variant{
		BusinessID:         biz.ID,
		Code:               req.Code,
		SKU:                req.SKU,
		CostPrice:          req.CostPrice,
		SalePrice:          req.SalePrice,
		Currency:           biz.Currency,
		StockQuantity:      req.StockQuantity,
		StockQuantityAlert: req.StockQuantityAlert,
	}
	if req.SKU == "" {
		variant.SKU = CreateProductSKU(biz.Descriptor, "", req.Code)
	}
	err := s.storage.variants.CreateOne(ctx, variant)
	if err != nil {
		return nil, err
	}
	return variant, nil
}

func (s *Service) UpdateProduct(ctx context.Context, actor *account.User, biz *business.Business, product *Product, req *UpdateProductRequest) error {
	return s.atomicProcessor.Exec(ctx, func(tctx context.Context) error {
		if req.Name != "" {
			product.Name = req.Name
			variants, err := s.GetProductVariants(tctx, actor, biz, product.ID)
			if err != nil {
				return err
			}
			for _, variant := range variants {
				variant.Name = fmt.Sprintf("%s - %s", product.Name, variant.Code)
			}
			err = s.storage.variants.UpdateMany(tctx, variants)
			if err != nil {
				return err
			}
		}
		if req.Description != "" {
			product.Description = req.Description
		}
		if req.CategoryID != "" {
			product.CategoryID = req.CategoryID
		}
		return s.storage.products.UpdateOne(ctx, product)
	})
}

func (s *Service) UpdateVariant(ctx context.Context, actor *account.User, biz *business.Business, variantID string, req *UpdateVariantRequest) error {
	variant, err := s.storage.variants.FindByID(ctx, variantID)
	if err != nil {
		return err
	}
	if req.Code != "" {
		variant.Code = req.Code
		product, err := s.storage.products.FindByID(ctx, variant.ProductID)
		if err != nil {
			return err
		}
		variant.Name = fmt.Sprintf("%s - %s", product.Name, variant.Code)
	}
	if req.SKU != "" {
		variant.SKU = req.SKU
	}
	if !req.CostPrice.IsZero() {
		variant.CostPrice = req.CostPrice
	}
	if !req.SalePrice.IsZero() {
		variant.SalePrice = req.SalePrice
	}
	if req.Currency != "" {
		variant.Currency = req.Currency
	}
	if req.StockQuantity != 0 {
		variant.StockQuantity = req.StockQuantity
	}
	if req.StockQuantityAlert != 0 {
		variant.StockQuantityAlert = req.StockQuantityAlert
	}
	return s.storage.variants.UpdateOne(ctx, variant)
}

func (s *Service) UpdateCategory(ctx context.Context, actor *account.User, biz *business.Business, category *Category, req *UpdateCategoryRequest) error {
	if req.Name != "" {
		category.Name = req.Name
	}
	return s.storage.categories.UpdateOne(ctx, category)
}

func (s *Service) DeleteProduct(ctx context.Context, actor *account.User, biz *business.Business, id string) error {
	product, err := s.GetProductByID(ctx, actor, biz, id)
	if err != nil {
		return err
	}
	return s.storage.products.DeleteOne(ctx, product)
}

func (s *Service) DeleteVariant(ctx context.Context, actor *account.User, biz *business.Business, id string) error {
	variant, err := s.GetVariantByID(ctx, actor, biz, id)
	if err != nil {
		return err
	}
	return s.storage.variants.DeleteOne(ctx, variant)
}

func (s *Service) DeleteCategory(ctx context.Context, actor *account.User, biz *business.Business, id string) error {
	category, err := s.GetCategoryByID(ctx, actor, biz, id)
	if err != nil {
		return err
	}
	return s.storage.categories.DeleteOne(ctx, category)
}

func (s *Service) CountLowStockVariants(ctx context.Context, actor *account.User, biz *business.Business) (int64, error) {
	return s.storage.variants.Count(ctx,
		s.storage.variants.ScopeBusinessID(biz.ID),
		s.storage.ScopeLowStockVariants(),
	)
}

// CountOutOfStockVariants returns the number of variants with zero stock for the business.
func (s *Service) CountOutOfStockVariants(ctx context.Context, actor *account.User, biz *business.Business) (int64, error) {
	return s.storage.variants.Count(ctx,
		s.storage.variants.ScopeBusinessID(biz.ID),
		s.storage.variants.ScopeEquals(VariantSchema.StockQuantity, 0),
	)
}

// SumStockQuantity returns the total units in stock across all variants for the business.
func (s *Service) SumStockQuantity(ctx context.Context, actor *account.User, biz *business.Business) (int64, error) {
	sum, err := s.storage.variants.Sum(ctx, VariantSchema.StockQuantity,
		s.storage.variants.ScopeBusinessID(biz.ID),
	)
	if err != nil {
		return 0, err
	}
	return sum.IntPart(), nil
}

// SumInventoryValue returns the total inventory value (sum of cost_price * stock_quantity) for the business.
func (s *Service) SumInventoryValue(ctx context.Context, actor *account.User, biz *business.Business) (decimal.Decimal, error) {
	variants, err := s.storage.variants.FindMany(ctx,
		s.storage.variants.ScopeBusinessID(biz.ID),
	)
	if err != nil {
		return decimal.Zero, err
	}
	total := decimal.Zero
	for _, v := range variants {
		total = total.Add(v.CostPrice.Mul(decimal.NewFromInt(int64(v.StockQuantity))))
	}
	return total, nil
}

// ComputeTopProductsByInventoryValue returns the top N products by their on-hand inventory value (sum of variant cost * qty).
func (s *Service) ComputeTopProductsByInventoryValue(ctx context.Context, actor *account.User, biz *business.Business, limit int) ([]*Product, error) {
	if limit <= 0 {
		limit = 5
	}
	variants, err := s.storage.variants.FindMany(ctx,
		s.storage.variants.ScopeBusinessID(biz.ID),
		s.storage.variants.WithPreload(ProductStruct),
	)
	if err != nil {
		return nil, err
	}
	// aggregate by product id
	byProduct := map[string]decimal.Decimal{}
	productRef := map[string]*Product{}
	for _, v := range variants {
		if v.Product != nil {
			productRef[v.ProductID] = v.Product
		}
		byProduct[v.ProductID] = byProduct[v.ProductID].Add(v.CostPrice.Mul(decimal.NewFromInt(int64(v.StockQuantity))))
	}
	// rank product ids by value
	type kv struct {
		id    string
		value decimal.Decimal
	}
	arr := make([]kv, 0, len(byProduct))
	for id, val := range byProduct {
		arr = append(arr, kv{id: id, value: val})
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].value.GreaterThan(arr[j].value) })

	// pick top N and return products in that order
	n := min(len(arr), limit)
	out := make([]*Product, 0, n)
	for i := range n {
		if p, ok := productRef[arr[i].id]; ok {
			out = append(out, p)
		}
	}
	return out, nil
}
