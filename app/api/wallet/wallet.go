package wallet

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/n3k0fi5t/wallet/app/middleware"
	"github.com/n3k0fi5t/wallet/app/service/wallet"
)

// NewHandler ...
func NewHandler(w wallet.Service) *Handler {
	return &Handler{
		walletSrv: w,
	}
}

type Handler struct {
	walletSrv wallet.Service
}

func (h *Handler) Handle(routerGroup *gin.RouterGroup) {
	rg := routerGroup.Group("/wallet")

	// APIs are only for authed user
	rg.Use(middleware.GetUserAccount())

	// trade relative
	rg.Handle("POST", "/deposit", h.deposit)
	rg.Handle("POST", "/withdraw", h.withdraw)
	rg.Handle("POST", "/transfer", h.transfer)

	// account relative
	arg := rg.Group("/account")
	arg.Handle("GET", "", h.getAccountInfo)
}

type depositParam struct {
	Amount int64 `json:"amount"`
}

type depositResp struct {
	TradeID string `json:"tradeID"`
}

func (h *Handler) deposit(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	accountID := c.MustGet("accountID").(string)

	param := withdrawParam{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"errMessage": err.Error(),
		})
		return
	}

	tradeID, err := h.walletSrv.Deposit(ctx, accountID, param.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"errMessage": err.Error(),
		})
		return
	}

	resp := withdrawResp{
		TradeID: tradeID,
	}
	c.JSON(http.StatusOK, resp)
}

type withdrawParam struct {
	Amount int64 `json:"amount"`
}

type withdrawResp struct {
	TradeID string `json:"tradeID"`
}

func (h *Handler) withdraw(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	accountID := c.MustGet("accountID").(string)

	param := withdrawParam{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"errMessage": err.Error(),
		})
		return
	}

	tradeID, err := h.walletSrv.Withdraw(ctx, accountID, param.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"errMessage": err.Error(),
		})
		return
	}

	resp := withdrawResp{
		TradeID: tradeID,
	}
	c.JSON(http.StatusOK, resp)
}

type transferParam struct {
	Amount    int64  `json:"amount"`
	ToAccount string `json:"toAccount"`
}

type transferResp struct {
	TradeID string `json:"tradeID"`
}

func (h *Handler) transfer(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	accountID := c.MustGet("accountID").(string)

	param := transferParam{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"errMessage": err.Error(),
		})
		return
	}

	tradeID, err := h.walletSrv.Transfer(ctx, accountID, param.ToAccount, param.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"errMessage": err.Error(),
		})
		return
	}

	resp := transferResp{
		TradeID: tradeID,
	}
	c.JSON(http.StatusOK, resp)
}

type accountInfoResp struct {
	AccountID string `json:"accountID"`
	Balance   int64  `json:"balance"`
}

func (h *Handler) getAccountInfo(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	accountID := c.MustGet("accountID").(string)

	account, err := h.walletSrv.GetAccount(ctx, accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"errMessage": err.Error(),
		})
		return
	}

	resp := accountInfoResp{
		AccountID: account.AccountID,
		Balance:   account.Balance,
	}
	c.JSON(http.StatusOK, resp)
}
