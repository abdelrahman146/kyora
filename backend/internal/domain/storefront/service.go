package storefront

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/domain/order"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/atomic"
	"github.com/abdelrahman146/kyora/internal/platform/types/list"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/abdelrahman146/kyora/internal/platform/utils/throttle"
)

type Service struct {
	storage         *Storage
	atomicProcessor atomic.AtomicProcessor
	business        *business.Service
	inventory       *inventory.Service
	customer        *customer.Service
	orders          *order.Service
}

func NewService(storage *Storage, atomicProcessor atomic.AtomicProcessor, businessSvc *business.Service, inventorySvc *inventory.Service, customerSvc *customer.Service, orderSvc *order.Service) *Service {
	return &Service{
		storage:         storage,
		atomicProcessor: atomicProcessor,
		business:        businessSvc,
		inventory:       inventorySvc,
		customer:        customerSvc,
		orders:          orderSvc,
	}
}

type PublicBusiness struct {
	ID                 string                   `json:"id"`
	Name               string                   `json:"name"`
	Descriptor         string                   `json:"descriptor"`
	Brand              string                   `json:"brand,omitempty"`
	CountryCode        string                   `json:"countryCode"`
	Currency           string                   `json:"currency"`
	StorefrontPublicID string                   `json:"storefrontPublicId"`
	StorefrontEnabled  bool                     `json:"storefrontEnabled"`
	StorefrontTheme    business.StorefrontTheme `json:"storefrontTheme"`
	SupportEmail       string                   `json:"supportEmail,omitempty"`
	PhoneNumber        string                   `json:"phoneNumber,omitempty"`
	WhatsappNumber     string                   `json:"whatsappNumber,omitempty"`
	Address            string                   `json:"address,omitempty"`
	WebsiteURL         string                   `json:"websiteUrl,omitempty"`
	InstagramURL       string                   `json:"instagramUrl,omitempty"`
	FacebookURL        string                   `json:"facebookUrl,omitempty"`
	TikTokURL          string                   `json:"tiktokUrl,omitempty"`
	XURL               string                   `json:"xUrl,omitempty"`
	SnapchatURL        string                   `json:"snapchatUrl,omitempty"`
}

type PublicVariant struct {
	ID        string                 `json:"id"`
	ProductID string                 `json:"productId"`
	Code      string                 `json:"code"`
	Name      string                 `json:"name"`
	SKU       string                 `json:"sku"`
	SalePrice string                 `json:"salePrice"`
	Currency  string                 `json:"currency"`
	Photos    inventory.PhotoURLList `json:"photos,omitempty"`
}

type PublicProduct struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	CategoryID  string                 `json:"categoryId"`
	Photos      inventory.PhotoURLList `json:"photos,omitempty"`
	Variants    []PublicVariant        `json:"variants"`
}

type PublicCategory struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Descriptor string `json:"descriptor"`
}

type CatalogResponse struct {
	Business   PublicBusiness   `json:"business"`
	Categories []PublicCategory `json:"categories"`
	Products   []PublicProduct  `json:"products"`
}

func (s *Service) GetCatalog(ctx context.Context, storefrontPublicID string) (*CatalogResponse, error) {
	biz, err := s.business.GetBusinessByStorefrontPublicID(ctx, storefrontPublicID)
	if err != nil {
		return nil, ErrStorefrontNotFound(storefrontPublicID, err)
	}
	if !biz.StorefrontEnabled {
		return nil, ErrStorefrontDisabled(storefrontPublicID)
	}

	cats, err := s.inventory.ListCategories(ctx, nil, biz)
	if err != nil {
		return nil, err
	}
	prods, err := s.listAllProducts(ctx, biz)
	if err != nil {
		return nil, err
	}
	vars, err := s.listAllVariants(ctx, biz)
	if err != nil {
		return nil, err
	}

	variantsByProduct := map[string][]PublicVariant{}
	for _, v := range vars {
		variantsByProduct[v.ProductID] = append(variantsByProduct[v.ProductID], PublicVariant{
			ID:        v.ID,
			ProductID: v.ProductID,
			Code:      v.Code,
			Name:      v.Name,
			SKU:       v.SKU,
			SalePrice: v.SalePrice.String(),
			Currency:  v.Currency,
			Photos:    v.Photos,
		})
	}

	outCats := make([]PublicCategory, 0, len(cats))
	for _, c := range cats {
		outCats = append(outCats, PublicCategory{ID: c.ID, Name: c.Name, Descriptor: c.Descriptor})
	}

	outProds := make([]PublicProduct, 0, len(prods))
	for _, p := range prods {
		outProds = append(outProds, PublicProduct{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			CategoryID:  p.CategoryID,
			Photos:      p.Photos,
			Variants:    variantsByProduct[p.ID],
		})
	}

	resp := &CatalogResponse{
		Business: PublicBusiness{
			ID:                 biz.ID,
			Name:               biz.Name,
			Descriptor:         biz.Descriptor,
			Brand:              biz.Brand,
			CountryCode:        biz.CountryCode,
			Currency:           biz.Currency,
			StorefrontPublicID: biz.StorefrontPublicID,
			StorefrontEnabled:  biz.StorefrontEnabled,
			StorefrontTheme:    biz.StorefrontTheme,
			SupportEmail:       biz.SupportEmail,
			PhoneNumber:        biz.PhoneNumber,
			WhatsappNumber:     biz.WhatsappNumber,
			Address:            biz.Address,
			WebsiteURL:         biz.WebsiteURL,
			InstagramURL:       biz.InstagramURL,
			FacebookURL:        biz.FacebookURL,
			TikTokURL:          biz.TikTokURL,
			XURL:               biz.XURL,
			SnapchatURL:        biz.SnapchatURL,
		},
		Categories: outCats,
		Products:   outProds,
	}
	return resp, nil
}

func (s *Service) listAllProducts(ctx context.Context, biz *business.Business) ([]*inventory.Product, error) {
	const pageSize = 100
	const maxPages = 100

	all := make([]*inventory.Product, 0, pageSize)
	for page := 1; page <= maxPages; page++ {
		req := list.NewListRequest(page, pageSize, nil, "")
		items, err := s.inventory.ListProducts(ctx, nil, biz, req)
		if err != nil {
			return nil, err
		}
		all = append(all, items...)
		if len(items) < pageSize {
			break
		}
	}
	return all, nil
}

func (s *Service) listAllVariants(ctx context.Context, biz *business.Business) ([]*inventory.Variant, error) {
	const pageSize = 100
	const maxPages = 100

	all := make([]*inventory.Variant, 0, pageSize)
	for page := 1; page <= maxPages; page++ {
		req := list.NewListRequest(page, pageSize, nil, "")
		items, err := s.inventory.ListVariants(ctx, nil, biz, req)
		if err != nil {
			return nil, err
		}
		all = append(all, items...)
		if len(items) < pageSize {
			break
		}
	}
	return all, nil
}

type CreateOrderItem struct {
	VariantID      string `json:"variantId" binding:"required"`
	Quantity       int    `json:"quantity" binding:"required,gt=0"`
	SpecialRequest string `json:"specialRequest" binding:"omitempty,max=500"`
}

type CreateOrderCustomer struct {
	Email             string `json:"email" binding:"required,email"`
	Name              string `json:"name" binding:"required"`
	PhoneNumber       string `json:"phoneNumber" binding:"omitempty"`
	InstagramUsername string `json:"instagramUsername" binding:"omitempty"`
}

type CreateOrderShippingAddress struct {
	CountryCode string `json:"countryCode" binding:"required,len=2"`
	State       string `json:"state" binding:"required"`
	City        string `json:"city" binding:"required"`
	Street      string `json:"street" binding:"omitempty"`
	ZipCode     string `json:"zipCode" binding:"omitempty"`
	PhoneCode   string `json:"phoneCode" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
}

type CreateOrderRequest struct {
	Customer        CreateOrderCustomer        `json:"customer" binding:"required"`
	ShippingAddress CreateOrderShippingAddress `json:"shippingAddress" binding:"required"`
	Items           []CreateOrderItem          `json:"items" binding:"required,min=1,max=50,dive"`
}

type CreateOrderResponse struct {
	OrderID       string `json:"orderId"`
	OrderNumber   string `json:"orderNumber"`
	Status        string `json:"status"`
	PaymentStatus string `json:"paymentStatus"`
	Total         string `json:"total"`
	Currency      string `json:"currency"`
}

func (s *Service) CreatePendingOrder(ctx context.Context, storefrontPublicID, idempotencyKey string, requestBody []byte, clientIP string, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	idempotencyKey = strings.TrimSpace(idempotencyKey)
	if idempotencyKey == "" {
		return nil, ErrIdempotencyKeyRequired()
	}
	if len(idempotencyKey) > 128 {
		return nil, problem.BadRequest("Idempotency-Key is too long").With("header", "Idempotency-Key")
	}

	biz, err := s.business.GetBusinessByStorefrontPublicID(ctx, storefrontPublicID)
	if err != nil {
		return nil, ErrStorefrontNotFound(storefrontPublicID, err)
	}

	h := sha256.Sum256(requestBody)
	requestHash := Hash(h)

	// Fast-path idempotency: allow safe replays even if the caller is rate limited.
	if existing, err := s.storage.GetRequestByKey(ctx, biz.ID, idempotencyKey); err == nil {
		if existing.RequestHash != requestHash {
			return nil, ErrIdempotencyConflict()
		}
		if existing.OrderID == "" {
			return nil, ErrIdempotencyInProgress()
		}
		ord, oerr := s.orders.GetOrderByID(ctx, nil, biz, existing.OrderID)
		if oerr != nil {
			return nil, oerr
		}
		return &CreateOrderResponse{
			OrderID:       ord.ID,
			OrderNumber:   ord.OrderNumber,
			Status:        string(ord.Status),
			PaymentStatus: string(ord.PaymentStatus),
			Total:         ord.Total.String(),
			Currency:      ord.Currency,
		}, nil
	}

	ip := strings.TrimSpace(clientIP)
	if ip == "" {
		ip = "unknown"
	}
	throttleKey := fmt.Sprintf("storefront:%s:order:%s", biz.ID, ip)
	if !throttle.Allow(s.storage.Cache(), throttleKey, 1*time.Minute, 10, 1*time.Second) {
		return nil, problem.TooManyRequests("rate limit exceeded")
	}

	if _, err := mail.ParseAddress(req.Customer.Email); err != nil {
		return nil, problem.BadRequest("invalid email").With("field", "customer.email")
	}

	var out *CreateOrderResponse
	err = s.atomicProcessor.Exec(ctx, func(tctx context.Context) error {
		rec := &StorefrontRequest{
			BusinessID:     biz.ID,
			IdempotencyKey: idempotencyKey,
			RequestHash:    requestHash,
		}
		if err := s.storage.CreateRequest(tctx, rec); err != nil {
			if database.IsUniqueViolation(err) {
				existing, gerr := s.storage.GetRequestByKey(tctx, biz.ID, idempotencyKey)
				if gerr != nil {
					return gerr
				}
				if existing.RequestHash != requestHash {
					return ErrIdempotencyConflict()
				}
				if existing.OrderID == "" {
					return ErrIdempotencyInProgress()
				}
				ord, oerr := s.orders.GetOrderByID(tctx, nil, biz, existing.OrderID)
				if oerr != nil {
					return oerr
				}
				out = &CreateOrderResponse{
					OrderID:       ord.ID,
					OrderNumber:   ord.OrderNumber,
					Status:        string(ord.Status),
					PaymentStatus: string(ord.PaymentStatus),
					Total:         ord.Total.String(),
					Currency:      ord.Currency,
				}
				return nil
			}
			return err
		}

		cust, err := s.customer.UpsertCustomerByEmail(tctx, biz, &customer.UpsertCustomerByEmailInput{
			Email:             req.Customer.Email,
			Name:              req.Customer.Name,
			PhoneNumber:       req.Customer.PhoneNumber,
			PhoneCode:         req.ShippingAddress.PhoneCode,
			InstagramUsername: req.Customer.InstagramUsername,
		})
		if err != nil {
			return err
		}

		addr, err := s.customer.CreateCustomerAddress(tctx, nil, biz, cust.ID, &customer.CreateCustomerAddressRequest{
			CountryCode: req.ShippingAddress.CountryCode,
			State:       req.ShippingAddress.State,
			City:        req.ShippingAddress.City,
			PhoneCode:   req.ShippingAddress.PhoneCode,
			Phone:       req.ShippingAddress.PhoneNumber,
			Street:      req.ShippingAddress.Street,
			ZipCode:     req.ShippingAddress.ZipCode,
		})
		if err != nil {
			return err
		}

		qtyByVariant := map[string]int{}
		var noteLines []string
		for _, it := range req.Items {
			vid := strings.TrimSpace(it.VariantID)
			if vid == "" {
				return problem.BadRequest("variantId is required").With("field", "items.variantId")
			}
			qtyByVariant[vid] += it.Quantity
			if strings.TrimSpace(it.SpecialRequest) != "" {
				noteLines = append(noteLines, fmt.Sprintf("- %s: %s", vid, strings.TrimSpace(it.SpecialRequest)))
			}
		}

		note := ""
		if len(noteLines) > 0 {
			note = "Special requests:\n" + strings.Join(noteLines, "\n")
		}

		ord, err := s.orders.CreatePendingStorefrontOrder(tctx, biz, cust.ID, addr.ID, qtyByVariant, note)
		if err != nil {
			return err
		}

		rec.OrderID = ord.ID
		if err := s.storage.UpdateRequest(tctx, rec); err != nil {
			return err
		}

		out = &CreateOrderResponse{
			OrderID:       ord.ID,
			OrderNumber:   ord.OrderNumber,
			Status:        string(ord.Status),
			PaymentStatus: string(ord.PaymentStatus),
			Total:         ord.Total.String(),
			Currency:      ord.Currency,
		}
		return nil
	}, atomic.WithIsolationLevel(atomic.LevelSerializable), atomic.WithRetries(2))
	if err != nil {
		return nil, err
	}
	return out, nil
}
