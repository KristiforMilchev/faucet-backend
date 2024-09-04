package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	implementations "faucet/api/Implementations"
	"faucet/api/interfaces"
	"faucet/api/middlewhere"
	"faucet/api/repositories"
)

type AccountsController struct {
	AccountRepository repositories.Accounts
	Storage           interfaces.Storage
	PaymentProcessor  *implementations.PaymentProcessor
}

func (a *AccountsController) getAccountNextDrip(ctx *gin.Context) {
	address := ctx.MustGet("ID")
	a.AccountRepository.Storage.Open(a.AccountRepository.ConnectionString)
	defer a.AccountRepository.Storage.Close()

	user, exists := a.AccountRepository.UserExists(address.(string))
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	nextLease := user.LastLease.Add(time.Hour * 24)
	ctx.JSON(http.StatusOK, nextLease)
}

func (a *AccountsController) drip(ctx *gin.Context) {
	address := ctx.MustGet("ID")
	a.AccountRepository.Storage.Open(a.AccountRepository.ConnectionString)
	defer a.AccountRepository.Storage.Close()

	user, exists := a.AccountRepository.UserExists(address.(string))
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	nextLease := user.LastLease.Add(time.Hour * 24)
	if time.Now().Before(nextLease) {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Allownace exceeded, please try later!"})
		return
	}

	payment, err := a.PaymentProcessor.ProcessNative(user.Address)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Failed to process payment, please try later!"})
		return
	}

	a.AccountRepository.UpdateDrip(address.(string))
	ctx.JSON(http.StatusOK, payment)

}

func (ac *AccountsController) Init(r *gin.RouterGroup, authMiddlewhere *middlewhere.AuthenticationMiddlewhere) {
	accounts := r.Group("accounts")

	accounts.Use(authMiddlewhere.Authorize())
	accounts.GET("get-drip", ac.getAccountNextDrip)
	accounts.GET("drip", ac.drip)
}
