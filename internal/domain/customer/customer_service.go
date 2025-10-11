package customer

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"
)

type CustomerService struct {
	customers     *CustomerRepository
	addresses     *AddressRepository
	notes         *CustomerNoteRepository
	atomicProcess *db.AtomicProcess
}

func NewCustomerService(customers *CustomerRepository, addresses *AddressRepository, notes *CustomerNoteRepository, atomicProcess *db.AtomicProcess) *CustomerService {
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
	}

	if err := s.customers.CreateOne(ctx, newCustomer); err != nil {
		return nil, err
	}
	return newCustomer, nil
}

func (s *CustomerService) GetCustomerByID(ctx context.Context, id string) (*Customer, error) {
	return s.customers.FindByID(ctx, id, db.WithPreload(AddressStruct), db.WithPreload(CustomerNoteStruct))
}

func (s *CustomerService) ListCustomers(ctx context.Context, storeID string, page, pageSize int, orderBy string, ascending bool) ([]*Customer, error) {
	return s.customers.List(ctx, s.customers.ScopeStoreID(storeID), db.WithPagination(page, pageSize), db.WithSorting(orderBy, ascending))
}

func (s *CustomerService) CountCustomers(ctx context.Context, storeID string) (int64, error) {
	return s.customers.Count(ctx, s.customers.ScopeStoreID(storeID))
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

	if err := s.addresses.CreateOne(ctx, newAddress); err != nil {
		return nil, err
	}
	return newAddress, nil
}

func (s *CustomerService) AddNoteToCustomer(ctx context.Context, customerID string, note *CreateCustomerNoteRequest) (*CustomerNote, error) {
	newNote := &CustomerNote{
		CustomerID: customerID,
		Note:       note.Note,
	}

	if err := s.notes.CreateOne(ctx, newNote); err != nil {
		return nil, err
	}
	return newNote, nil
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, customerID string, updates *UpdateCustomerRequest) (*Customer, error) {
	existingCustomer, err := s.customers.FindByID(ctx, customerID)
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

	if err := s.customers.UpdateOne(ctx, existingCustomer); err != nil {
		return nil, err
	}
	return existingCustomer, nil
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, customerID string) error {
	return s.customers.DeleteOne(ctx, s.customers.ScopeID(customerID))
}

func (s *CustomerService) DeleteCustomers(ctx context.Context, customerIDs []string) error {
	return s.customers.DeleteMany(ctx, s.customers.ScopeIDs(customerIDs))
}

func (s *CustomerService) DeleteAllCustomersInStore(ctx context.Context, storeID string) error {
	return s.customers.DeleteMany(ctx, s.customers.ScopeStoreID(storeID))
}

func (s *CustomerService) ListAddressesOfCustomer(ctx context.Context, customerID string) ([]*Address, error) {
	return s.addresses.List(ctx, s.addresses.ScopeCustomerID(customerID))
}

func (s *CustomerService) ListNotesOfCustomer(ctx context.Context, customerID string) ([]*CustomerNote, error) {
	return s.notes.List(ctx, s.notes.ScopeCustomerID(customerID))
}

func (s *CustomerService) GetAddressByID(ctx context.Context, addressID string) (*Address, error) {
	return s.addresses.FindByID(ctx, addressID)
}

func (s *CustomerService) GetNoteByID(ctx context.Context, noteID string) (*CustomerNote, error) {
	return s.notes.FindByID(ctx, noteID)
}
func (s *CustomerService) DeleteAddress(ctx context.Context, addressID string) error {
	return s.addresses.DeleteOne(ctx, s.addresses.ScopeID(addressID))
}

func (s *CustomerService) DeleteNote(ctx context.Context, noteID string) error {
	return s.notes.DeleteOne(ctx, s.notes.ScopeID(noteID))
}
func (s *CustomerService) DeleteAllAddressesOfCustomer(ctx context.Context, customerID string) error {
	return s.addresses.DeleteMany(ctx, s.addresses.ScopeCustomerID(customerID))
}

func (s *CustomerService) DeleteAllNotesOfCustomer(ctx context.Context, customerID string) error {
	return s.notes.DeleteMany(ctx, s.notes.ScopeCustomerID(customerID))
}

func (s *CustomerService) UpdateAddress(ctx context.Context, addressID string, updates *UpdateAddressRequest) (*Address, error) {
	existingAddress, err := s.addresses.FindByID(ctx, addressID)
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

	if err := s.addresses.UpdateOne(ctx, existingAddress); err != nil {
		return nil, err
	}
	return existingAddress, nil
}

func (s *CustomerService) UpdateNote(ctx context.Context, noteID string, updates *UpdateCustomerNoteRequest) (*CustomerNote, error) {
	existingNote, err := s.notes.FindByID(ctx, noteID)
	if err != nil {
		return nil, err
	}

	if updates.Note != "" {
		existingNote.Note = updates.Note
	}

	if err := s.notes.UpdateOne(ctx, existingNote); err != nil {
		return nil, err
	}
	return existingNote, nil
}
