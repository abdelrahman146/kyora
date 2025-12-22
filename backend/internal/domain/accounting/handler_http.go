package accounting

import (
	"net/http"
	"time"

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

// HttpHandler handles HTTP requests for accounting domain operations
type HttpHandler struct {
	service         *Service
	businessService *business.Service
}

// NewHttpHandler creates a new HTTP handler for accounting operations
func NewHttpHandler(service *Service, businessService *business.Service) *HttpHandler {
	return &HttpHandler{
		service:         service,
		businessService: businessService,
	}
}

// getBusinessForWorkspace is a helper that gets the first business for a workspace
// In future, this might need to support multiple businesses per workspace
func (h *HttpHandler) getBusinessForWorkspace(c *gin.Context, actor *account.User) (*business.Business, error) {
	businesses, err := h.businessService.ListBusinesses(c.Request.Context(), actor)
	if err != nil {
		return nil, err
	}
	if len(businesses) == 0 {
		return nil, problem.NotFound("no business found for this workspace")
	}
	return businesses[0], nil
}

// Asset endpoints

type listAssetsQuery struct {
	Page     int      `form:"page" binding:"omitempty,min=1"`
	PageSize int      `form:"pageSize" binding:"omitempty,min=1,max=100"`
	OrderBy  []string `form:"orderBy" binding:"omitempty"`
}

// ListAssets returns a paginated list of assets for the workspace
//
// @Summary      List assets
// @Description  Returns a paginated list of all assets for the authenticated workspace
// @Tags         accounting
// @Produce      json
// @Param        page query int false "Page number (default: 1)"
// @Param        pageSize query int false "Page size (default: 20, max: 100)"
// @Param        orderBy query []string false "Sort order (e.g., -value, name)"
// @Success      200 {object} list.ListResponse[Asset]
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/assets [get]
// @Security     BearerAuth
func (h *HttpHandler) ListAssets(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var query listAssetsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, problem.BadRequest("invalid query parameters").WithError(err))
		return
	}

	// Get business for the workspace
	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	// Set defaults
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}

	listReq := list.NewListRequest(query.Page, query.PageSize, query.OrderBy, "")
	assets, err := h.service.ListAssets(c.Request.Context(), actor, biz, listReq)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	totalCount, err := h.service.CountAssets(c.Request.Context(), actor, biz)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	hasMore := int64(query.Page*query.PageSize) < totalCount
	listResp := list.NewListResponse(assets, query.Page, query.PageSize, totalCount, hasMore)
	response.SuccessJSON(c, http.StatusOK, listResp)
}

// GetAsset returns a specific asset by ID
//
// @Summary      Get asset
// @Description  Returns a specific asset by ID for the authenticated workspace
// @Tags         accounting
// @Produce      json
// @Param        assetId path string true "Asset ID"
// @Success      200 {object} Asset
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/assets/{assetId} [get]
// @Security     BearerAuth
func (h *HttpHandler) GetAsset(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	assetID := c.Param("assetId")
	if assetID == "" {
		response.Error(c, problem.BadRequest("assetId is required"))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	asset, err := h.service.GetAssetByID(c.Request.Context(), actor, biz, assetID)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrAssetNotFound(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessJSON(c, http.StatusOK, asset)
}

// CreateAsset creates a new asset
//
// @Summary      Create asset
// @Description  Creates a new asset for the authenticated workspace
// @Tags         accounting
// @Accept       json
// @Produce      json
// @Param        request body CreateAssetRequest true "Asset data"
// @Success      201 {object} Asset
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/assets [post]
// @Security     BearerAuth
func (h *HttpHandler) CreateAsset(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req CreateAssetRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	asset, err := h.service.CreateAsset(c.Request.Context(), actor, biz, &req)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessJSON(c, http.StatusCreated, asset)
}

// UpdateAsset updates an existing asset
//
// @Summary      Update asset
// @Description  Updates an existing asset for the authenticated workspace
// @Tags         accounting
// @Accept       json
// @Produce      json
// @Param        assetId path string true "Asset ID"
// @Param        request body UpdateAssetRequest true "Asset update data"
// @Success      200 {object} Asset
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/assets/{assetId} [patch]
// @Security     BearerAuth
func (h *HttpHandler) UpdateAsset(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	assetID := c.Param("assetId")
	if assetID == "" {
		response.Error(c, problem.BadRequest("assetId is required"))
		return
	}

	var req UpdateAssetRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	asset, err := h.service.UpdateAsset(c.Request.Context(), actor, biz, assetID, &req)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrAssetNotFound(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessJSON(c, http.StatusOK, asset)
}

// DeleteAsset deletes an asset
//
// @Summary      Delete asset
// @Description  Deletes an asset for the authenticated workspace
// @Tags         accounting
// @Produce      json
// @Param        assetId path string true "Asset ID"
// @Success      204
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/assets/{assetId} [delete]
// @Security     BearerAuth
func (h *HttpHandler) DeleteAsset(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	assetID := c.Param("assetId")
	if assetID == "" {
		response.Error(c, problem.BadRequest("assetId is required"))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	err = h.service.DeleteAsset(c.Request.Context(), actor, biz, assetID)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrAssetNotFound(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessEmpty(c, http.StatusNoContent)
}

// Investment endpoints

// ListInvestments returns a paginated list of investments for the workspace
//
// @Summary      List investments
// @Description  Returns a paginated list of all investments for the authenticated workspace
// @Tags         accounting
// @Produce      json
// @Param        page query int false "Page number (default: 1)"
// @Param        pageSize query int false "Page size (default: 20, max: 100)"
// @Param        orderBy query []string false "Sort order (e.g., -amount, investedAt)"
// @Success      200 {object} list.ListResponse[Investment]
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/investments [get]
// @Security     BearerAuth
func (h *HttpHandler) ListInvestments(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var query listAssetsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, problem.BadRequest("invalid query parameters").WithError(err))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
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

	listReq := list.NewListRequest(query.Page, query.PageSize, query.OrderBy, "")
	investments, err := h.service.ListInvestments(c.Request.Context(), actor, biz, listReq)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	totalCount, err := h.service.CountInvestments(c.Request.Context(), actor, biz)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	listResp := list.NewListResponse(investments, query.Page, query.PageSize, totalCount, (int64(query.Page*query.PageSize) < totalCount))
	response.SuccessJSON(c, http.StatusOK, listResp)
}

// GetInvestment returns a specific investment by ID
//
// @Summary      Get investment
// @Description  Returns a specific investment by ID for the authenticated workspace
// @Tags         accounting
// @Produce      json
// @Param        investmentId path string true "Investment ID"
// @Success      200 {object} Investment
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/investments/{investmentId} [get]
// @Security     BearerAuth
func (h *HttpHandler) GetInvestment(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	investmentID := c.Param("investmentId")
	if investmentID == "" {
		response.Error(c, problem.BadRequest("investmentId is required"))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	investment, err := h.service.GetInvestmentByID(c.Request.Context(), actor, biz, investmentID)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrInvestmentNotFound(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessJSON(c, http.StatusOK, investment)
}

// CreateInvestment creates a new investment
//
// @Summary      Create investment
// @Description  Creates a new investment for the authenticated workspace
// @Tags         accounting
// @Accept       json
// @Produce      json
// @Param        request body CreateInvestmentRequest true "Investment data"
// @Success      201 {object} Investment
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/investments [post]
// @Security     BearerAuth
func (h *HttpHandler) CreateInvestment(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req CreateInvestmentRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	investment, err := h.service.CreateInvestment(c.Request.Context(), actor, biz, &req)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessJSON(c, http.StatusCreated, investment)
}

// UpdateInvestment updates an existing investment
//
// @Summary      Update investment
// @Description  Updates an existing investment for the authenticated workspace
// @Tags         accounting
// @Accept       json
// @Produce      json
// @Param        investmentId path string true "Investment ID"
// @Param        request body UpdateInvestmentRequest true "Investment update data"
// @Success      200 {object} Investment
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/investments/{investmentId} [patch]
// @Security     BearerAuth
func (h *HttpHandler) UpdateInvestment(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	investmentID := c.Param("investmentId")
	if investmentID == "" {
		response.Error(c, problem.BadRequest("investmentId is required"))
		return
	}

	var req UpdateInvestmentRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	investment, err := h.service.UpdateInvestment(c.Request.Context(), actor, biz, investmentID, &req)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrInvestmentNotFound(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessJSON(c, http.StatusOK, investment)
}

// DeleteInvestment deletes an investment
//
// @Summary      Delete investment
// @Description  Deletes an investment for the authenticated workspace
// @Tags         accounting
// @Produce      json
// @Param        investmentId path string true "Investment ID"
// @Success      204
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/investments/{investmentId} [delete]
// @Security     BearerAuth
func (h *HttpHandler) DeleteInvestment(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	investmentID := c.Param("investmentId")
	if investmentID == "" {
		response.Error(c, problem.BadRequest("investmentId is required"))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	err = h.service.DeleteInvestment(c.Request.Context(), actor, biz, investmentID)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrInvestmentNotFound(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessEmpty(c, http.StatusNoContent)
}

// Withdrawal endpoints

// ListWithdrawals returns a paginated list of withdrawals for the workspace
//
// @Summary      List withdrawals
// @Description  Returns a paginated list of all withdrawals for the authenticated workspace
// @Tags         accounting
// @Produce      json
// @Param        page query int false "Page number (default: 1)"
// @Param        pageSize query int false "Page size (default: 20, max: 100)"
// @Param        orderBy query []string false "Sort order (e.g., -amount, withdrawnAt)"
// @Success      200 {object} list.ListResponse[Withdrawal]
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/withdrawals [get]
// @Security     BearerAuth
func (h *HttpHandler) ListWithdrawals(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var query listAssetsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, problem.BadRequest("invalid query parameters").WithError(err))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
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

	listReq := list.NewListRequest(query.Page, query.PageSize, query.OrderBy, "")
	withdrawals, err := h.service.ListWithdrawals(c.Request.Context(), actor, biz, listReq)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	totalCount, err := h.service.CountWithdrawals(c.Request.Context(), actor, biz)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	listResp := list.NewListResponse(withdrawals, query.Page, query.PageSize, totalCount, (int64(query.Page*query.PageSize) < totalCount))
	response.SuccessJSON(c, http.StatusOK, listResp)
}

// GetWithdrawal returns a specific withdrawal by ID
//
// @Summary      Get withdrawal
// @Description  Returns a specific withdrawal by ID for the authenticated workspace
// @Tags         accounting
// @Produce      json
// @Param        withdrawalId path string true "Withdrawal ID"
// @Success      200 {object} Withdrawal
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/withdrawals/{withdrawalId} [get]
// @Security     BearerAuth
func (h *HttpHandler) GetWithdrawal(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	withdrawalID := c.Param("withdrawalId")
	if withdrawalID == "" {
		response.Error(c, problem.BadRequest("withdrawalId is required"))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	withdrawal, err := h.service.GetWithdrawalByID(c.Request.Context(), actor, biz, withdrawalID)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrWithdrawalNotFound(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessJSON(c, http.StatusOK, withdrawal)
}

// CreateWithdrawal creates a new withdrawal
//
// @Summary      Create withdrawal
// @Description  Creates a new withdrawal for the authenticated workspace
// @Tags         accounting
// @Accept       json
// @Produce      json
// @Param        request body CreateWithdrawalRequest true "Withdrawal data"
// @Success      201 {object} Withdrawal
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/withdrawals [post]
// @Security     BearerAuth
func (h *HttpHandler) CreateWithdrawal(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req CreateWithdrawalRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	withdrawal, err := h.service.CreateWithdrawal(c.Request.Context(), actor, biz, &req)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessJSON(c, http.StatusCreated, withdrawal)
}

// UpdateWithdrawal updates an existing withdrawal
//
// @Summary      Update withdrawal
// @Description  Updates an existing withdrawal for the authenticated workspace
// @Tags         accounting
// @Accept       json
// @Produce      json
// @Param        withdrawalId path string true "Withdrawal ID"
// @Param        request body UpdateWithdrawalRequest true "Withdrawal update data"
// @Success      200 {object} Withdrawal
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/withdrawals/{withdrawalId} [patch]
// @Security     BearerAuth
func (h *HttpHandler) UpdateWithdrawal(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	withdrawalID := c.Param("withdrawalId")
	if withdrawalID == "" {
		response.Error(c, problem.BadRequest("withdrawalId is required"))
		return
	}

	var req UpdateWithdrawalRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	withdrawal, err := h.service.UpdateWithdrawal(c.Request.Context(), actor, biz, withdrawalID, &req)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrWithdrawalNotFound(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessJSON(c, http.StatusOK, withdrawal)
}

// DeleteWithdrawal deletes a withdrawal
//
// @Summary      Delete withdrawal
// @Description  Deletes a withdrawal for the authenticated workspace
// @Tags         accounting
// @Produce      json
// @Param        withdrawalId path string true "Withdrawal ID"
// @Success      204
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/withdrawals/{withdrawalId} [delete]
// @Security     BearerAuth
func (h *HttpHandler) DeleteWithdrawal(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	withdrawalID := c.Param("withdrawalId")
	if withdrawalID == "" {
		response.Error(c, problem.BadRequest("withdrawalId is required"))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	err = h.service.DeleteWithdrawal(c.Request.Context(), actor, biz, withdrawalID)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrWithdrawalNotFound(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessEmpty(c, http.StatusNoContent)
}

// Expense endpoints

// ListExpenses returns a paginated list of expenses for the workspace
//
// @Summary      List expenses
// @Description  Returns a paginated list of all expenses for the authenticated workspace
// @Tags         accounting
// @Produce      json
// @Param        page query int false "Page number (default: 1)"
// @Param        pageSize query int false "Page size (default: 20, max: 100)"
// @Param        orderBy query []string false "Sort order (e.g., -amount, occurredOn)"
// @Success      200 {object} list.ListResponse[Expense]
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/expenses [get]
// @Security     BearerAuth
func (h *HttpHandler) ListExpenses(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var query listAssetsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, problem.BadRequest("invalid query parameters").WithError(err))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
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

	listReq := list.NewListRequest(query.Page, query.PageSize, query.OrderBy, "")
	expenses, err := h.service.ListExpenses(c.Request.Context(), actor, biz, listReq)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	totalCount, err := h.service.CountExpenses(c.Request.Context(), actor, biz)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	listResp := list.NewListResponse(expenses, query.Page, query.PageSize, totalCount, (int64(query.Page*query.PageSize) < totalCount))
	response.SuccessJSON(c, http.StatusOK, listResp)
}

// GetExpense returns a specific expense by ID
//
// @Summary      Get expense
// @Description  Returns a specific expense by ID for the authenticated workspace
// @Tags         accounting
// @Produce      json
// @Param        expenseId path string true "Expense ID"
// @Success      200 {object} Expense
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/expenses/{expenseId} [get]
// @Security     BearerAuth
func (h *HttpHandler) GetExpense(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	expenseID := c.Param("expenseId")
	if expenseID == "" {
		response.Error(c, problem.BadRequest("expenseId is required"))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	expense, err := h.service.GetExpenseByID(c.Request.Context(), actor, biz, expenseID)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrExpenseNotFound(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessJSON(c, http.StatusOK, expense)
}

// CreateExpense creates a new expense
//
// @Summary      Create expense
// @Description  Creates a new expense for the authenticated workspace
// @Tags         accounting
// @Accept       json
// @Produce      json
// @Param        request body CreateExpenseRequest true "Expense data"
// @Success      201 {object} Expense
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/expenses [post]
// @Security     BearerAuth
func (h *HttpHandler) CreateExpense(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req CreateExpenseRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	expense, err := h.service.CreateExpense(c.Request.Context(), actor, biz, &req)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessJSON(c, http.StatusCreated, expense)
}

// UpdateExpense updates an existing expense
//
// @Summary      Update expense
// @Description  Updates an existing expense for the authenticated workspace
// @Tags         accounting
// @Accept       json
// @Produce      json
// @Param        expenseId path string true "Expense ID"
// @Param        request body UpdateExpenseRequest true "Expense update data"
// @Success      200 {object} Expense
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/expenses/{expenseId} [patch]
// @Security     BearerAuth
func (h *HttpHandler) UpdateExpense(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	expenseID := c.Param("expenseId")
	if expenseID == "" {
		response.Error(c, problem.BadRequest("expenseId is required"))
		return
	}

	var req UpdateExpenseRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	expense, err := h.service.UpdateExpense(c.Request.Context(), actor, biz, expenseID, &req)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrExpenseNotFound(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessJSON(c, http.StatusOK, expense)
}

// DeleteExpense deletes an expense
//
// @Summary      Delete expense
// @Description  Deletes an expense for the authenticated workspace
// @Tags         accounting
// @Produce      json
// @Param        expenseId path string true "Expense ID"
// @Success      204
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/expenses/{expenseId} [delete]
// @Security     BearerAuth
func (h *HttpHandler) DeleteExpense(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	expenseID := c.Param("expenseId")
	if expenseID == "" {
		response.Error(c, problem.BadRequest("expenseId is required"))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	err = h.service.DeleteExpense(c.Request.Context(), actor, biz, expenseID)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrExpenseNotFound(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessEmpty(c, http.StatusNoContent)
}

// Recurring Expense endpoints

// ListRecurringExpenses returns a paginated list of recurring expenses for the workspace
//
// @Summary      List recurring expenses
// @Description  Returns a paginated list of all recurring expenses for the authenticated workspace
// @Tags         accounting
// @Produce      json
// @Param        page query int false "Page number (default: 1)"
// @Param        pageSize query int false "Page size (default: 20, max: 100)"
// @Param        orderBy query []string false "Sort order (e.g., -amount, nextRecurringDate)"
// @Success      200 {object} list.ListResponse[RecurringExpense]
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/recurring-expenses [get]
// @Security     BearerAuth
func (h *HttpHandler) ListRecurringExpenses(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var query listAssetsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, problem.BadRequest("invalid query parameters").WithError(err))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
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

	listReq := list.NewListRequest(query.Page, query.PageSize, query.OrderBy, "")
	recurringExpenses, err := h.service.ListRecurringExpenses(c.Request.Context(), actor, biz, listReq)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	totalCount, err := h.service.CountRecurringExpenses(c.Request.Context(), actor, biz)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	listResp := list.NewListResponse(recurringExpenses, query.Page, query.PageSize, totalCount, (int64(query.Page*query.PageSize) < totalCount))
	response.SuccessJSON(c, http.StatusOK, listResp)
}

// GetRecurringExpense returns a specific recurring expense by ID
//
// @Summary      Get recurring expense
// @Description  Returns a specific recurring expense by ID for the authenticated workspace
// @Tags         accounting
// @Produce      json
// @Param        recurringExpenseId path string true "Recurring Expense ID"
// @Success      200 {object} RecurringExpense
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/recurring-expenses/{recurringExpenseId} [get]
// @Security     BearerAuth
func (h *HttpHandler) GetRecurringExpense(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	recurringExpenseID := c.Param("recurringExpenseId")
	if recurringExpenseID == "" {
		response.Error(c, problem.BadRequest("recurringExpenseId is required"))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	recurringExpense, err := h.service.GetRecurringExpenseByID(c.Request.Context(), actor, biz, recurringExpenseID)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrRecurringExpenseNotFound(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessJSON(c, http.StatusOK, recurringExpense)
}

// CreateRecurringExpense creates a new recurring expense
//
// @Summary      Create recurring expense
// @Description  Creates a new recurring expense for the authenticated workspace
// @Tags         accounting
// @Accept       json
// @Produce      json
// @Param        request body CreateRecurringExpenseRequest true "Recurring expense data"
// @Success      201 {object} RecurringExpense
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/recurring-expenses [post]
// @Security     BearerAuth
func (h *HttpHandler) CreateRecurringExpense(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req CreateRecurringExpenseRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	recurringExpense, err := h.service.CreateRecurringExpense(c.Request.Context(), actor, biz, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessJSON(c, http.StatusCreated, recurringExpense)
}

// UpdateRecurringExpense updates an existing recurring expense
//
// @Summary      Update recurring expense
// @Description  Updates an existing recurring expense for the authenticated workspace
// @Tags         accounting
// @Accept       json
// @Produce      json
// @Param        recurringExpenseId path string true "Recurring Expense ID"
// @Param        request body UpdateRecurringExpenseRequest true "Recurring expense update data"
// @Success      200 {object} RecurringExpense
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/recurring-expenses/{recurringExpenseId} [patch]
// @Security     BearerAuth
func (h *HttpHandler) UpdateRecurringExpense(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	recurringExpenseID := c.Param("recurringExpenseId")
	if recurringExpenseID == "" {
		response.Error(c, problem.BadRequest("recurringExpenseId is required"))
		return
	}

	var req UpdateRecurringExpenseRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	recurringExpense, err := h.service.UpdateRecurringExpense(c.Request.Context(), actor, biz, recurringExpenseID, &req)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrRecurringExpenseNotFound(err))
			return
		}
		response.Error(c, err)
		return
	}

	response.SuccessJSON(c, http.StatusOK, recurringExpense)
}

// DeleteRecurringExpense deletes a recurring expense
//
// @Summary      Delete recurring expense
// @Description  Deletes a recurring expense for the authenticated workspace
// @Tags         accounting
// @Produce      json
// @Param        recurringExpenseId path string true "Recurring Expense ID"
// @Success      204
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/recurring-expenses/{recurringExpenseId} [delete]
// @Security     BearerAuth
func (h *HttpHandler) DeleteRecurringExpense(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	recurringExpenseID := c.Param("recurringExpenseId")
	if recurringExpenseID == "" {
		response.Error(c, problem.BadRequest("recurringExpenseId is required"))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	err = h.service.DeleteRecurringExpense(c.Request.Context(), actor, biz, recurringExpenseID)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrRecurringExpenseNotFound(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessEmpty(c, http.StatusNoContent)
}

type updateRecurringExpenseStatusRequest struct {
	Status RecurringExpenseStatus `json:"status" binding:"required,oneof=active paused ended canceled"`
}

// UpdateRecurringExpenseStatus updates the status of a recurring expense
//
// @Summary      Update recurring expense status
// @Description  Updates the status of a recurring expense (active, paused, ended, canceled)
// @Tags         accounting
// @Accept       json
// @Produce      json
// @Param        recurringExpenseId path string true "Recurring Expense ID"
// @Param        request body updateRecurringExpenseStatusRequest true "Status update data"
// @Success      200 {object} RecurringExpense
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      409 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/recurring-expenses/{recurringExpenseId}/status [patch]
// @Security     BearerAuth
func (h *HttpHandler) UpdateRecurringExpenseStatus(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	recurringExpenseID := c.Param("recurringExpenseId")
	if recurringExpenseID == "" {
		response.Error(c, problem.BadRequest("recurringExpenseId is required"))
		return
	}

	var req updateRecurringExpenseStatusRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	recurringExpense, err := h.service.UpdateRecurringExpenseStatus(c.Request.Context(), actor, biz, recurringExpenseID, req.Status)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrRecurringExpenseNotFound(err))
			return
		}
		response.Error(c, err)
		return
	}

	response.SuccessJSON(c, http.StatusOK, recurringExpense)
}

// GetRecurringExpenseOccurrences returns all expense occurrences for a recurring expense
//
// @Summary      Get recurring expense occurrences
// @Description  Returns all expense occurrences (instances) for a specific recurring expense
// @Tags         accounting
// @Produce      json
// @Param        recurringExpenseId path string true "Recurring Expense ID"
// @Success      200 {array} Expense
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/recurring-expenses/{recurringExpenseId}/occurrences [get]
// @Security     BearerAuth
func (h *HttpHandler) GetRecurringExpenseOccurrences(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	recurringExpenseID := c.Param("recurringExpenseId")
	if recurringExpenseID == "" {
		response.Error(c, problem.BadRequest("recurringExpenseId is required"))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	// Verify recurring expense exists and belongs to workspace
	_, err = h.service.GetRecurringExpenseByID(c.Request.Context(), actor, biz, recurringExpenseID)
	if err != nil {
		if database.IsRecordNotFound(err) {
			response.Error(c, ErrRecurringExpenseNotFound(err))
			return
		}
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	occurrences, err := h.service.GetRecurringExpenseOccurrences(c.Request.Context(), actor, biz, recurringExpenseID)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessJSON(c, http.StatusOK, occurrences)
}

// Summary endpoints

type summaryQuery struct {
	From string `form:"from" binding:"omitempty"`
	To   string `form:"to" binding:"omitempty"`
}

type summaryResponse struct {
	TotalAssetValue  string `json:"totalAssetValue"`
	TotalInvestments string `json:"totalInvestments"`
	TotalWithdrawals string `json:"totalWithdrawals"`
	TotalExpenses    string `json:"totalExpenses"`
	SafeToDrawAmount string `json:"safeToDrawAmount"`
	Currency         string `json:"currency"`
	From             string `json:"from,omitempty"`
	To               string `json:"to,omitempty"`
}

// GetAccountingSummary returns a summary of accounting metrics
//
// @Summary      Get accounting summary
// @Description  Returns a summary of key accounting metrics for the authenticated workspace
// @Tags         accounting
// @Produce      json
// @Param        from query string false "Start date (YYYY-MM-DD)"
// @Param        to query string false "End date (YYYY-MM-DD)"
// @Success      200 {object} summaryResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/accounting/summary [get]
// @Security     BearerAuth
func (h *HttpHandler) GetAccountingSummary(c *gin.Context) {
	actor, err := account.ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var query summaryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, problem.BadRequest("invalid query parameters").WithError(err))
		return
	}

	biz, err := h.getBusinessForWorkspace(c, actor)
	if err != nil {
		response.Error(c, err)
		return
	}

	// Parse dates
	var from, to time.Time
	if query.From != "" {
		from, err = time.Parse("2006-01-02", query.From)
		if err != nil {
			response.Error(c, problem.BadRequest("invalid from date format, use YYYY-MM-DD").WithError(err))
			return
		}
	}
	if query.To != "" {
		to, err = time.Parse("2006-01-02", query.To)
		if err != nil {
			response.Error(c, problem.BadRequest("invalid to date format, use YYYY-MM-DD").WithError(err))
			return
		}
	}

	// Get metrics
	totalAssetValue, err := h.service.SumAssetsValue(c.Request.Context(), actor, biz, from, to)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	totalInvestments, err := h.service.SumInvestmentsAmount(c.Request.Context(), actor, biz, from, to)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	totalWithdrawals, err := h.service.SumWithdrawalsAmount(c.Request.Context(), actor, biz, from, to)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	totalExpenses, err := h.service.SumExpensesAmount(c.Request.Context(), actor, biz, from, to)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	// For safe to draw, we need total income and COGS.
	// For now, we treat investments as income and assume COGS = 0 (until orders/COGS are wired in).
	safeToDrawAmount, err := h.service.ComputeSafeToDrawAmount(c.Request.Context(), actor, biz, totalInvestments, decimal.Zero, from, to)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	summary := summaryResponse{
		TotalAssetValue:  totalAssetValue.String(),
		TotalInvestments: totalInvestments.String(),
		TotalWithdrawals: totalWithdrawals.String(),
		TotalExpenses:    totalExpenses.String(),
		SafeToDrawAmount: safeToDrawAmount.String(),
		Currency:         biz.Currency,
		From:             query.From,
		To:               query.To,
	}

	response.SuccessJSON(c, http.StatusOK, summary)
}
