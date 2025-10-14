package inventory

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/types"
	"github.com/shopspring/decimal"
)

type InventoryService struct {
	products      *productRepository
	variants      *variantRepository
	store         *store.StoreService
	atomicProcess *db.AtomicProcess
}

func NewInventoryService(products *productRepository, variants *variantRepository, store *store.StoreService, atomicProcess *db.AtomicProcess) *InventoryService {
	return &InventoryService{
		products:      products,
		variants:      variants,
		store:         store,
		atomicProcess: atomicProcess,
	}
}

func (s *InventoryService) CreateProduct(ctx context.Context, storeID string, productReq *CreateProductRequest, variantsReq []*CreateVariantRequest) (*Product, error) {
	var product *Product
	err := s.atomicProcess.Exec(ctx, func(ctx context.Context) error {
		// Load store to access standardized store Code for SKU generation
		st, err := s.store.GetStoreByID(ctx, storeID)
		if err != nil {
			return err
		}
		product = &Product{
			StoreID:     storeID,
			Name:        productReq.Name,
			Description: productReq.Description,
			Tags:        productReq.Tags,
		}
		if err := s.products.createOne(ctx, product); err != nil {
			return err
		}
		variants := s.toVariants(storeID, st.Code, product.ID, productReq.Name, variantsReq)
		if len(variants) > 0 {
			if err := s.createVariantsWithRetry(ctx, st.Code, productReq.Name, variants); err != nil {
				return err
			}
		}
		product.Variants = variants
		return nil
	})
	return product, err
}

func (s *InventoryService) AddVariantToProduct(ctx context.Context, productID string, variantReq *CreateVariantRequest) (*Variant, error) {
	var variant *Variant
	err := s.atomicProcess.Exec(ctx, func(ctx context.Context) error {
		// Ensure StoreID via product
		prod, err := s.products.findByID(ctx, productID)
		if err != nil {
			return err
		}
		st, err := s.store.GetStoreByID(ctx, prod.StoreID)
		if err != nil {
			return err
		}
		variant = s.toVariant(prod.StoreID, st.Code, productID, prod.Name, variantReq)
		if err := s.createVariantWithRetry(ctx, st.Code, prod.Name, variant); err != nil {
			return err
		}
		return nil
	})
	return variant, err
}

// toVariants converts create requests to Variant models, generating SKUs when missing.
func (s *InventoryService) toVariants(storeID, storeCode, productID, productName string, reqs []*CreateVariantRequest) []*Variant {
	variants := make([]*Variant, 0, len(reqs))
	for _, v := range reqs {
		variant := s.toVariant(storeID, storeCode, productID, productName, v)
		variants = append(variants, variant)
	}
	return variants
}

func (s *InventoryService) toVariant(storeID, storeCode, productID, productName string, v *CreateVariantRequest) *Variant {
	sku := strings.TrimSpace(v.SKU)
	if sku == "" {
		sku = GenerateSku(storeCode, productName, v.Code)
	}
	return &Variant{
		Name:          fmt.Sprintf("%s - %s", productName, v.Code),
		Code:          v.Code,
		SKU:           sku,
		StoreID:       storeID,
		ProductID:     productID,
		CostPrice:     v.CostPrice,
		SalePrice:     v.SalePrice,
		StockQuantity: v.StockQuantity,
		StockAlert:    v.StockAlert,
	}
}

// createVariantsWithRetry inserts variants and retries on unique SKU conflicts by regenerating SKUs.
func (s *InventoryService) createVariantsWithRetry(ctx context.Context, storeCode, productName string, variants []*Variant) error {
	const maxAttempts = 5
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := s.variants.createManyStrict(ctx, variants); err != nil {
			if db.IsUniqueViolation(err) {
				for _, variant := range variants {
					variant.SKU = GenerateSku(storeCode, productName, variant.Name)
				}
				if attempt == maxAttempts {
					return err
				}
				continue
			}
			return err
		}
		return nil
	}
	return errors.New("max attempts reached")
}

// createVariantWithRetry inserts a single variant and retries on unique SKU conflicts by regenerating SKU.
func (s *InventoryService) createVariantWithRetry(ctx context.Context, storeCode, productName string, variant *Variant) error {
	const maxAttempts = 5
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := s.variants.createOne(ctx, variant); err != nil {
			if db.IsUniqueViolation(err) {
				variant.SKU = GenerateSku(storeCode, productName, variant.Name)
				if attempt == maxAttempts {
					return err
				}
				continue
			}
			return err
		}
		return nil
	}
	return errors.New("max attempts reached")
}

func (s *InventoryService) GetProductByID(ctx context.Context, productID string) (*Product, error) {
	product, err := s.products.findByID(ctx, productID, db.WithPreload(VariantStruct))
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *InventoryService) GetVariantByID(ctx context.Context, variantID string) (*Variant, error) {
	variant, err := s.variants.findByID(ctx, variantID, db.WithPreload(ProductStruct))
	if err != nil {
		return nil, err
	}
	return variant, nil
}

func (s *InventoryService) GetVariantBySKU(ctx context.Context, sku string) (*Variant, error) {
	variant, err := s.variants.findOne(ctx, s.variants.scopeSKU(sku), db.WithPreload(ProductStruct))
	if err != nil {
		return nil, err
	}
	return variant, nil
}

func (s *InventoryService) ListVariantsByProductID(ctx context.Context, productID string) ([]*Variant, error) {
	variants, err := s.variants.list(ctx, s.variants.scopeProductID(productID))
	if err != nil {
		return nil, err
	}
	return variants, nil
}

func (s *InventoryService) ListVariants(ctx context.Context, filter *VariantFilter, page int, pageSize int, orderBy string, ascending bool) ([]*Variant, error) {
	variants, err := s.variants.list(ctx, s.variants.scopeFilter(filter), db.WithPagination(page, pageSize), db.WithPreload(ProductStruct), db.WithSorting(orderBy, ascending))
	if err != nil {
		return nil, err
	}
	return variants, nil
}

func (s *InventoryService) CountVariants(ctx context.Context, filter *VariantFilter) (int64, error) {
	count, err := s.variants.count(ctx, s.variants.scopeFilter(filter))
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *InventoryService) ListProducts(ctx context.Context, filter *ProductFilter, page int, pageSize int, orderBy string, ascending bool) ([]*Product, error) {
	products, err := s.products.list(ctx, s.products.scopeFilter(filter), db.WithPagination(page, pageSize), db.WithPreload(VariantStruct), db.WithSorting(orderBy, ascending))
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *InventoryService) CountProducts(ctx context.Context, filter *ProductFilter) (int64, error) {
	count, err := s.products.count(ctx, s.products.scopeFilter(filter))
	if err != nil {
		return 0, err
	}
	return count, nil
}

type UpdateVariant struct {
	ID     string
	Update *UpdateVariantRequest
}

func (s *InventoryService) UpdateProduct(ctx context.Context, productID string, updateReq *UpdateProductRequest, updateVariantsReq []*UpdateVariant) (*Product, error) {
	var product *Product
	err := s.atomicProcess.Exec(ctx, func(ctx context.Context) error {
		var err error
		product, err = s.products.findByID(ctx, productID)
		if err != nil {
			return err
		}
		updateProductFields(product, updateReq)
		if err := s.products.updateOne(ctx, product); err != nil {
			return err
		}
		for _, uv := range updateVariantsReq {
			variant, err := s.variants.findByID(ctx, uv.ID)
			if err != nil {
				return err
			}
			updateVariantFields(variant, product.Name, uv.Update)
			if err := s.variants.updateOne(ctx, variant); err != nil {
				return err
			}
		}
		// Reload product with variants
		product, err = s.products.findByID(ctx, productID, db.WithPreload(VariantStruct))
		if err != nil {
			return err
		}
		return nil
	})
	return product, err
}

func updateProductFields(product *Product, updateReq *UpdateProductRequest) {
	if updateReq.Name != "" {
		product.Name = updateReq.Name
	}
	if updateReq.Description != "" {
		product.Description = updateReq.Description
	}
	if updateReq.Tags != nil {
		product.Tags = updateReq.Tags
	}
}

func (s *InventoryService) UpdateVariant(ctx context.Context, variantID string, updateReq *UpdateVariantRequest) (*Variant, error) {
	var variant *Variant
	err := s.atomicProcess.Exec(ctx, func(ctx context.Context) error {
		var err error
		variant, err = s.variants.findByID(ctx, variantID, db.WithPreload(ProductStruct))
		if err != nil {
			return err
		}
		updateVariantFields(variant, variant.Product.Name, updateReq)
		if err := s.variants.updateOne(ctx, variant); err != nil {
			return err
		}
		return nil
	})
	return variant, err
}

// updateVariantFields updates the fields of a Variant based on UpdateVariantRequest.
func updateVariantFields(variant *Variant, productName string, updateReq *UpdateVariantRequest) {
	if updateReq.Code != "" {
		variant.Code = updateReq.Code
	}
	variant.Name = fmt.Sprintf("%s - %s", productName, variant.Code)
	if updateReq.SKU != "" {
		variant.SKU = updateReq.SKU
	}
	if updateReq.CostPrice.Sign() > 0 || updateReq.CostPrice.IsZero() {
		variant.CostPrice = updateReq.CostPrice
	}
	if updateReq.SalePrice.Sign() > 0 || updateReq.SalePrice.IsZero() {
		variant.SalePrice = updateReq.SalePrice
	}
	if updateReq.Currency != "" {
		variant.Currency = updateReq.Currency
	}
	if updateReq.StockQuantity >= 0 {
		variant.StockQuantity = updateReq.StockQuantity
	}
	if updateReq.StockAlert >= 0 {
		variant.StockAlert = updateReq.StockAlert
	}
}

func (s *InventoryService) DeleteVariant(ctx context.Context, variantID string) error {
	return s.atomicProcess.Exec(ctx, func(ctx context.Context) error {
		variant, err := s.variants.findByID(ctx, variantID)
		if err != nil {
			return err
		}
		if err := s.variants.deleteOne(ctx, s.variants.scopeID(variant.ID)); err != nil {
			return err
		}
		return nil
	})
}

func (s *InventoryService) DeleteVariantsByProductID(ctx context.Context, productID string) error {
	return s.atomicProcess.Exec(ctx, func(ctx context.Context) error {
		variants, err := s.variants.list(ctx, s.variants.scopeProductID(productID))
		if err != nil {
			return err
		}
		for _, variant := range variants {
			if err := s.variants.deleteOne(ctx, s.variants.scopeID(variant.ID)); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *InventoryService) DeleteProduct(ctx context.Context, productID string) error {
	return s.atomicProcess.Exec(ctx, func(ctx context.Context) error {
		if err := s.DeleteVariantsByProductID(ctx, productID); err != nil {
			return err
		}
		if err := s.products.deleteOne(ctx, s.products.scopeID(productID)); err != nil {
			return err
		}
		return nil
	})
}

func (s *InventoryService) DeleteProducts(ctx context.Context, productIDs []string) error {
	return s.atomicProcess.Exec(ctx, func(ctx context.Context) error {
		for _, productID := range productIDs {
			if err := s.DeleteVariantsByProductID(ctx, productID); err != nil {
				return err
			}
			if err := s.products.deleteOne(ctx, s.products.scopeID(productID)); err != nil {
				return err
			}
		}
		return nil
	})
}

// ---- Analytics wrappers ----

// InventoryTotals returns aggregate inventory metrics for the store.
func (s *InventoryService) InventoryTotals(ctx context.Context, storeID string) (totalValue decimal.Decimal, totalUnits int64, lowStock int64, outOfStock int64, err error) {
	totalValue, err = s.variants.sumInventoryValue(ctx, s.variants.scopeStoreID(storeID))
	if err != nil {
		return
	}
	totalUnits, err = s.variants.sumStockQuantity(ctx, s.variants.scopeStoreID(storeID))
	if err != nil {
		return
	}
	lowStock, err = s.variants.countLowStock(ctx, s.variants.scopeStoreID(storeID))
	if err != nil {
		return
	}
	outOfStock, err = s.variants.countOutOfStock(ctx, s.variants.scopeStoreID(storeID))
	if err != nil {
		return
	}
	return
}

// TopProductsByInventoryValue returns top-N products by inventory value.
func (s *InventoryService) TopProductsByInventoryValue(ctx context.Context, storeID string, limit int) ([]types.KeyValue, error) {
	rows, err := s.variants.topProductsByInventoryValue(ctx, limit, s.variants.scopeStoreID(storeID))
	if err != nil {
		return nil, err
	}
	return rows, nil
}
