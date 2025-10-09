package loginPage

import (
	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/abdelrahman146/kyora/internal/web/middleware"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	authService *account.AuthenticationService
}

func NewLoginHandler(authService *account.AuthenticationService) *LoginHandler {
	return &LoginHandler{authService}
}

func (h *LoginHandler) RegisterRoutes(r gin.IRoutes) {
	r.GET("/login", middleware.GuestRequiredMiddleware, h.Index)
	r.POST("/login", middleware.GuestRequiredMiddleware, h.Login)
}

func (h *LoginHandler) Index(c *gin.Context) {
	info := webcontext.PageInfo{
		Title:       "Login",
		Description: "Login page ",
		Keywords:    "login, zard",
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, Login())
}

func (h *LoginHandler) Login(c *gin.Context) {
	var form struct {
		Email    string `form:"email" binding:"required,email"`
		Password string `form:"password" binding:"required"`
	}
	if err := c.ShouldBind(&form); err != nil {
		webutils.Render(c, 400, LoginFailed("Invalid form data"))
		return
	}
	_, jwtToken, err := h.authService.Authenticate(c.Request.Context(), form.Email, form.Password)
	if err != nil {
		webutils.Render(c, 401, LoginFailed("Invalid email or password"))
		return
	}
	utils.JWT.SetJwtCookie(c, jwtToken)
	webutils.Redirect(c, "/")
}
