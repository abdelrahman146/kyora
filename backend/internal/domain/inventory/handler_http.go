package inventory

import (
	"net/http"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/request"
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/abdelrahman146/kyora/internal/platform/types/list"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

// HttpHandler handles HTTP requests for inventory domain operations.
type HttpHandler struct {
	service *Service
}

// NewHttpHandler creates a new HTTP handler for inventory operations.
func NewHttpHandler(service *Service) *HttpHandler {
	return &HttpHandler{service: service}
}
func (h *HttpHandler) getBusinessForRequest(c *gin.Context, actor *account.User) (*business.Business, error) {
	_ = actor
	return business.BusinessFromContext(c)
}

type listInventoryQuery struct {
	Page        int      `form:"page" binding:"omitempty,min=1"`
	PageSize    int      `form:"pageSize" binding:"omitempty,min=1,max=100"`
	OrderBy     []string `form:"orderBy" binding:"omitempty"`
	SearchTerm  string   `form:"search" binding:"omitempty"`
	ProductID   string   `form:"productId" binding:"omitempty"`
	TopLimit    int      `form:"topLimit" binding:"omitempty,min=1,max=50"`
	DetailLimit int      `form:"limit" binding:"omitempty,min=1,max=50"`
}

// ListProducts returns a paginated list of products.
//
// @Summary      List products
// @Description  Returns a paginated list of products for the authenticated workspace/business. Each product includes its variants.
// @Tags         inventory
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        page query int false "Page number (default: 1)"
// @Param        pageSize query int false "Page size (default: 20, max: 100)"
// @Param        orderBy query []string false "Sort order (e.g., -name, createdAt)"
// @Param        search query string false "Search term for product name"
// @Success      200 {object} list.ListResponse[Product] "Products with their variants included"
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/inventory/products [get]
// @Security     BearerAuth
func (h *HttpHandler) ListProducts(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var query listInventoryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, problem.BadRequest("invalid query parameters").WithError(err))
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}
	if query.SearchTerm != "" {
		term, err := list.NormalizeSearchTerm(query.SearchTerm)
		if err != nil {
			response.Error(c, problem.BadRequest("invalid search term"))
			return
		}
		query.SearchTerm = term
	}
	listReq := list.NewListRequest(query.Page, query.PageSize, query.OrderBy, query.SearchTerm)
	items, err := h.service.ListProducts(c.Request.Context(), actor, biz, listReq)
	if err != nil {
		response.Error(c, err)
		return
	}
	total, err := h.service.CountProducts(c.Request.Context(), actor, biz)
	if err != nil {
		response.Error(c, err)
		return
	}
	hasMore := int64(query.Page*query.PageSize) < total
	response.SuccessJSON(c, http.StatusOK, list.NewListResponse(items, query.Page, query.PageSize, total, hasMore))
}

// GetProduct returns a product by ID including its variants.
//
// @Summary      Get product
// @Description  Returns a product by ID including its variants
// @Tags         inventory
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        productId path string true "Product ID"
// @Success      200 {object} Product
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/inventory/products/{productId} [get]
// @Security     BearerAuth
func (h *HttpHandler) GetProduct(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	id := c.Param("productId")
	product, err := h.service.GetProductByID(c.Request.Context(), actor, biz, id)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrProductNotFound(err).With("productId", id))
			return
		}
		response.Error(c, err)
		return
	}
	variants, err := h.service.GetProductVariants(c.Request.Context(), actor, biz, product.ID)
	if err != nil {
		response.Error(c, err)
		return
	}
	product.Variants = variants
	response.SuccessJSON(c, http.StatusOK, product)
}

// CreateProduct creates a new product.
//
// @Summary      Create product
// @Description  Creates a new product
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        body body CreateProductRequest true "Product"
// @Success      201 {object} Product
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/inventory/products [post]
// @Security     BearerAuth
func (h *HttpHandler) CreateProduct(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	var req CreateProductRequest
	if err := request.ValidBody(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	product, err := h.service.CreateProduct(c.Request.Context(), actor, biz, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusCreated, product)
}

// CreateProductWithVariants creates a product and its variants atomically.
//
// @Summary      Create product with variants
// @Description  Creates a product and its variants in a single request
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        body body CreateProductWithVariantsRequest true "Product with variants"
// @Success      201 {object} Product
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/inventory/products/with-variants [post]
// @Security     BearerAuth
func (h *HttpHandler) CreateProductWithVariants(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	var req CreateProductWithVariantsRequest
	if err := request.ValidBody(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	product, err := h.service.CreateProductWithVariants(c.Request.Context(), actor, biz, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusCreated, product)
}

// UpdateProduct updates an existing product.
//
// @Summary      Update product
// @Description  Updates an existing product
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        productId path string true "Product ID"
// @Param        body body UpdateProductRequest true "Updates"
// @Success      200 {object} Product
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/inventory/products/{productId} [patch]
// @Security     BearerAuth
func (h *HttpHandler) UpdateProduct(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	id := c.Param("productId")
	product, err := h.service.GetProductByID(c.Request.Context(), actor, biz, id)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrProductNotFound(err).With("productId", id))
			return
		}
		response.Error(c, err)
		return
	}
	var req UpdateProductRequest
	if err := request.ValidBody(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	if err := h.service.UpdateProduct(c.Request.Context(), actor, biz, product, &req); err != nil {
		response.Error(c, err)
		return
	}
	updated, err := h.service.GetProductByID(c.Request.Context(), actor, biz, id)
	if err != nil {
		response.Error(c, err)
		return
	}
	variants, err := h.service.GetProductVariants(c.Request.Context(), actor, biz, updated.ID)
	if err != nil {
		response.Error(c, err)
		return
	}
	updated.Variants = variants
	response.SuccessJSON(c, http.StatusOK, updated)
}

// DeleteProduct deletes a product.
//
// @Summary      Delete product
// @Description  Deletes a product
// @Tags         inventory
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        productId path string true "Product ID"
// @Success      204
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/inventory/products/{productId} [delete]
// @Security     BearerAuth
func (h *HttpHandler) DeleteProduct(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	id := c.Param("productId")
	if err := h.service.DeleteProduct(c.Request.Context(), actor, biz, id); err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrProductNotFound(err).With("productId", id))
			return
		}
		response.Error(c, err)
		return
	}
	response.SuccessEmpty(c, http.StatusNoContent)
}

// ListProductVariants returns a paginated list of variants for a product.
//
// @Summary      List product variants
// @Description  Returns a paginated list of variants for a given product
// @Tags         inventory
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        productId path string true "Product ID"
// @Param        page query int false "Page number (default: 1)"
// @Param        pageSize query int false "Page size (default: 20, max: 100)"
// @Param        orderBy query []string false "Sort order (e.g., -name, createdAt)"
// @Param        search query string false "Search term for variant name"
// @Success      200 {object} list.ListResponse[Variant]
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/inventory/products/{productId}/variants [get]
// @Security     BearerAuth
func (h *HttpHandler) ListProductVariants(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var query listInventoryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, problem.BadRequest("invalid query parameters").WithError(err))
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	productID := c.Param("productId")
	if _, err := h.service.GetProductByID(c.Request.Context(), actor, biz, productID); err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrProductNotFound(err).With("productId", productID))
			return
		}
		response.Error(c, err)
		return
	}
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}
	if query.SearchTerm != "" {
		term, err := list.NormalizeSearchTerm(query.SearchTerm)
		if err != nil {
			response.Error(c, problem.BadRequest("invalid search term"))
			return
		}
		query.SearchTerm = term
	}
	listReq := list.NewListRequest(query.Page, query.PageSize, query.OrderBy, query.SearchTerm)
	items, total, err := h.service.ListProductVariants(c.Request.Context(), actor, biz, productID, listReq)
	if err != nil {
		response.Error(c, err)
		return
	}
	hasMore := int64(query.Page*query.PageSize) < total
	response.SuccessJSON(c, http.StatusOK, list.NewListResponse(items, query.Page, query.PageSize, total, hasMore))
}

// ListVariants returns a paginated list of variants.
//
// @Summary      List variants
// @Description  Returns a paginated list of variants
// @Tags         inventory
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        page query int false "Page number (default: 1)"
// @Param        pageSize query int false "Page size (default: 20, max: 100)"
// @Param        orderBy query []string false "Sort order (e.g., -name, createdAt)"
// @Param        search query string false "Search term for variant name"
// @Success      200 {object} list.ListResponse[Variant]
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/inventory/variants [get]
// @Security     BearerAuth
func (h *HttpHandler) ListVariants(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var query listInventoryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, problem.BadRequest("invalid query parameters").WithError(err))
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}
	if query.SearchTerm != "" {
		term, err := list.NormalizeSearchTerm(query.SearchTerm)
		if err != nil {
			response.Error(c, problem.BadRequest("invalid search term"))
			return
		}
		query.SearchTerm = term
	}
	listReq := list.NewListRequest(query.Page, query.PageSize, query.OrderBy, query.SearchTerm)
	items, err := h.service.ListVariants(c.Request.Context(), actor, biz, listReq)
	if err != nil {
		response.Error(c, err)
		return
	}
	total, err := h.service.CountVariants(c.Request.Context(), actor, biz)
	if err != nil {
		response.Error(c, err)
		return
	}
	hasMore := int64(query.Page*query.PageSize) < total
	response.SuccessJSON(c, http.StatusOK, list.NewListResponse(items, query.Page, query.PageSize, total, hasMore))
}

// GetVariant returns a variant by ID.
//
// @Summary      Get variant
// @Description  Returns a variant by ID
// @Tags         inventory
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        variantId path string true "Variant ID"
// @Success      200 {object} Variant
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/inventory/variants/{variantId} [get]
// @Security     BearerAuth
func (h *HttpHandler) GetVariant(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	id := c.Param("variantId")
	variant, err := h.service.GetVariantByID(c.Request.Context(), actor, biz, id)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrVariantNotFound(err).With("variantId", id))
			return
		}
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, variant)
}

// CreateVariant creates a new variant.
//
// @Summary      Create variant
// @Description  Creates a new variant under a product
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        body body CreateVariantRequest true "Variant"
// @Success      201 {object} Variant
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/inventory/variants [post]
// @Security     BearerAuth
func (h *HttpHandler) CreateVariant(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	var req CreateVariantRequest
	if err := request.ValidBody(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	variant, err := h.service.CreateVariant(c.Request.Context(), actor, biz, &req)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrProductNotFound(err).With("productId", req.ProductID))
			return
		}
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusCreated, variant)
}

// UpdateVariant updates a variant.
//
// @Summary      Update variant
// @Description  Updates a variant
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        variantId path string true "Variant ID"
// @Param        body body UpdateVariantRequest true "Updates"
// @Success      200 {object} Variant
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/inventory/variants/{variantId} [patch]
// @Security     BearerAuth
func (h *HttpHandler) UpdateVariant(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	id := c.Param("variantId")
	var req UpdateVariantRequest
	if err := request.ValidBody(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	if err := h.service.UpdateVariant(c.Request.Context(), actor, biz, id, &req); err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrVariantNotFound(err).With("variantId", id))
			return
		}
		response.Error(c, err)
		return
	}
	updated, err := h.service.GetVariantByID(c.Request.Context(), actor, biz, id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, updated)
}

// DeleteVariant deletes a variant.
//
// @Summary      Delete variant
// @Description  Deletes a variant
// @Tags         inventory
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        variantId path string true "Variant ID"
// @Success      204
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/inventory/variants/{variantId} [delete]
// @Security     BearerAuth
func (h *HttpHandler) DeleteVariant(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	id := c.Param("variantId")
	if err := h.service.DeleteVariant(c.Request.Context(), actor, biz, id); err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrVariantNotFound(err).With("variantId", id))
			return
		}
		response.Error(c, err)
		return
	}
	response.SuccessEmpty(c, http.StatusNoContent)
}

// ListCategories returns all categories.
//
// @Summary      List categories
// @Description  Returns all categories
// @Tags         inventory
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Success      200 {array} Category
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/inventory/categories [get]
// @Security     BearerAuth
func (h *HttpHandler) ListCategories(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	items, err := h.service.ListCategories(c.Request.Context(), actor, biz)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, items)
}

// GetCategory returns a category by ID.
//
// @Summary      Get category
// @Description  Returns a category by ID
// @Tags         inventory
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        categoryId path string true "Category ID"
// @Success      200 {object} Category
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/inventory/categories/{categoryId} [get]
// @Security     BearerAuth
func (h *HttpHandler) GetCategory(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	id := c.Param("categoryId")
	cat, err := h.service.GetCategoryByID(c.Request.Context(), actor, biz, id)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrCategoryNotFound(err).With("categoryId", id))
			return
		}
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, cat)
}

// CreateCategory creates a category.
//
// @Summary      Create category
// @Description  Creates a category
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        body body CreateCategoryRequest true "Category"
// @Success      201 {object} Category
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/inventory/categories [post]
// @Security     BearerAuth
func (h *HttpHandler) CreateCategory(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	var req CreateCategoryRequest
	if err := request.ValidBody(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	cat, err := h.service.CreateCategory(c.Request.Context(), actor, biz, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusCreated, cat)
}

// UpdateCategory updates a category.
//
// @Summary      Update category
// @Description  Updates a category
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        categoryId path string true "Category ID"
// @Param        body body UpdateCategoryRequest true "Updates"
// @Success      200 {object} Category
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/inventory/categories/{categoryId} [patch]
// @Security     BearerAuth
func (h *HttpHandler) UpdateCategory(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	id := c.Param("categoryId")
	cat, err := h.service.GetCategoryByID(c.Request.Context(), actor, biz, id)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrCategoryNotFound(err).With("categoryId", id))
			return
		}
		response.Error(c, err)
		return
	}
	var req UpdateCategoryRequest
	if err := request.ValidBody(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	if err := h.service.UpdateCategory(c.Request.Context(), actor, biz, cat, &req); err != nil {
		response.Error(c, err)
		return
	}
	updated, err := h.service.GetCategoryByID(c.Request.Context(), actor, biz, id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, updated)
}

// DeleteCategory deletes a category.
//
// @Summary      Delete category
// @Description  Deletes a category
// @Tags         inventory
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        categoryId path string true "Category ID"
// @Success      204
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/inventory/categories/{categoryId} [delete]
// @Security     BearerAuth
func (h *HttpHandler) DeleteCategory(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	id := c.Param("categoryId")
	if err := h.service.DeleteCategory(c.Request.Context(), actor, biz, id); err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrCategoryNotFound(err).With("categoryId", id))
			return
		}
		response.Error(c, err)
		return
	}
	response.SuccessEmpty(c, http.StatusNoContent)
}

type inventorySummaryResponse struct {
	ProductsCount           int64           `json:"productsCount"`
	VariantsCount           int64           `json:"variantsCount"`
	CategoriesCount         int64           `json:"categoriesCount"`
	LowStockVariantsCount   int64           `json:"lowStockVariantsCount"`
	OutOfStockVariantsCount int64           `json:"outOfStockVariantsCount"`
	TotalStockUnits         int64           `json:"totalStockUnits"`
	InventoryValue          decimal.Decimal `json:"inventoryValue"`
	TopProducts             []*Product      `json:"topProductsByInventoryValue"`
}

// GetInventorySummary returns inventory summary metrics.
//
// @Summary      Inventory summary
// @Description  Returns inventory metrics for the current business
// @Tags         inventory
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        topLimit query int false "Top products limit (default: 5, max: 50)"
// @Success      200 {object} inventorySummaryResponse
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/inventory/summary [get]
// @Security     BearerAuth
func (h *HttpHandler) GetInventorySummary(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var query listInventoryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, problem.BadRequest("invalid query parameters").WithError(err))
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	productsCount, err := h.service.CountProducts(c.Request.Context(), actor, biz)
	if err != nil {
		response.Error(c, err)
		return
	}
	variantsCount, err := h.service.CountVariants(c.Request.Context(), actor, biz)
	if err != nil {
		response.Error(c, err)
		return
	}
	categoriesCount, err := h.service.CountCategories(c.Request.Context(), actor, biz)
	if err != nil {
		response.Error(c, err)
		return
	}
	lowStockCount, err := h.service.CountLowStockVariants(c.Request.Context(), actor, biz)
	if err != nil {
		response.Error(c, err)
		return
	}
	outOfStockCount, err := h.service.CountOutOfStockVariants(c.Request.Context(), actor, biz)
	if err != nil {
		response.Error(c, err)
		return
	}
	totalUnits, err := h.service.SumStockQuantity(c.Request.Context(), actor, biz)
	if err != nil {
		response.Error(c, err)
		return
	}
	invValue, err := h.service.SumInventoryValue(c.Request.Context(), actor, biz)
	if err != nil {
		response.Error(c, err)
		return
	}
	limit := query.TopLimit
	if limit == 0 {
		limit = 5
	}
	topProducts, err := h.service.ComputeTopProductsByInventoryValue(c.Request.Context(), actor, biz, limit)
	if err != nil {
		response.Error(c, err)
		return
	}
	resp := inventorySummaryResponse{
		ProductsCount:           productsCount,
		VariantsCount:           variantsCount,
		CategoriesCount:         categoriesCount,
		LowStockVariantsCount:   lowStockCount,
		OutOfStockVariantsCount: outOfStockCount,
		TotalStockUnits:         totalUnits,
		InventoryValue:          invValue,
		TopProducts:             topProducts,
	}
	response.SuccessJSON(c, http.StatusOK, resp)
}

// GetTopProductsByInventoryValue returns top products including the computed inventory value.
//
// @Summary      Top products by inventory value
// @Description  Returns the top products by inventory value (cost * stock)
// @Tags         inventory
// @Produce      json
// @Param        businessDescriptor path string true "Business descriptor"
// @Param        limit query int false "Top limit (default: 5, max: 50)"
// @Success      200 {array} TopProductByInventoryValue
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/businesses/{businessDescriptor}/inventory/top-products [get]
// @Security     BearerAuth
func (h *HttpHandler) GetTopProductsByInventoryValue(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var query listInventoryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, problem.BadRequest("invalid query parameters").WithError(err))
		return
	}
	biz, err := h.getBusinessForRequest(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}
	limit := query.DetailLimit
	if limit == 0 {
		limit = 5
	}
	items, err := h.service.ComputeTopProductsByInventoryValueDetailed(c.Request.Context(), actor, biz, limit)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, items)
}
