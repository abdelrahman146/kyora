package handlers

import (
	"github.com/abdelrahman146/kyora/internal/domain/analytics"
	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

type dashboardHandler struct {
	storeDomain     *store.StoreDomain
	analyticsDomain *analytics.AnalyticsDomain
}

func AddDashboardRoutes(
	r *gin.RouterGroup,
	storeDomain *store.StoreDomain,
	analyticsDomain *analytics.AnalyticsDomain,
) {
	h := &dashboardHandler{
		storeDomain:     storeDomain,
		analyticsDomain: analyticsDomain,
	}
	h.registerRoutes(r)
}

func (h *dashboardHandler) registerRoutes(c *gin.RouterGroup) {
	c.GET("/dashboard", h.index)
}

func (h *dashboardHandler) index(c *gin.Context) {
	storeId := c.Param("storeId")
	_, err := h.analyticsDomain.Service.GenerateDashboardAnalytics(c.Request.Context(), storeId)
	if err != nil {
		webutils.Render(c, 500, pages.ErrorPage(500, "Failed to load dashboard analytics"))
		return
	}
	webutils.Render(c, 200, pages.NotImplemented("Dashboard"))
}
