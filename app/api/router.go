package api

import (
	"github.com/gin-gonic/gin"
	"github.com/n3k0fi5t/wallet/app/api/wallet"
	"github.com/n3k0fi5t/wallet/app/middleware"
	"github.com/n3k0fi5t/wallet/app/repository/bank"
	wSrv "github.com/n3k0fi5t/wallet/app/service/wallet"
	"github.com/n3k0fi5t/wallet/app/setup/mysql"
)

func BuildWalletHandler() *wallet.Handler {
	// It's could be better if we use DI container
	db := mysql.GetMySQL()
	b := bank.NewBank(db)
	walletSrv := wSrv.NewWallet(b)
	return wallet.NewHandler(walletSrv)
}

func BuildRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1")

	// set context for following process
	api.Use(middleware.SetHandleContext())

	walletHandler := BuildWalletHandler()
	walletHandler.Handle(api)

	return router
}
