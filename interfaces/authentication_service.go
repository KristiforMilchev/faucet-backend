package interfaces

import (
	"github.com/google/uuid"

	"faucet/api/dtos"
)

type AuthenticationService interface {
	Start()
	GetMessage(email *string, id *uuid.UUID) dtos.InitAuthReponse
	VerifyEVMSignature(current uuid.UUID, signature string) (string, error)
}
