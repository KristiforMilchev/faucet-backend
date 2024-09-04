package implementations

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"golang.org/x/sync/syncmap"

	"faucet/api/dtos"
)

type AuthenticationService struct {
	AuthRequests syncmap.Map
}

func (authService *AuthenticationService) GetMessage(email *string, id *uuid.UUID) dtos.InitAuthReponse {

	uuid := uuid.New()
	code := generateRandomString(32)

	authResponse := dtos.StoredAuthRequest{
		Id:   *id,
		Name: *email,
		Code: code,
		Uuid: uuid.String(),
		Time: time.Time{}.Add(time.Duration(time.Minute * 5)),
	}

	authService.AuthRequests.LoadOrStore(uuid, authResponse)

	fmt.Println(&authService.AuthRequests)

	response := dtos.InitAuthReponse{
		Code: authResponse.Code,
		Uuid: authResponse.Uuid,
	}

	return response
}

func (authService *AuthenticationService) GetRequestById(id uuid.UUID) uuid.UUID {
	request, _ := authService.AuthRequests.Load(id)

	return request.(dtos.StoredAuthRequest).Id
}

func (authService *AuthenticationService) VerifyEVMSignature(current uuid.UUID, signature string) (string, error) {
	request, _ := authService.AuthRequests.Load(current)
	authRequest := request.(dtos.StoredAuthRequest)

	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(authRequest.Code), authRequest.Code)
	msgHash := crypto.Keccak256Hash([]byte(msg))
	sig, err := hexutil.Decode(signature)
	if err != nil {
		return "", fmt.Errorf("invalid signature: %v", err)
	}

	if len(sig) != 65 {
		return "", fmt.Errorf("invalid signature length")
	}

	r := sig[:32]
	s := sig[32:64]
	v := sig[64]

	if v < 27 {
		v += 27
	}

	pubKey, err := crypto.Ecrecover(msgHash.Bytes(), append(r, append(s, v-27)...))
	if err != nil {
		return "", fmt.Errorf("signature verification failed: %v", err)
	}

	recoveredPubKey, err := crypto.UnmarshalPubkey(pubKey)
	if err != nil {
		return "", fmt.Errorf("unable to unmarshal public key: %v", err)
	}

	recoveredAddr := crypto.PubkeyToAddress(*recoveredPubKey).Hex()
	if strings.ToLower(recoveredAddr) != strings.ToLower(authRequest.Name) {
		log.Printf("Recovered address (%s) does not match provided address (%s)", recoveredAddr, authRequest.Name)
		return "", nil
	}

	return recoveredAddr, nil
}

func (authService *AuthenticationService) Start() {
	// Create a channel to receive signals
	signal := make(chan struct{})

	// Start a goroutine to send signals at regular intervals
	go func() {
		for {
			time.Sleep(30 * time.Second) // Wait for 30 seconds
			signal <- struct{}{}         // Send a signal to the channel
		}
	}()

	for range signal {
		fmt.Println("Something has happened at", time.Now())
		authService.AuthRequests.Range(func(key, value any) bool {
			storedRequest := value.(dtos.StoredAuthRequest)
			expired := storedRequest.Time.UTC().After(time.Now().UTC())
			fmt.Println(storedRequest.Name)
			fmt.Println(storedRequest.Code)
			if expired {
				fmt.Println("Deleting")

				authService.AuthRequests.Delete(key)
			}

			return true
		})
	}

}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result string
	for i := 0; i < length; i++ {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result += string(charset[randomIndex.Int64()])
	}
	return result
}
