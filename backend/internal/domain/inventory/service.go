package inventory

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/asset"
	"github.com/abdelrahman146/kyora/internal/platform/types/atomic"
	"github.com/abdelrahman146/kyora/internal/platform/types/list"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func normalizeCategoryDescriptor(v string) string {
	return strings.TrimSpace(strings.ToLower(v))
}

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

// ListProductsFilters contains optional filters for product listing.
type ListProductsFilters struct {
	CategoryID  string
	StockStatus StockStatus
}

func (s *Service) ListProducts(ctx context.Context, actor *account.User, biz *business.Business, req *list.ListRequest, filters *ListProductsFilters) ([]*Product, int64, error) {
	baseScopes := []func(db *gorm.DB) *gorm.DB{
		s.storage.products.ScopeWhere("products.business_id = ?", biz.ID),
	}

	needsVariantJoin := false
	useInnerJoin := false // INNER JOIN only if stock filter requires it

	// Check if sorting requires aggregation from variants
	needsVariantsCountSort := false
	needsCostPriceSort := false
	needsStockSort := false
	for _, ob := range req.OrderBy() {
		switch {
		case ob == "variantsCount" || ob == "-variantsCount":
			needsVariantsCountSort = true
			needsVariantJoin = true
		case ob == "costPrice" || ob == "-costPrice":
			needsCostPriceSort = true
			needsVariantJoin = true
		case ob == "stock" || ob == "-stock":
			needsStockSort = true
			needsVariantJoin = true
		}
	}

	if filters != nil {
		if filters.CategoryID != "" {
			baseScopes = append(baseScopes, s.storage.products.ScopeEquals(ProductSchema.CategoryID, filters.CategoryID))
		}

		if filters.StockStatus != "" {
			needsVariantJoin = true
			useInnerJoin = true // Stock filtering requires products to have variants
			switch filters.StockStatus {
			case StockStatusInStock:
				baseScopes = append(baseScopes,
					s.storage.variants.ScopeWhere("variants.stock_quantity > variants.stock_alert"),
				)
			case StockStatusLowStock:
				baseScopes = append(baseScopes,
					s.storage.variants.ScopeWhere("variants.stock_quantity > 0 AND variants.stock_quantity <= variants.stock_alert"),
				)
			case StockStatusOutOfStock:
				baseScopes = append(baseScopes,
					s.storage.variants.ScopeWhere("variants.stock_quantity = 0"),
				)
			}
		}
	}

	var listExtra []func(db *gorm.DB) *gorm.DB
	if req.SearchTerm() != "" {
		term := req.SearchTerm()
		like := "%" + term + "%"

		needsVariantJoin = true
		// Don't override useInnerJoin - keep LEFT JOIN for search-only scenarios
		baseScopes = append(baseScopes,
			s.storage.products.WithJoins("LEFT JOIN categories ON categories.id = products.category_id AND categories.deleted_at IS NULL"),
			s.storage.products.ScopeWhere(
				"(products.search_vector @@ websearch_to_tsquery('simple', ?) OR variants.search_vector @@ websearch_to_tsquery('simple', ?) OR categories.search_vector @@ websearch_to_tsquery('simple', ?) OR variants.sku ILIKE ?)",
				term,
				term,
				term,
				like,
			),
		)

		if !req.HasExplicitOrderBy() {
			rankExpr, err := database.WebSearchRankOrder(term, "products.search_vector", "variants.search_vector", "categories.search_vector")
			if err != nil {
				return nil, 0, err
			}
			listExtra = append(listExtra, s.storage.products.WithOrderByExpr(rankExpr))
		}
	}

	// If sorting by aggregated fields, add appropriate aggregations
	needsAggregation := needsVariantsCountSort || needsCostPriceSort || needsStockSort
	if needsAggregation {
		// Build aggregation SQL with all needed fields
		var aggFields []string
		if needsVariantsCountSort {
			aggFields = append(aggFields, "COUNT(*)::int as variants_count")
		}
		if needsCostPriceSort {
			aggFields = append(aggFields, "AVG(cost_price)::numeric as avg_cost_price")
		}
		if needsStockSort {
			aggFields = append(aggFields, "SUM(stock_quantity)::int as total_stock")
		}

		joinSQL := fmt.Sprintf(`
			LEFT JOIN LATERAL (
				SELECT %s
				FROM variants
				WHERE variants.product_id = products.id 
					AND variants.deleted_at IS NULL
			) AS product_agg ON true
		`, strings.Join(aggFields, ", "))

		baseScopes = append([]func(db *gorm.DB) *gorm.DB{
			s.storage.products.WithJoins(joinSQL),
		}, baseScopes...)

		// Parse custom ordering for aggregated fields
		customOrders := []string{}
		for _, ob := range req.OrderBy() {
			switch ob {
			case "variantsCount":
				customOrders = append(customOrders, "product_agg.variants_count ASC")
			case "-variantsCount":
				customOrders = append(customOrders, "product_agg.variants_count DESC")
			case "costPrice":
				customOrders = append(customOrders, "product_agg.avg_cost_price ASC")
			case "-costPrice":
				customOrders = append(customOrders, "product_agg.avg_cost_price DESC")
			case "stock":
				customOrders = append(customOrders, "product_agg.total_stock ASC")
			case "-stock":
				customOrders = append(customOrders, "product_agg.total_stock DESC")
			default:
				// Parse normal schema fields
				field, desc, found := list.ParseOrderField(ob, ProductSchema)
				if found {
					direction := "ASC"
					if desc {
						direction = "DESC"
					}
					customOrders = append(customOrders, field.Column()+" "+direction)
				}
			}
		}

		if len(customOrders) > 0 {
			listExtra = append(listExtra, s.storage.products.WithOrderBy(customOrders))
		}

		// If stock filtering is also needed, add regular JOIN after LATERAL
		if useInnerJoin {
			baseScopes = append(baseScopes,
				s.storage.products.WithJoins("INNER JOIN variants ON variants.product_id = products.id AND variants.deleted_at IS NULL"),
			)
		}
	} else if needsVariantJoin {
		joinType := "LEFT"
		if useInnerJoin {
			joinType = "INNER"
		}
		baseScopes = append([]func(db *gorm.DB) *gorm.DB{
			s.storage.products.WithJoins(fmt.Sprintf("%s JOIN variants ON variants.product_id = products.id AND variants.deleted_at IS NULL", joinType)),
		}, baseScopes...)
	}

	// Add GROUP BY if we have variant joins to prevent duplicates
	if needsVariantJoin || useInnerJoin {
		// Build GROUP BY clause including aggregated columns if needed
		groupByColumns := []string{"products.id"}

		// Add aggregated columns to GROUP BY when sorting by them
		if needsAggregation {
			if needsVariantsCountSort {
				groupByColumns = append(groupByColumns, "product_agg.variants_count")
			}
			if needsCostPriceSort {
				groupByColumns = append(groupByColumns, "product_agg.avg_cost_price")
			}
			if needsStockSort {
				groupByColumns = append(groupByColumns, "product_agg.total_stock")
			}
		}

		listExtra = append(listExtra, func(db *gorm.DB) *gorm.DB {
			return db.Group(strings.Join(groupByColumns, ", "))
		})
	}

	findOpts := append([]func(*gorm.DB) *gorm.DB{}, baseScopes...)
	findOpts = append(findOpts, listExtra...)
	findOpts = append(findOpts,
		s.storage.products.WithPagination(req.Offset(), req.Limit()),
	)

	// Only add parsed order by if we're not doing custom aggregation sort
	if !needsAggregation {
		findOpts = append(findOpts, s.storage.products.WithOrderBy(req.ParsedOrderBy(ProductSchema)))
	}

	findOpts = append(findOpts, s.storage.products.WithPreload("Variants"))

	items, err := s.storage.products.FindMany(ctx, findOpts...)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.storage.products.Count(ctx, baseScopes...)
	if err != nil {
		return nil, 0, err
	}

	return items, count, nil
}

func (s *Service) ListVariants(ctx context.Context, actor *account.User, biz *business.Business, req *list.ListRequest) ([]*Variant, error) {
	return s.storage.variants.FindMany(ctx,
		s.storage.variants.ScopeBusinessID(biz.ID),
		s.storage.variants.WithPagination(req.Offset(), req.Limit()),
		s.storage.variants.WithOrderBy(req.ParsedOrderBy(VariantSchema)),
	)
}

func (s *Service) GetProductVariants(ctx context.Context, actor *account.User, biz *business.Business, productID string) ([]*Variant, error) {
	return s.storage.variants.FindMany(ctx,
		s.storage.variants.ScopeBusinessID(biz.ID),
		s.storage.variants.ScopeEquals(VariantSchema.ProductID, productID),
	)
}

func (s *Service) ListProductVariants(ctx context.Context, actor *account.User, biz *business.Business, productID string, req *list.ListRequest) ([]*Variant, int64, error) {
	scopes := []func(db *gorm.DB) *gorm.DB{
		s.storage.variants.ScopeBusinessID(biz.ID),
		s.storage.variants.ScopeEquals(VariantSchema.ProductID, productID),
	}
	items, err := s.storage.variants.FindMany(ctx,
		append(scopes,
			s.storage.variants.WithPagination(req.Offset(), req.Limit()),
			s.storage.variants.WithOrderBy(req.ParsedOrderBy(VariantSchema)),
		)...,
	)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.storage.variants.Count(ctx, scopes...)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
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
	Product  CreateProductRequest          `json:"product" binding:"required"`
	Variants []CreateProductVariantRequest `json:"variants" binding:"required,min=1,dive"`
}

// CreateProductVariantRequest is used when creating a product with variants in a single request.
// It intentionally does not include ProductID because the product is created atomically.
type CreateProductVariantRequest struct {
	Code               string                 `json:"code" binding:"required"`
	SKU                string                 `json:"sku" binding:"omitempty"`
	Photos             []asset.AssetReference `json:"photos" binding:"omitempty,max=10,dive"`
	CostPrice          *decimal.Decimal       `json:"costPrice" binding:"required"`
	SalePrice          *decimal.Decimal       `json:"salePrice" binding:"required"`
	StockQuantity      *int                   `json:"stockQuantity" binding:"required,gte=0"`
	StockQuantityAlert *int                   `json:"stockQuantityAlert" binding:"required,gte=0"`
}

func (s *Service) CreateProductWithVariants(ctx context.Context, actor *account.User, biz *business.Business, req *CreateProductWithVariantsRequest) (*Product, error) {
	var product *Product
	err := s.atomicProcessor.Exec(ctx, func(txCtx context.Context) error {
		var err error
		// Validate category ownership to prevent cross-tenant reference.
		if _, err := s.GetCategoryByID(txCtx, actor, biz, req.Product.CategoryID); err != nil {
			return err
		}

		photos := AssetReferenceList(req.Product.Photos)
		product = &Product{
			BusinessID:  biz.ID,
			Name:        req.Product.Name,
			Description: req.Product.Description,
			Photos:      photos,
			CategoryID:  req.Product.CategoryID,
		}
		err = s.storage.products.CreateOne(txCtx, product)
		if err != nil {
			return err
		}
		variants := make([]*Variant, len(req.Variants))
		for i, variantReq := range req.Variants {
			if variantReq.CostPrice == nil {
				return problem.BadRequest("costPrice is required").With("field", "costPrice")
			}
			if variantReq.SalePrice == nil {
				return problem.BadRequest("salePrice is required").With("field", "salePrice")
			}
			if variantReq.StockQuantity == nil {
				return problem.BadRequest("stockQuantity is required").With("field", "stockQuantity")
			}
			if variantReq.StockQuantityAlert == nil {
				return problem.BadRequest("stockQuantityAlert is required").With("field", "stockQuantityAlert")
			}
			if variantReq.CostPrice.IsNegative() {
				return problem.BadRequest("costPrice must be >= 0").With("field", "costPrice")
			}
			if variantReq.SalePrice.IsNegative() {
				return problem.BadRequest("salePrice must be >= 0").With("field", "salePrice")
			}
			photos := AssetReferenceList(variantReq.Photos)
			sku := strings.TrimSpace(variantReq.SKU)
			if sku == "" {
				sku = CreateProductSKU(biz.Descriptor, product.Name, variantReq.Code)
			}
			variants[i] = &Variant{
				BusinessID:         biz.ID,
				ProductID:          product.ID,
				Code:               variantReq.Code,
				Name:               fmt.Sprintf("%s - %s", product.Name, variantReq.Code),
				SKU:                sku,
				SalePrice:          *variantReq.SalePrice,
				CostPrice:          *variantReq.CostPrice,
				Currency:           biz.Currency,
				Photos:             photos,
				StockQuantity:      *variantReq.StockQuantity,
				StockQuantityAlert: *variantReq.StockQuantityAlert,
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
	descriptor := normalizeCategoryDescriptor(req.Descriptor)
	if descriptor == "" {
		return nil, problem.BadRequest("descriptor is required").With("field", "descriptor")
	}
	category := &Category{
		BusinessID: biz.ID,
		Name:       req.Name,
		Descriptor: descriptor,
	}
	err := s.storage.categories.CreateOne(ctx, category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *Service) CreateProduct(ctx context.Context, actor *account.User, biz *business.Business, req *CreateProductRequest) (*Product, error) {
	// Validate category ownership to prevent cross-tenant reference.
	if _, err := s.GetCategoryByID(ctx, actor, biz, req.CategoryID); err != nil {
		return nil, err
	}

	photos := AssetReferenceList(req.Photos)
	product := &Product{
		BusinessID:  biz.ID,
		Name:        req.Name,
		Description: req.Description,
		Photos:      photos,
		CategoryID:  req.CategoryID,
	}
	err := s.storage.products.CreateOne(ctx, product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *Service) CreateVariant(ctx context.Context, actor *account.User, biz *business.Business, req *CreateVariantRequest) (*Variant, error) {
	product, err := s.GetProductByID(ctx, actor, biz, req.ProductID)
	if err != nil {
		return nil, err
	}
	if req.CostPrice == nil {
		return nil, problem.BadRequest("costPrice is required").With("field", "costPrice")
	}
	if req.SalePrice == nil {
		return nil, problem.BadRequest("salePrice is required").With("field", "salePrice")
	}
	if req.StockQuantity == nil {
		return nil, problem.BadRequest("stockQuantity is required").With("field", "stockQuantity")
	}
	if req.StockQuantityAlert == nil {
		return nil, problem.BadRequest("stockQuantityAlert is required").With("field", "stockQuantityAlert")
	}
	if req.CostPrice.IsNegative() {
		return nil, problem.BadRequest("costPrice must be >= 0").With("field", "costPrice")
	}
	if req.SalePrice.IsNegative() {
		return nil, problem.BadRequest("salePrice must be >= 0").With("field", "salePrice")
	}
	sku := strings.TrimSpace(req.SKU)
	if sku == "" {
		sku = CreateProductSKU(biz.Descriptor, product.Name, req.Code)
	}
	photos := AssetReferenceList(req.Photos)
	variant := &Variant{
		BusinessID:         biz.ID,
		ProductID:          product.ID,
		Code:               req.Code,
		Name:               fmt.Sprintf("%s - %s", product.Name, req.Code),
		SKU:                sku,
		CostPrice:          *req.CostPrice,
		SalePrice:          *req.SalePrice,
		Currency:           biz.Currency,
		Photos:             photos,
		StockQuantity:      *req.StockQuantity,
		StockQuantityAlert: *req.StockQuantityAlert,
	}
	err = s.storage.variants.CreateOne(ctx, variant)
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
		if req.Photos != nil {
			product.Photos = AssetReferenceList(req.Photos)
		}
		if req.CategoryID != "" {
			if _, err := s.GetCategoryByID(tctx, actor, biz, req.CategoryID); err != nil {
				return err
			}
			product.CategoryID = req.CategoryID
		}
		return s.storage.products.UpdateOne(tctx, product)
	})
}

func (s *Service) UpdateVariant(ctx context.Context, actor *account.User, biz *business.Business, variantID string, req *UpdateVariantRequest) error {
	variant, err := s.GetVariantByID(ctx, actor, biz, variantID)
	if err != nil {
		return err
	}
	if req.Code != nil {
		variant.Code = strings.TrimSpace(*req.Code)
		product, err := s.GetProductByID(ctx, actor, biz, variant.ProductID)
		if err != nil {
			return err
		}
		variant.Name = fmt.Sprintf("%s - %s", product.Name, variant.Code)
	}
	if req.SKU != nil {
		variant.SKU = strings.TrimSpace(*req.SKU)
	}
	if req.Photos != nil {
		variant.Photos = AssetReferenceList(req.Photos)
	}
	if req.CostPrice != nil {
		if req.CostPrice.IsNegative() {
			return problem.BadRequest("costPrice must be >= 0").With("field", "costPrice")
		}
		variant.CostPrice = *req.CostPrice
	}
	if req.SalePrice != nil {
		if req.SalePrice.IsNegative() {
			return problem.BadRequest("salePrice must be >= 0").With("field", "salePrice")
		}
		variant.SalePrice = *req.SalePrice
	}
	if req.Currency != nil {
		variant.Currency = strings.TrimSpace(strings.ToUpper(*req.Currency))
	}
	if req.StockQuantity != nil {
		variant.StockQuantity = *req.StockQuantity
	}
	if req.StockQuantityAlert != nil {
		variant.StockQuantityAlert = *req.StockQuantityAlert
	}
	return s.storage.variants.UpdateOne(ctx, variant)
}

func (s *Service) UpdateCategory(ctx context.Context, actor *account.User, biz *business.Business, category *Category, req *UpdateCategoryRequest) error {
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Descriptor != "" {
		category.Descriptor = normalizeCategoryDescriptor(req.Descriptor)
	}
	return s.storage.categories.UpdateOne(ctx, category)
}

func (s *Service) DeleteProduct(ctx context.Context, actor *account.User, biz *business.Business, id string) error {
	return s.atomicProcessor.Exec(ctx, func(tctx context.Context) error {
		product, err := s.GetProductByID(tctx, actor, biz, id)
		if err != nil {
			return err
		}
		if err := s.storage.variants.DeleteMany(tctx,
			s.storage.variants.ScopeBusinessID(biz.ID),
			s.storage.variants.ScopeEquals(VariantSchema.ProductID, product.ID),
		); err != nil {
			return err
		}
		return s.storage.products.DeleteOne(tctx, product)
	})
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

type TopProductByInventoryValue struct {
	Product        *Product        `json:"product"`
	InventoryValue decimal.Decimal `json:"inventoryValue"`
}

// ComputeTopProductsByInventoryValueDetailed returns the top N products by their on-hand inventory value
// and includes the computed inventory value per product.
func (s *Service) ComputeTopProductsByInventoryValueDetailed(ctx context.Context, actor *account.User, biz *business.Business, limit int) ([]*TopProductByInventoryValue, error) {
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
	byProduct := map[string]decimal.Decimal{}
	productRef := map[string]*Product{}
	for _, v := range variants {
		if v.Product != nil {
			productRef[v.ProductID] = v.Product
		}
		byProduct[v.ProductID] = byProduct[v.ProductID].Add(v.CostPrice.Mul(decimal.NewFromInt(int64(v.StockQuantity))))
	}
	type kv struct {
		id    string
		value decimal.Decimal
	}
	arr := make([]kv, 0, len(byProduct))
	for id, val := range byProduct {
		arr = append(arr, kv{id: id, value: val})
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].value.GreaterThan(arr[j].value) })

	n := min(len(arr), limit)
	out := make([]*TopProductByInventoryValue, 0, n)
	for i := range n {
		p, ok := productRef[arr[i].id]
		if !ok {
			continue
		}
		out = append(out, &TopProductByInventoryValue{Product: p, InventoryValue: arr[i].value})
	}
	return out, nil
}

// ComputeTopProductsByInventoryValue returns the top N products by their on-hand inventory value (sum of variant cost * qty).
func (s *Service) ComputeTopProductsByInventoryValue(ctx context.Context, actor *account.User, biz *business.Business, limit int) ([]*Product, error) {
	detailed, err := s.ComputeTopProductsByInventoryValueDetailed(ctx, actor, biz, limit)
	if err != nil {
		return nil, err
	}
	out := make([]*Product, 0, len(detailed))
	for _, d := range detailed {
		out = append(out, d.Product)
	}
	return out, nil
}
