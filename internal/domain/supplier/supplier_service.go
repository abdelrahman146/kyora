package supplier

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/store"
)

type SupplierService struct {
	storeService *store.StoreService
	supplierRepo *SupplierRepository
}

func NewSupplierService(storeService *store.StoreService, supplierRepo *SupplierRepository) *SupplierService {
	return &SupplierService{storeService: storeService, supplierRepo: supplierRepo}
}

func (s *SupplierService) GetSupplierByID(ctx context.Context, storeID, id string, opts ...db.PostgresOptions) (*Supplier, error) {
	return s.supplierRepo.FindByID(ctx, id, opts...)
}

func (s *SupplierService) ListSuppliers(ctx context.Context, storeID string, page, pageSize int, orderBy string, ascending bool) ([]*Supplier, error) {
	return s.supplierRepo.List(ctx, s.supplierRepo.ScopeStoreID(storeID), db.WithPagination(page, pageSize), db.WithSorting(orderBy, ascending))
}

func (s *SupplierService) CountSuppliers(ctx context.Context, storeID string) (int64, error) {
	return s.supplierRepo.Count(ctx, s.supplierRepo.ScopeStoreID(storeID))
}

func (s *SupplierService) CreateSupplier(ctx context.Context, storeID string, supplier *CreateSupplierRequest) (*Supplier, error) {
	store, err := s.storeService.GetStoreByID(ctx, storeID)
	if err != nil {
		return nil, err
	}

	newSupplier := &Supplier{
		StoreID:     store.ID,
		Name:        supplier.Name,
		Contact:     supplier.Contact,
		Email:       supplier.Email,
		Phone:       supplier.Phone,
		CountryCode: supplier.CountryCode,
		Website:     supplier.Website,
	}

	if err := s.supplierRepo.CreateOne(ctx, newSupplier); err != nil {
		return nil, err
	}
	return newSupplier, nil
}

func (s *SupplierService) UpdateSupplier(ctx context.Context, storeID, id string, supplier *UpdateSupplierRequest) (*Supplier, error) {
	existingSupplier, err := s.supplierRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if supplier.Name != "" {
		existingSupplier.Name = supplier.Name
	}
	if supplier.Contact != "" {
		existingSupplier.Contact = supplier.Contact
	}
	if supplier.Email != "" {
		existingSupplier.Email = supplier.Email
	}
	if supplier.Phone != "" {
		existingSupplier.Phone = supplier.Phone
	}
	if supplier.Website != "" {
		existingSupplier.Website = supplier.Website
	}

	if supplier.CountryCode != "" {
		existingSupplier.CountryCode = supplier.CountryCode
	}

	if err := s.supplierRepo.UpdateOne(ctx, existingSupplier); err != nil {
		return nil, err
	}
	return existingSupplier, nil
}

func (s *SupplierService) DeleteSupplier(ctx context.Context, storeID, id string) error {
	_, err := s.supplierRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return s.supplierRepo.DeleteOne(ctx, s.supplierRepo.ScopeID(id), s.supplierRepo.ScopeStoreID(storeID))
}
