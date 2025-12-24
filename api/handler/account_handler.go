package handler

import (
	"net/http"
	"strconv"

	"simple_bank/internal/service"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountService  service.AccountService
	transferService service.TransferService
}

func NewAccountHandler(router *gin.RouterGroup,
	accountService service.AccountService,
	transferService service.TransferService) {

	handler := &AccountHandler{
		accountService:  accountService,
		transferService: transferService,
	}

	accounts := router.Group("/accounts")
	{
		// Static routes first
		accounts.POST("", handler.CreateAccount)
		accounts.GET("", handler.ListAccounts)
		accounts.GET("/owner/:owner", handler.GetAccountsByOwner)

		// Dynamic routes with specific names
		accounts.GET("/:account_id", handler.GetAccount)       // Changed from :id to :account_id
		accounts.PUT("/:account_id", handler.UpdateAccount)    // Changed from :id to :account_id
		accounts.DELETE("/:account_id", handler.DeleteAccount) // Changed from :id to :account_id

		// Transfer routes (use different param name)
		accounts.POST("/:from_account_id/transfer", handler.CreateTransfer)
		accounts.GET("/:account_id/transfers", handler.ListTransfers)
	}

	transfers := router.Group("/transfers")
	{
		transfers.GET("/:transfer_id", handler.GetTransfer) // Changed from :id to :transfer_id
	}
}

// Request/Response structures
type CreateAccountRequest struct {
	Owner          string `json:"owner" binding:"required"`
	Currency       string `json:"currency" binding:"required"`
	InitialBalance int64  `json:"initial_balance" binding:"min=0"`
}

type CreateAccountResponse struct {
	ID        int64  `json:"id"`
	Owner     string `json:"owner"`
	Balance   int64  `json:"balance"`
	Currency  string `json:"currency"`
	CreatedAt string `json:"created_at"`
}

type UpdateAccountRequest struct {
	Balance int64 `json:"balance" binding:"min=0"`
}

type CreateTransferRequest struct {
	ToAccountID int64 `json:"to_account_id" binding:"required,gt=0"`
	Amount      int64 `json:"amount" binding:"required,gt=0"`
}

// Handler methods
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.accountService.CreateAccount(c.Request.Context(),
		req.Owner, req.Currency, req.InitialBalance)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateAccountResponse{
		ID:        account.ID,
		Owner:     account.Owner,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt.Format("2006-01-02 15:04:05"),
	})
}

func (h *AccountHandler) GetAccount(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	account, err := h.accountService.GetAccount(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (h *AccountHandler) ListAccounts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	accounts, err := h.accountService.ListAccounts(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts":  accounts,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *AccountHandler) GetAccountsByOwner(c *gin.Context) {
	owner := c.Param("owner")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	accounts, err := h.accountService.GetAccountsByOwner(c.Request.Context(), owner, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts":  accounts,
		"owner":     owner,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	var req UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.accountService.UpdateAccount(c.Request.Context(), id, req.Balance)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	err = h.accountService.DeleteAccount(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

func (h *AccountHandler) CreateTransfer(c *gin.Context) {
	var req CreateTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fromAccountID, err := strconv.ParseInt(c.Param("from_account_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	transfer, err := h.transferService.CreateTransfer(c.Request.Context(), fromAccountID, req.ToAccountID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, transfer)
}

func (h *AccountHandler) GetTransfer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transfer ID"})
		return
	}

	transfer, err := h.transferService.GetTransfer(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transfer not found"})
		return
	}

	c.JSON(http.StatusOK, transfer)
}

func (h *AccountHandler) ListTransfers(c *gin.Context) {
	accountID, err := strconv.ParseInt(c.Param("account_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	transfers, err := h.transferService.ListTransfers(c.Request.Context(), accountID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transfers": transfers,
		"page":      page,
		"page_size": pageSize,
	})
}
