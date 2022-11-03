package route

import (
	"github.com/gin-gonic/gin"
	"miniWalletExercise/controller"
	"miniWalletExercise/middleware"
	"miniWalletExercise/repository"
)

func NewHttpHandlerFromMethod(walletInterface repository.WalletInterface, transactionInterface repository.TransactionInterface) controller.RepositoryHandler {
	return controller.RepositoryHandler{
		Wallet:      walletInterface,
		Transaction: transactionInterface,
	}
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	walletHandler := repository.NewWalletRepository()
	transactionHandler := repository.NewTransactionRepository()

	r.POST("/api/v1/init", NewHttpHandlerFromMethod(walletHandler, transactionHandler).InitializeAccount)

	miniWalletRoute := r.Group("/api/v1/wallet")
	miniWalletRoute.Use(middleware.JwtAuthMiddleware())
	miniWalletRoute.POST("/", NewHttpHandlerFromMethod(walletHandler, transactionHandler).EnableWallet)
	miniWalletRoute.GET("/", NewHttpHandlerFromMethod(walletHandler, transactionHandler).ViewWallet)
	miniWalletRoute.POST("/deposits", NewHttpHandlerFromMethod(walletHandler, transactionHandler).DepositWallet)
	miniWalletRoute.POST("/withdrawals", NewHttpHandlerFromMethod(walletHandler, transactionHandler).WithdrawalWallet)
	miniWalletRoute.PATCH("/", NewHttpHandlerFromMethod(walletHandler, transactionHandler).DisableWallet)

	return r
}
