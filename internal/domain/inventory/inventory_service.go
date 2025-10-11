package inventory

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"
)

type InventoryService struct {
	products      *ProductRepository
	variants      *VariantRepository
	atomicProcess *db.AtomicProcess
}

func NewInventoryService(products *ProductRepository, variants *VariantRepository, atomicProcess *db.AtomicProcess) *InventoryService {
	return &InventoryService{
		products:      products,
		variants:      variants,
		atomicProcess: atomicProcess,
	}
}

func (s *InventoryService) CreateProduct(ctx context.Context, productReq *CreateProductRequest, variantsReq []*CreateVariantRequest) (*Product, error) {
	var product *Product
	err := s.atomicProcess.Exec(ctx, func(ctx context.Context) error {
		product = &Product{
			Name:        productReq.Name,
			Description: productReq.Description,
			Tags:        productReq.Tags,
		}
		if err := s.products.CreateOne(ctx, product); err != nil {
			return db.HandleDBError(err)
		}
		var variants []*Variant
		for _, v := range variantsReq {
			variant := &Variant{
				Name:          v.Name,
				SKU:           v.SKU,
				ProductID:     product.ID,
				CostPrice:     v.CostPrice,
				SalePrice:     v.SalePrice,
				StockQuantity: v.StockQuantity,
				StockAlert:    v.StockAlert,
			}
			variants = append(variants, variant)
		}
		if len(variants) > 0 {
			if err := s.variants.CreateMany(ctx, variants); err != nil {
				return db.HandleDBError(err)
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
		variant = &Variant{
			Name:          variantReq.Name,
			SKU:           variantReq.SKU,
			ProductID:     productID,
			CostPrice:     variantReq.CostPrice,
			SalePrice:     variantReq.SalePrice,
			StockQuantity: variantReq.StockQuantity,
			StockAlert:    variantReq.StockAlert,
		}
		if err := s.variants.CreateOne(ctx, variant); err != nil {
			return db.HandleDBError(err)
		}
		return nil
	})
	return variant, err
}

func (s *InventoryService) GetProductByID(ctx context.Context, productID string) (*Product, error) {
	product, err := s.products.FindByID(ctx, productID, db.WithPreload(VariantStruct))
	if err != nil {
		return nil, db.HandleDBError(err)
	}
	return product, nil
}

func (s *InventoryService) GetVariantByID(ctx context.Context, variantID string) (*Variant, error) {
	variant, err := s.variants.FindByID(ctx, variantID, db.WithPreload(ProductStruct))
	if err != nil {
		return nil, db.HandleDBError(err)
	}
	return variant, nil
}

func (s *InventoryService) GetVariantBySKU(ctx context.Context, sku string) (*Variant, error) {
	variant, err := s.variants.FindOne(ctx, s.variants.ScopeSKU(sku), db.WithPreload(ProductStruct))
	if err != nil {
		return nil, db.HandleDBError(err)
	}
	return variant, nil
}

func (s *InventoryService) ListVariantsByProductID(ctx context.Context, productID string) ([]*Variant, error) {
	variants, err := s.variants.List(ctx, s.variants.ScopeProductID(productID))
	if err != nil {
		return nil, db.HandleDBError(err)
	}
	return variants, nil
}

func (s *InventoryService) ListVariants(ctx context.Context, filter *VariantFilter, page int, pageSize int, orderBy string, ascending bool) ([]*Variant, error) {
	variants, err := s.variants.List(ctx, s.variants.ScopeFilter(filter), db.WithPagination(page, pageSize), db.WithPreload(ProductStruct), db.WithSorting(orderBy, ascending))
	if err != nil {
		return nil, db.HandleDBError(err)
	}
	return variants, nil
}

func (s *InventoryService) CountVariants(ctx context.Context, filter *VariantFilter) (int64, error) {
	count, err := s.variants.Count(ctx, s.variants.ScopeFilter(filter))
	if err != nil {
		return 0, db.HandleDBError(err)
	}
	return count, nil
}

func (s *InventoryService) ListProducts(ctx context.Context, filter *ProductFilter, page int, pageSize int, orderBy string, ascending bool) ([]*Product, error) {
	products, err := s.products.List(ctx, s.products.ScopeFilter(filter), db.WithPagination(page, pageSize), db.WithPreload(VariantStruct), db.WithSorting(orderBy, ascending))
	if err != nil {
		return nil, db.HandleDBError(err)
	}
	return products, nil
}

func (s *InventoryService) CountProducts(ctx context.Context, filter *ProductFilter) (int64, error) {
	count, err := s.products.Count(ctx, s.products.ScopeFilter(filter))
	if err != nil {
		return 0, db.HandleDBError(err)
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
		product, err = s.products.FindByID(ctx, productID)
		if err != nil {
			return db.HandleDBError(err)
		}
		updateProductFields(product, updateReq)
		if err := s.products.UpdateOne(ctx, product); err != nil {
			return db.HandleDBError(err)
		}
		for _, uv := range updateVariantsReq {
			variant, err := s.variants.FindByID(ctx, uv.ID)
			if err != nil {
				return db.HandleDBError(err)
			}
			updateVariantFields(variant, uv.Update)
			if err := s.variants.UpdateOne(ctx, variant); err != nil {
				return db.HandleDBError(err)
			}
		}
		// Reload product with variants
		product, err = s.products.FindByID(ctx, productID, db.WithPreload(VariantStruct))
		if err != nil {
			return db.HandleDBError(err)
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
		variant, err = s.variants.FindByID(ctx, variantID)
		if err != nil {
			return db.HandleDBError(err)
		}
		updateVariantFields(variant, updateReq)
		if err := s.variants.UpdateOne(ctx, variant); err != nil {
			return db.HandleDBError(err)
		}
		return nil
	})
	return variant, err
}

// updateVariantFields updates the fields of a Variant based on UpdateVariantRequest.
func updateVariantFields(variant *Variant, updateReq *UpdateVariantRequest) {
	if updateReq.Name != "" {
		variant.Name = updateReq.Name
	}
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
		variant, err := s.variants.FindByID(ctx, variantID)
		if err != nil {
			return db.HandleDBError(err)
		}
		if err := s.variants.DeleteOne(ctx, s.variants.ScopeID(variant.ID)); err != nil {
			return db.HandleDBError(err)
		}
		return nil
	})
}

func (s *InventoryService) DeleteVariantsByProductID(ctx context.Context, productID string) error {
	return s.atomicProcess.Exec(ctx, func(ctx context.Context) error {
		variants, err := s.variants.List(ctx, s.variants.ScopeProductID(productID))
		if err != nil {
			return db.HandleDBError(err)
		}
		for _, variant := range variants {
			if err := s.variants.DeleteOne(ctx, s.variants.ScopeID(variant.ID)); err != nil {
				return db.HandleDBError(err)
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
		if err := s.products.DeleteOne(ctx, s.products.ScopeID(productID)); err != nil {
			return db.HandleDBError(err)
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
			if err := s.products.DeleteOne(ctx, s.products.ScopeID(productID)); err != nil {
				return db.HandleDBError(err)
			}
		}
		return nil
	})
}
