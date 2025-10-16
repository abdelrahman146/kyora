package customer

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/types"
)

type CustomerService struct {
	customers     *customerRepository
	addresses     *addressRepository
	notes         *customerNoteRepository
	atomicProcess *db.AtomicProcess
}

func NewCustomerService(customers *customerRepository, addresses *addressRepository, notes *customerNoteRepository, atomicProcess *db.AtomicProcess) *CustomerService {
	return &CustomerService{
		customers:     customers,
		addresses:     addresses,
		notes:         notes,
		atomicProcess: atomicProcess,
	}
}

func (s *CustomerService) CreateCustomer(ctx context.Context, storeID string, customer *CreateCustomerRequest) (*Customer, error) {
	newCustomer := &Customer{
		StoreID:           storeID,
		FirstName:         customer.FirstName,
		LastName:          customer.LastName,
		Gender:            customer.Gender,
		CountryCode:       customer.CountryCode,
		Email:             customer.Email,
		Phone:             customer.Phone,
		TikTokUsername:    customer.TikTokUsername,
		InstagramUsername: customer.InstagramUsername,
		FacebookUsername:  customer.FacebookUsername,
		XUsername:         customer.XUsername,
		SnapchatUsername:  customer.SnapchatUsername,
		WhatsappNumber:    customer.WhatsappNumber,
		JoinedAt:          customer.JoinedAt,
	}

	if err := s.customers.createOne(ctx, newCustomer); err != nil {
		return nil, err
	}
	return newCustomer, nil
}

func (s *CustomerService) GetCustomerByID(ctx context.Context, storeId string, id string) (*Customer, error) {
	return s.customers.findOne(ctx, s.customers.scopeID(id), s.customers.scopeStoreID(storeId), db.WithPreload(AddressStruct), db.WithPreload(CustomerNoteStruct))
}

func (s *CustomerService) ListCustomers(ctx context.Context, storeID string, listReq *types.ListRequest) ([]*Customer, error) {
	return s.customers.list(ctx, s.customers.scopeStoreID(storeID), db.WithPagination(listReq.Page, listReq.PageSize), db.WithOrderBy(listReq.OrderBy))
}

func (s *CustomerService) CountCustomers(ctx context.Context, storeID string) (int64, error) {
	return s.customers.count(ctx, s.customers.scopeStoreID(storeID))
}

// CountNewCustomersInRange returns number of customers created in the provided range.
func (s *CustomerService) CountNewCustomersInRange(ctx context.Context, storeID string, from, to time.Time) (int64, error) {
	return s.customers.count(ctx, s.customers.scopeStoreID(storeID), s.customers.scopeJoinedAt(from, to))
}

// NewCustomersTimeSeries returns time series of new customer counts.
func (s *CustomerService) NewCustomersTimeSeries(ctx context.Context, storeID string, from, to time.Time, bucket string) ([]types.TimeSeriesRow, error) {
	return s.customers.newCustomersTimeSeries(ctx, bucket, from, to, s.customers.scopeStoreID(storeID))
}

func (s *CustomerService) AddAddressToCustomer(ctx context.Context, customerID string, address *CreateAddressRequest) (*Address, error) {
	newAddress := &Address{
		CustomerID:  customerID,
		Street:      address.Street,
		City:        address.City,
		State:       address.State,
		CountryCode: address.CountryCode,
		Phone:       address.Phone,
		ZipCode:     address.ZipCode,
	}

	if err := s.addresses.createOne(ctx, newAddress); err != nil {
		return nil, err
	}
	return newAddress, nil
}

func (s *CustomerService) AddNoteToCustomer(ctx context.Context, customerID string, note *CreateCustomerNoteRequest) (*CustomerNote, error) {
	newNote := &CustomerNote{
		CustomerID: customerID,
		Note:       note.Note,
	}

	if err := s.notes.createOne(ctx, newNote); err != nil {
		return nil, err
	}
	return newNote, nil
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, customerID string, updates *UpdateCustomerRequest) (*Customer, error) {
	existingCustomer, err := s.customers.findByID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	if updates.FirstName != "" {
		existingCustomer.FirstName = updates.FirstName
	}
	if updates.LastName != "" {
		existingCustomer.LastName = updates.LastName
	}
	if updates.Gender != "" {
		existingCustomer.Gender = updates.Gender
	}
	if updates.CountryCode != "" {
		existingCustomer.CountryCode = updates.CountryCode
	}
	if updates.Email != "" {
		existingCustomer.Email = updates.Email
	}
	if updates.Phone != "" {
		existingCustomer.Phone = updates.Phone
	}
	if updates.TikTokUsername != "" {
		existingCustomer.TikTokUsername = updates.TikTokUsername
	}
	if updates.InstagramUsername != "" {
		existingCustomer.InstagramUsername = updates.InstagramUsername
	}
	if updates.FacebookUsername != "" {
		existingCustomer.FacebookUsername = updates.FacebookUsername
	}
	if updates.XUsername != "" {
		existingCustomer.XUsername = updates.XUsername
	}
	if updates.SnapchatUsername != "" {
		existingCustomer.SnapchatUsername = updates.SnapchatUsername
	}
	if updates.WhatsappNumber != "" {
		existingCustomer.WhatsappNumber = updates.WhatsappNumber
	}
	if !updates.JoinedAt.IsZero() {
		existingCustomer.JoinedAt = updates.JoinedAt
	}

	if err := s.customers.updateOne(ctx, existingCustomer); err != nil {
		return nil, err
	}
	return existingCustomer, nil
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, customerID string) error {
	return s.customers.deleteOne(ctx, s.customers.scopeID(customerID))
}

func (s *CustomerService) DeleteCustomers(ctx context.Context, customerIDs []string) error {
	return s.customers.deleteMany(ctx, s.customers.scopeIDs(customerIDs))
}

func (s *CustomerService) DeleteAllCustomersInStore(ctx context.Context, storeID string) error {
	return s.customers.deleteMany(ctx, s.customers.scopeStoreID(storeID))
}

func (s *CustomerService) ListAddressesOfCustomer(ctx context.Context, customerID string) ([]*Address, error) {
	return s.addresses.list(ctx, s.addresses.scopeCustomerID(customerID))
}

func (s *CustomerService) ListNotesOfCustomer(ctx context.Context, customerID string) ([]*CustomerNote, error) {
	return s.notes.list(ctx, s.notes.scopeCustomerID(customerID))
}

func (s *CustomerService) GetAddressByID(ctx context.Context, addressID string) (*Address, error) {
	return s.addresses.findByID(ctx, addressID)
}

func (s *CustomerService) GetNoteByID(ctx context.Context, noteID string) (*CustomerNote, error) {
	return s.notes.findByID(ctx, noteID)
}
func (s *CustomerService) DeleteAddress(ctx context.Context, addressID string) error {
	return s.addresses.deleteOne(ctx, s.addresses.scopeID(addressID))
}

func (s *CustomerService) DeleteNote(ctx context.Context, noteID string) error {
	return s.notes.deleteOne(ctx, s.notes.scopeID(noteID))
}
func (s *CustomerService) DeleteAllAddressesOfCustomer(ctx context.Context, customerID string) error {
	return s.addresses.deleteMany(ctx, s.addresses.scopeCustomerID(customerID))
}

func (s *CustomerService) DeleteAllNotesOfCustomer(ctx context.Context, customerID string) error {
	return s.notes.deleteMany(ctx, s.notes.scopeCustomerID(customerID))
}

func (s *CustomerService) UpdateAddress(ctx context.Context, addressID string, updates *UpdateAddressRequest) (*Address, error) {
	existingAddress, err := s.addresses.findByID(ctx, addressID)
	if err != nil {
		return nil, err
	}

	if updates.Street != "" {
		existingAddress.Street = updates.Street
	}
	if updates.City != "" {
		existingAddress.City = updates.City
	}
	if updates.State != "" {
		existingAddress.State = updates.State
	}
	if updates.CountryCode != "" {
		existingAddress.CountryCode = updates.CountryCode
	}
	if updates.Phone != "" {
		existingAddress.Phone = updates.Phone
	}
	if updates.ZipCode != "" {
		existingAddress.ZipCode = updates.ZipCode
	}

	if err := s.addresses.updateOne(ctx, existingAddress); err != nil {
		return nil, err
	}
	return existingAddress, nil
}

func (s *CustomerService) UpdateNote(ctx context.Context, noteID string, updates *UpdateCustomerNoteRequest) (*CustomerNote, error) {
	existingNote, err := s.notes.findByID(ctx, noteID)
	if err != nil {
		return nil, err
	}

	if updates.Note != "" {
		existingNote.Note = updates.Note
	}

	if err := s.notes.updateOne(ctx, existingNote); err != nil {
		return nil, err
	}
	return existingNote, nil
}
