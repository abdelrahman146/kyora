package handlers

import (
	"github.com/abdelrahman146/kyora/internal/domain/customer"
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

type customerHandler struct {
	customerDomain *customer.CustomerDomain
}

func AddCustomerRoutes(r *gin.RouterGroup, customerDomain *customer.CustomerDomain) {
	h := &customerHandler{customerDomain: customerDomain}
	h.registerRoutes(r)
}

func (h *customerHandler) registerRoutes(c *gin.RouterGroup) {
	r := c.Group("/customers")
	{
		r.GET("/", h.index)
		r.POST("/", h.create)
		r.GET("/:id", h.show)      // View customer details
		r.PUT("/:id", h.update)    // Update customer details
		r.DELETE("/:id", h.delete) // Delete customer
	}
}

func (h *customerHandler) index(c *gin.Context) {
	storeId := c.Param("storeId")
	listReq := webutils.GetPaginationParams(c)
	_, err := h.customerDomain.CustomerService.ListCustomers(c.Request.Context(), storeId, listReq)
	if err != nil {
		utils.Log.FromContext(c.Request.Context()).Error("failed to list customers", "error", err, "storeId", storeId)
		webutils.Render(c, 500, pages.ErrorPage(500, "Failed to load customers"))
		return
	}
	_, err = h.customerDomain.CustomerService.CountCustomers(c.Request.Context(), storeId)
	if err != nil {
		utils.Log.FromContext(c.Request.Context()).Error("failed to count customers", "error", err, "storeId", storeId)
		webutils.Render(c, 500, pages.ErrorPage(500, "Failed to load customers"))
		return
	}
	webutils.Render(c, 200, pages.NotImplemented("Customers List"))
}

func (h *customerHandler) create(c *gin.Context) {
	_ = c.Param("storeId")
	// receive form data and validate customer.CreateCustomerRequest
	c.String(200, "not implemented")
}

func (h *customerHandler) show(c *gin.Context) {
	storeId := c.Param("storeId")
	id := c.Param("id")
	_, err := h.customerDomain.CustomerService.GetCustomerByID(c.Request.Context(), storeId, id)
	if err != nil {
		utils.Log.FromContext(c.Request.Context()).Error("failed to get customer", "error", err, "storeId", storeId, "customerId", id)
		webutils.Render(c, 500, pages.ErrorPage(500, "Failed to load customer"))
		return
	}
	webutils.Render(c, 200, pages.NotImplemented("Customer Details"))
}

func (h *customerHandler) update(c *gin.Context) {
	_ = c.Param("storeId")
	_ = c.Param("id")
	// receive form data and validate customer.UpdateCustomerRequest
	c.String(200, "not implemented")
}

func (h *customerHandler) delete(c *gin.Context) {
	_ = c.Param("storeId")
	_ = c.Param("id")
	// perform delete operation
	c.String(200, "not implemented")
}
