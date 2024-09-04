package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"

	implementations "faucet/api/Implementations"
	"faucet/api/dtos"
	"faucet/api/interfaces"
	"faucet/api/repositories"
)

type AuthenticationController struct {
	AuthenticationService interfaces.AuthenticationService
	AccountRepository     repositories.Accounts
	Storage               interfaces.Storage
	JwtService            interfaces.JwtService
}

func (ac *AuthenticationController) begin(ctx *gin.Context) {
	var address dtos.BeginRequest
	if err := ctx.BindJSON(&address); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request!"})
		return
	}

	id, err := uuid.NewUUID()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Message": "Failed to get an identifier, reverting!"})
		return
	}

	response := ac.AuthenticationService.GetMessage(&address.Key, &id)
	ctx.JSON(http.StatusOK, response)
}

func (a *AuthenticationController) verifySignature(ctx *gin.Context) {
	var signatureData dtos.EVMAuthRequest
	if err := ctx.BindJSON(&signatureData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad input data!"})
		return
	}

	address, err := a.AuthenticationService.VerifyEVMSignature(signatureData.Id, signatureData.Signature)
	if err != nil || address == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Message": "Internal Error, login rejected. Bad Signature!"})
		return
	}

	a.checkNewUser(address)
	token := a.JwtService.IssueToken("user", address)
	ctx.JSON(http.StatusOK, token)
}

func (a *AuthenticationController) checkNewUser(address string) bool {
	a.AccountRepository.Storage.Open(a.AccountRepository.ConnectionString)
	defer a.AccountRepository.Storage.Close()

	_, exists := a.AccountRepository.UserExists(address)
	if exists {
		return true
	}

	_, err := a.AccountRepository.CreateUser(address)
	if err != nil {
		return false
	}

	return true
}

func (ac *AuthenticationController) Init(r *gin.RouterGroup) {
	ac.AuthenticationService = &implementations.AuthenticationService{}
	println("initializing Authentication Controller")
	go ac.AuthenticationService.Start()
	// Router group
	v1 := r.Group("/authentication")
	{
		v1.POST("/begin", ac.begin)
		v1.POST("/finish", ac.verifySignature)
	}
}
