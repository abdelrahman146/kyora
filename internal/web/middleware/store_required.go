package middleware

import (
	"fmt"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/store"

	"github.com/gin-gonic/gin"
)

func StoreRequired(storeService *store.StoreService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exist := c.Get(UserKey)
		if !exist {
			c.Redirect(302, loginPath)
			c.Abort()
			return
		}
		u := user.(*account.User)
		storeID := c.Param("storeId")
		userOrgStores, err := storeService.ListOrganizationStores(c.Request.Context(), u.OrganizationID)
		if err != nil || len(userOrgStores) == 0 {
			c.Redirect(302, "/onboarding")
			c.Abort()
			return
		}
		if storeID == "" && len(userOrgStores) > 0 {
			c.Redirect(302, fmt.Sprintf("/%s/dashboard", userOrgStores[0].ID))
			c.Abort()
			return
		}
		var storeFound *store.Store
		for _, s := range userOrgStores {
			if s.ID == storeID {
				storeFound = s
				break
			}
		}
		if storeFound == nil {
			c.String(404, "Store not found")
			c.Abort()
			return
		}
		c.Set(StoresListKey, userOrgStores)
		c.Set(StoreKey, storeFound)
		c.Next()
	}
}
