package customer

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/abdelrahman146/kyora/internal/platform/types/atomic"
	"github.com/abdelrahman146/kyora/internal/platform/types/list"
	"github.com/abdelrahman146/kyora/internal/platform/types/timeseries"
	"github.com/abdelrahman146/kyora/internal/platform/utils/transformer"
	"gorm.io/gorm"
)

type Service struct {
	storage *Storage
	bus     *bus.Bus
}

func NewService(storage *Storage, atomicProcessor atomic.AtomicProcessor, bus *bus.Bus) *Service {
	return &Service{
		bus:     bus,
		storage: storage,
	}
}

func (s *Service) GetCustomerByID(ctx context.Context, actor *account.User, biz *business.Business, id string) (*Customer, error) {
	return s.storage.customer.FindOne(ctx,
		s.storage.customer.ScopeBusinessID(biz.ID),
		s.storage.customer.ScopeID(id), s.storage.customer.WithPreload(CustomerAddressStruct),
		s.storage.customer.WithPreload(CustomerAddressStruct),
		s.storage.customer.WithPreload(CustomerNoteStruct),
	)
}

func (s *Service) CreateCustomer(ctx context.Context, actor *account.User, biz *business.Business, req *CreateCustomerRequest) (*Customer, error) {
	customer := &Customer{
		BusinessID:        biz.ID,
		Name:              req.Name,
		Email:             transformer.ToNullString(req.Email),
		CountryCode:       req.CountryCode,
		Gender:            req.Gender,
		PhoneNumber:       transformer.ToNullString(req.PhoneNumber),
		PhoneCode:         transformer.ToNullString(req.PhoneCode),
		TikTokUsername:    transformer.ToNullString(req.TikTokUsername),
		InstagramUsername: transformer.ToNullString(req.InstagramUsername),
		FacebookUsername:  transformer.ToNullString(req.FacebookUsername),
		XUsername:         transformer.ToNullString(req.XUsername),
		SnapchatUsername:  transformer.ToNullString(req.SnapchatUsername),
		WhatsappNumber:    transformer.ToNullString(req.WhatsappNumber),
	}
	err := s.storage.customer.CreateOne(ctx, customer)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *Service) UpdateCustomer(ctx context.Context, actor *account.User, biz *business.Business, id string, req *UpdateCustomerRequest) (*Customer, error) {
	customer, err := s.GetCustomerByID(ctx, actor, biz, id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		customer.Name = req.Name
	}
	if req.CountryCode != "" {
		customer.CountryCode = req.CountryCode
	}
	if req.Gender != "" {
		customer.Gender = req.Gender
	}
	if req.Email != "" {
		customer.Email = transformer.ToNullString(req.Email)
	}
	if req.PhoneNumber != "" {
		customer.PhoneNumber = transformer.ToNullString(req.PhoneNumber)
	}
	if req.PhoneCode != "" {
		customer.PhoneCode = transformer.ToNullString(req.PhoneCode)
	}
	if req.TikTokUsername != "" {
		customer.TikTokUsername = transformer.ToNullString(req.TikTokUsername)
	}
	if req.InstagramUsername != "" {
		customer.InstagramUsername = transformer.ToNullString(req.InstagramUsername)
	}
	if req.FacebookUsername != "" {
		customer.FacebookUsername = transformer.ToNullString(req.FacebookUsername)
	}
	if req.XUsername != "" {
		customer.XUsername = transformer.ToNullString(req.XUsername)
	}
	if req.SnapchatUsername != "" {
		customer.SnapchatUsername = transformer.ToNullString(req.SnapchatUsername)
	}
	if req.WhatsappNumber != "" {
		customer.WhatsappNumber = transformer.ToNullString(req.WhatsappNumber)
	}
	err = s.storage.customer.UpdateOne(ctx, customer)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *Service) DeleteCustomer(ctx context.Context, actor *account.User, biz *business.Business, id string) error {
	customer, err := s.GetCustomerByID(ctx, actor, biz, id)
	if err != nil {
		return err
	}
	return s.storage.customer.DeleteOne(ctx, customer)
}

func (s *Service) ListCustomers(ctx context.Context, actor *account.User, biz *business.Business, req *list.ListRequest) ([]*Customer, int64, error) {
	scopes := []func(*gorm.DB) *gorm.DB{
		s.storage.customer.ScopeBusinessID(biz.ID),
	}

	// Apply search if provided
	if req.SearchTerm != "" {
		searchScopes := []func(*gorm.DB) *gorm.DB{
			s.storage.customer.ScopeSearchTerm(req.SearchTerm, CustomerSchema.Name, CustomerSchema.Email),
		}
		scopes = append(scopes, searchScopes...)
	}

	customers, err := s.storage.customer.FindMany(ctx,
		append(scopes,
			s.storage.customer.WithPagination(req.Offset(), req.Limit()),
			s.storage.customer.WithOrderBy(req.ParsedOrderBy(CustomerSchema)),
		)...,
	)
	if err != nil {
		return nil, 0, err
	}

	// Count with same filters (excluding pagination)
	totalCount, err := s.storage.customer.Count(ctx, scopes...)
	if err != nil {
		return nil, 0, err
	}

	return customers, totalCount, nil
}

func (s *Service) CountCustomers(ctx context.Context, actor *account.User, biz *business.Business) (int64, error) {
	return s.storage.customer.Count(ctx,
		s.storage.customer.ScopeBusinessID(biz.ID),
	)
}

func (s *Service) CountCustomersByDateRange(ctx context.Context, actor *account.User, biz *business.Business, from, to time.Time) (int64, error) {
	return s.storage.customer.Count(ctx,
		s.storage.customer.ScopeBusinessID(biz.ID),
		s.storage.customer.ScopeTime(CustomerSchema.CreatedAt, from, to),
	)
}

// GetCustomersByIDs fetches customers by a list of IDs scoped to the business.
func (s *Service) GetCustomersByIDs(ctx context.Context, actor *account.User, biz *business.Business, ids []any) ([]*Customer, error) {
	if len(ids) == 0 {
		return []*Customer{}, nil
	}
	return s.storage.customer.FindMany(ctx,
		s.storage.customer.ScopeBusinessID(biz.ID),
		s.storage.customer.ScopeIDs(ids),
	)
}

func (s *Service) CreateCustomerNote(ctx context.Context, actor *account.User, biz *business.Business, customerID string, content string) (*CustomerNote, error) {
	// ensure customer exists within the business
	_, err := s.GetCustomerByID(ctx, actor, biz, customerID)
	if err != nil {
		return nil, err
	}
	note := &CustomerNote{
		CustomerID: customerID,
		Content:    content,
	}
	err = s.storage.customerNote.CreateOne(ctx, note)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (s *Service) ListCustomerNotes(ctx context.Context, actor *account.User, biz *business.Business, customerID string) ([]*CustomerNote, error) {
	// ensure customer exists within the business
	_, err := s.GetCustomerByID(ctx, actor, biz, customerID)
	if err != nil {
		return nil, err
	}
	return s.storage.customerNote.FindMany(ctx,
		s.storage.customerNote.ScopeEquals(CustomerNoteSchema.CustomerID, customerID),
	)
}

func (s *Service) CountCustomerNotes(ctx context.Context, actor *account.User, biz *business.Business, customerID string) (int64, error) {
	// ensure customer exists within the business
	_, err := s.GetCustomerByID(ctx, actor, biz, customerID)
	if err != nil {
		return 0, err
	}
	return s.storage.customerNote.Count(ctx,
		s.storage.customerNote.ScopeEquals(CustomerNoteSchema.CustomerID, customerID),
	)
}

func (s *Service) DeleteCustomerNote(ctx context.Context, actor *account.User, biz *business.Business, customerID string, noteID string) error {
	// ensure customer exists within the business
	_, err := s.GetCustomerByID(ctx, actor, biz, customerID)
	if err != nil {
		return err
	}
	// ensure note exists within the customer
	note, err := s.storage.customerNote.FindByID(ctx, noteID)
	if err != nil {
		return err
	}
	if note.CustomerID != customerID {
		return gorm.ErrRecordNotFound
	}
	return s.storage.customerNote.DeleteOne(ctx, note)
}

func (s *Service) CreateCustomerAddress(ctx context.Context, actor *account.User, biz *business.Business, customerID string, req *CreateCustomerAddressRequest) (*CustomerAddress, error) {
	// ensure customer exists within the business
	_, err := s.GetCustomerByID(ctx, actor, biz, customerID)
	if err != nil {
		return nil, err
	}
	address := &CustomerAddress{
		CustomerID:  customerID,
		Street:      transformer.ToNullString(req.Street),
		City:        req.City,
		State:       req.State,
		ZipCode:     transformer.ToNullString(req.ZipCode),
		CountryCode: req.CountryCode,
		PhoneCode:   req.PhoneCode,
		PhoneNumber: req.Phone,
	}
	err = s.storage.customerAddress.CreateOne(ctx, address)
	if err != nil {
		return nil, err
	}
	return address, nil
}

func (s *Service) UpdateCustomerAddress(ctx context.Context, actor *account.User, biz *business.Business, customerID string, addressID string, req *UpdateCustomerAddressRequest) (*CustomerAddress, error) {
	// ensure customer exists within the business
	_, err := s.GetCustomerByID(ctx, actor, biz, customerID)
	if err != nil {
		return nil, err
	}
	// ensure address exists within the customer
	address, err := s.storage.customerAddress.FindByID(ctx, addressID)
	if err != nil {
		return nil, err
	}
	if address.CustomerID != customerID {
		return nil, gorm.ErrRecordNotFound
	}
	if req.Street != "" {
		address.Street = transformer.ToNullString(req.Street)
	}
	if req.City != "" {
		address.City = req.City
	}
	if req.State != "" {
		address.State = req.State
	}
	if req.ZipCode != "" {
		address.ZipCode = transformer.ToNullString(req.ZipCode)
	}
	if req.CountryCode != "" {
		address.CountryCode = req.CountryCode
	}
	err = s.storage.customerAddress.UpdateOne(ctx, address)
	if err != nil {
		return nil, err
	}
	return address, nil
}

func (s *Service) ListCustomerAddresses(ctx context.Context, actor *account.User, biz *business.Business, customerID string) ([]*CustomerAddress, error) {
	// ensure customer exists within the business
	_, err := s.GetCustomerByID(ctx, actor, biz, customerID)
	if err != nil {
		return nil, err
	}
	return s.storage.customerAddress.FindMany(ctx,
		s.storage.customerAddress.ScopeEquals(CustomerAddressSchema.CustomerID, customerID),
	)
}

func (s *Service) DeleteCustomerAddress(ctx context.Context, actor *account.User, biz *business.Business, customerID string, addressID string) error {
	// ensure customer exists within the business
	_, err := s.GetCustomerByID(ctx, actor, biz, customerID)
	if err != nil {
		return err
	}
	// ensure address exists within the customer
	address, err := s.storage.customerAddress.FindByID(ctx, addressID)
	if err != nil {
		return err
	}
	if address.CustomerID != customerID {
		return gorm.ErrRecordNotFound
	}
	return s.storage.customerAddress.DeleteOne(ctx, address)
}

func (s *Service) CountCustomerAddresses(ctx context.Context, actor *account.User, biz *business.Business, customerID string) (int64, error) {
	// ensure customer exists within the business
	_, err := s.GetCustomerByID(ctx, actor, biz, customerID)
	if err != nil {
		return 0, err
	}
	return s.storage.customerAddress.Count(ctx,
		s.storage.customerAddress.ScopeEquals(CustomerAddressSchema.CustomerID, customerID),
	)
}

func (s *Service) ComputeCustomersTimeSeries(ctx context.Context, actor *account.User, biz *business.Business, from, to time.Time) (*timeseries.TimeSeries, error) {
	granularity := timeseries.GetTimeGranularityByDateRange(from, to)
	return s.storage.customer.TimeSeriesCount(ctx, CustomerSchema.JoinedAt, granularity,
		s.storage.customer.ScopeBusinessID(biz.ID),
		s.storage.customer.ScopeTime(CustomerSchema.JoinedAt, from, to),
	)
}
