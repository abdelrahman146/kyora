package handlers

import (
	"net/http"

	"github.com/abdelrahman146/kyora/internal/domain/owner"
	"github.com/gin-gonic/gin"
)

type ownerHandler struct {
	ownerDomain *owner.OwnerDomain
}

func AddOwnerRoutes(r *gin.Engine, ownerDomain *owner.OwnerDomain) {
	h := &ownerHandler{ownerDomain: ownerDomain}
	h.registerRoutes(r)
}

func (h *ownerHandler) registerRoutes(r *gin.Engine) {
	r.GET("/", h.index)
}

func (h *ownerHandler) index(c *gin.Context) {
	c.String(http.StatusOK, "not implemented")
}
