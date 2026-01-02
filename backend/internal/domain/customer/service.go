package customer

import (
	"context"
	"strings"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/atomic"
	"github.com/abdelrahman146/kyora/internal/platform/types/list"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/abdelrahman146/kyora/internal/platform/types/timeseries"
	"github.com/abdelrahman146/kyora/internal/platform/utils/transformer"
	"gorm.io/gorm"
)

type Service struct {
	storage *Storage
	bus     *bus.Bus
}

// UpsertCustomerByEmailInput is used by public storefront order submissions.
// Only a minimal set of fields are supported to keep mass-assignment impossible.
type UpsertCustomerByEmailInput struct {
	Email             string
	Name              string
	PhoneNumber       string
	PhoneCode         string
	InstagramUsername string
	TikTokUsername    string
	FacebookUsername  string
	XUsername         string
	SnapchatUsername  string
	WhatsappNumber    string
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
		s.storage.customer.ScopeID(id),
	)
}

// GetCustomerAddressByID returns a customer address by ID after enforcing:
// - customer exists in this business
// - address belongs to that customer
//
// This helper is used by other domains (e.g., orders) to prevent cross-customer or cross-business reference attacks.
func (s *Service) GetCustomerAddressByID(ctx context.Context, actor *account.User, biz *business.Business, customerID string, addressID string) (*CustomerAddress, error) {
	_, err := s.GetCustomerByID(ctx, actor, biz, customerID)
	if err != nil {
		return nil, err
	}
	address, err := s.storage.customerAddress.FindOne(ctx,
		s.storage.customerAddress.ScopeID(addressID),
		s.storage.customerAddress.ScopeEquals(CustomerAddressSchema.CustomerID, customerID),
	)
	if err != nil {
		return nil, err
	}
	return address, nil
}

func (s *Service) CreateCustomer(ctx context.Context, actor *account.User, biz *business.Business, req *CreateCustomerRequest) (*Customer, error) {
	countryCode := strings.ToUpper(strings.TrimSpace(req.CountryCode))
	if countryCode == "" {
		countryCode = req.CountryCode
	}
	customer := &Customer{
		BusinessID:        biz.ID,
		Name:              req.Name,
		Email:             transformer.ToNullableString(req.Email),
		CountryCode:       countryCode,
		Gender:            req.Gender,
		PhoneNumber:       transformer.ToNullableString(req.PhoneNumber),
		PhoneCode:         transformer.ToNullableString(req.PhoneCode),
		TikTokUsername:    transformer.ToNullableString(req.TikTokUsername),
		InstagramUsername: transformer.ToNullableString(req.InstagramUsername),
		FacebookUsername:  transformer.ToNullableString(req.FacebookUsername),
		XUsername:         transformer.ToNullableString(req.XUsername),
		SnapchatUsername:  transformer.ToNullableString(req.SnapchatUsername),
		WhatsappNumber:    transformer.ToNullableString(req.WhatsappNumber),
	}
	err := s.storage.customer.CreateOne(ctx, customer)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

// UpsertCustomerByEmail creates or updates a customer using (businessId, email) as a natural key.
// It is designed for unauthenticated storefront flows.
func (s *Service) UpsertCustomerByEmail(ctx context.Context, biz *business.Business, in *UpsertCustomerByEmailInput) (*Customer, error) {
	if biz == nil {
		return nil, problem.InternalError().With("reason", "business is required")
	}
	if in == nil {
		return nil, problem.BadRequest("customer is required")
	}
	email := strings.TrimSpace(strings.ToLower(in.Email))
	if email == "" {
		return nil, problem.BadRequest("email is required").With("field", "email")
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, problem.BadRequest("name is required").With("field", "name")
	}

	existing, err := s.storage.customer.FindOne(ctx,
		s.storage.customer.ScopeBusinessID(biz.ID),
		s.storage.customer.ScopeEquals(CustomerSchema.Email, email),
	)
	if err != nil {
		if database.IsRecordNotFound(err) {
			cust := &Customer{
				BusinessID:        biz.ID,
				Name:              name,
				Email:             transformer.ToNullableString(email),
				CountryCode:       strings.TrimSpace(strings.ToUpper(biz.CountryCode)),
				Gender:            GenderOther,
				PhoneNumber:       transformer.ToNullableString(strings.TrimSpace(in.PhoneNumber)),
				PhoneCode:         transformer.ToNullableString(strings.TrimSpace(in.PhoneCode)),
				TikTokUsername:    transformer.ToNullableString(strings.TrimSpace(in.TikTokUsername)),
				InstagramUsername: transformer.ToNullableString(strings.TrimSpace(in.InstagramUsername)),
				FacebookUsername:  transformer.ToNullableString(strings.TrimSpace(in.FacebookUsername)),
				XUsername:         transformer.ToNullableString(strings.TrimSpace(in.XUsername)),
				SnapchatUsername:  transformer.ToNullableString(strings.TrimSpace(in.SnapchatUsername)),
				WhatsappNumber:    transformer.ToNullableString(strings.TrimSpace(in.WhatsappNumber)),
				JoinedAt:          time.Now().UTC(),
			}
			if err := s.storage.customer.CreateOne(ctx, cust); err != nil {
				return nil, err
			}
			return cust, nil
		}
		return nil, err
	}

	// Update only when new values are provided.
	existing.Name = name
	if strings.TrimSpace(in.PhoneNumber) != "" {
		existing.PhoneNumber = transformer.ToNullableString(strings.TrimSpace(in.PhoneNumber))
	}
	if strings.TrimSpace(in.PhoneCode) != "" {
		existing.PhoneCode = transformer.ToNullableString(strings.TrimSpace(in.PhoneCode))
	}
	if strings.TrimSpace(in.TikTokUsername) != "" {
		existing.TikTokUsername = transformer.ToNullableString(strings.TrimSpace(in.TikTokUsername))
	}
	if strings.TrimSpace(in.InstagramUsername) != "" {
		existing.InstagramUsername = transformer.ToNullableString(strings.TrimSpace(in.InstagramUsername))
	}
	if strings.TrimSpace(in.FacebookUsername) != "" {
		existing.FacebookUsername = transformer.ToNullableString(strings.TrimSpace(in.FacebookUsername))
	}
	if strings.TrimSpace(in.XUsername) != "" {
		existing.XUsername = transformer.ToNullableString(strings.TrimSpace(in.XUsername))
	}
	if strings.TrimSpace(in.SnapchatUsername) != "" {
		existing.SnapchatUsername = transformer.ToNullableString(strings.TrimSpace(in.SnapchatUsername))
	}
	if strings.TrimSpace(in.WhatsappNumber) != "" {
		existing.WhatsappNumber = transformer.ToNullableString(strings.TrimSpace(in.WhatsappNumber))
	}

	if err := s.storage.customer.UpdateOne(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
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
		customer.CountryCode = strings.ToUpper(strings.TrimSpace(req.CountryCode))
	}
	if req.Gender != "" {
		customer.Gender = req.Gender
	}
	if req.Email != "" {
		customer.Email = transformer.ToNullableString(req.Email)
	}
	if req.PhoneNumber != "" {
		customer.PhoneNumber = transformer.ToNullableString(req.PhoneNumber)
	}
	if req.PhoneCode != "" {
		customer.PhoneCode = transformer.ToNullableString(req.PhoneCode)
	}
	if req.TikTokUsername != "" {
		customer.TikTokUsername = transformer.ToNullableString(req.TikTokUsername)
	}
	if req.InstagramUsername != "" {
		customer.InstagramUsername = transformer.ToNullableString(req.InstagramUsername)
	}
	if req.FacebookUsername != "" {
		customer.FacebookUsername = transformer.ToNullableString(req.FacebookUsername)
	}
	if req.XUsername != "" {
		customer.XUsername = transformer.ToNullableString(req.XUsername)
	}
	if req.SnapchatUsername != "" {
		customer.SnapchatUsername = transformer.ToNullableString(req.SnapchatUsername)
	}
	if req.WhatsappNumber != "" {
		customer.WhatsappNumber = transformer.ToNullableString(req.WhatsappNumber)
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

type ListCustomersFilters struct {
	CountryCode     string
	HasOrders       *bool
	SocialPlatforms []string
}

func (s *Service) ListCustomers(ctx context.Context, actor *account.User, biz *business.Business, req *list.ListRequest, filters *ListCustomersFilters) ([]*CustomerResponse, int64, error) {
	scopes := []func(*gorm.DB) *gorm.DB{
		s.storage.customer.ScopeBusinessID(biz.ID),
	}

	// Apply filters
	if filters != nil {
		// Filter by country
		if filters.CountryCode != "" {
			scopes = append(scopes,
				s.storage.customer.ScopeEquals(CustomerSchema.CountryCode, strings.ToUpper(filters.CountryCode)),
			)
		}

		// Filter by hasOrders
		if filters.HasOrders != nil {
			if *filters.HasOrders {
				scopes = append(scopes,
					s.storage.customer.ScopeWhere("EXISTS (SELECT 1 FROM orders WHERE orders.customer_id = customers.id AND orders.deleted_at IS NULL)"),
				)
			} else {
				scopes = append(scopes,
					s.storage.customer.ScopeWhere("NOT EXISTS (SELECT 1 FROM orders WHERE orders.customer_id = customers.id AND orders.deleted_at IS NULL)"),
				)
			}
		}

		// Filter by social media platforms
		if len(filters.SocialPlatforms) > 0 {
			conditions := []string{}
			for _, platform := range filters.SocialPlatforms {
				switch strings.ToLower(platform) {
				case "instagram":
					conditions = append(conditions, "customers.instagram_username IS NOT NULL AND customers.instagram_username != ''")
				case "tiktok":
					conditions = append(conditions, "customers.tiktok_username IS NOT NULL AND customers.tiktok_username != ''")
				case "facebook":
					conditions = append(conditions, "customers.facebook_username IS NOT NULL AND customers.facebook_username != ''")
				case "x":
					conditions = append(conditions, "customers.x_username IS NOT NULL AND customers.x_username != ''")
				case "snapchat":
					conditions = append(conditions, "customers.snapchat_username IS NOT NULL AND customers.snapchat_username != ''")
				case "whatsapp":
					conditions = append(conditions, "customers.whatsapp_number IS NOT NULL AND customers.whatsapp_number != ''")
				}
			}
			if len(conditions) > 0 {
				scopes = append(scopes,
					s.storage.customer.ScopeWhere("("+strings.Join(conditions, " OR ")+")"),
				)
			}
		}
	}

	// Apply search if provided
	var listOpts []func(*gorm.DB) *gorm.DB
	if req.SearchTerm() != "" {
		term := req.SearchTerm()
		like := "%" + term + "%"
		scopes = append(scopes,
			s.storage.customer.ScopeWhere(
				"(customers.search_vector @@ websearch_to_tsquery('simple', ?) OR customers.name ILIKE ? OR customers.email ILIKE ?)",
				term,
				like,
				like,
			),
		)
		if !req.HasExplicitOrderBy() {
			rankExpr, err := database.WebSearchRankOrder(term, "customers.search_vector")
			if err != nil {
				return nil, 0, problem.InternalError().WithError(err)
			}
			listOpts = append(listOpts, s.storage.customer.WithOrderByExpr(rankExpr))
		}
	}

	findOpts := append([]func(*gorm.DB) *gorm.DB{}, scopes...)
	findOpts = append(findOpts, listOpts...)
	findOpts = append(findOpts,
		s.storage.customer.WithPagination(req.Offset(), req.Limit()),
		s.storage.customer.WithOrderBy(req.ParsedOrderBy(CustomerSchema)),
	)

	// Fetch customers using repository
	customers, err := s.storage.customer.FindMany(ctx, findOpts...)
	if err != nil {
		return nil, 0, err
	}

	// Fetch aggregations in single query
	var aggMap map[string]CustomerAggregation
	if len(customers) > 0 {
		customerIDs := make([]string, len(customers))
		for i, c := range customers {
			customerIDs[i] = c.ID
		}
		aggMap, err = s.storage.GetCustomerAggregations(ctx, biz.ID, customerIDs)
		if err != nil {
			return nil, 0, err
		}
	}

	// Build response DTOs with aggregation data
	responses := make([]*CustomerResponse, len(customers))
	for i, customer := range customers {
		ordersCount := 0
		totalSpent := 0.0
		if agg, found := aggMap[customer.ID]; found {
			ordersCount = agg.OrdersCount
			totalSpent = agg.TotalSpent
		}
		responses[i] = customer.ToResponse(ordersCount, totalSpent)
	}

	// Count with same filters (excluding pagination)
	totalCount, err := s.storage.customer.Count(ctx, scopes...)
	if err != nil {
		return nil, 0, err
	}

	return responses, totalCount, nil
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
	note, err := s.storage.customerNote.FindOne(ctx,
		s.storage.customerNote.ScopeID(noteID),
		s.storage.customerNote.ScopeEquals(CustomerNoteSchema.CustomerID, customerID),
	)
	if err != nil {
		return err
	}
	return s.storage.customerNote.DeleteOne(ctx, note)
}

func (s *Service) CreateCustomerAddress(ctx context.Context, actor *account.User, biz *business.Business, customerID string, req *CreateCustomerAddressRequest) (*CustomerAddress, error) {
	// ensure customer exists within the business
	_, err := s.GetCustomerByID(ctx, actor, biz, customerID)
	if err != nil {
		return nil, err
	}
	countryCode := strings.ToUpper(strings.TrimSpace(req.CountryCode))
	if countryCode == "" {
		countryCode = req.CountryCode
	}
	address := &CustomerAddress{
		CustomerID:  customerID,
		Street:      transformer.ToNullableString(req.Street),
		City:        req.City,
		State:       req.State,
		ZipCode:     transformer.ToNullableString(req.ZipCode),
		CountryCode: countryCode,
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
	address, err := s.storage.customerAddress.FindOne(ctx,
		s.storage.customerAddress.ScopeID(addressID),
		s.storage.customerAddress.ScopeEquals(CustomerAddressSchema.CustomerID, customerID),
	)
	if err != nil {
		return nil, err
	}
	if req.Street != "" {
		address.Street = transformer.ToNullableString(req.Street)
	}
	if req.City != "" {
		address.City = req.City
	}
	if req.State != "" {
		address.State = req.State
	}
	if req.ZipCode != "" {
		address.ZipCode = transformer.ToNullableString(req.ZipCode)
	}
	if req.CountryCode != "" {
		address.CountryCode = strings.ToUpper(strings.TrimSpace(req.CountryCode))
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
	address, err := s.storage.customerAddress.FindOne(ctx,
		s.storage.customerAddress.ScopeID(addressID),
		s.storage.customerAddress.ScopeEquals(CustomerAddressSchema.CustomerID, customerID),
	)
	if err != nil {
		return err
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
