package middlewhere

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"faucet/api/interfaces"
)

type UnsignedResponse struct {
	Message interface{} `json:"message"`
}

type AuthenticationMiddlewhere struct {
	JwtService interfaces.JwtService
}

func (aw *AuthenticationMiddlewhere) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken, err := extractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
				Message: err.Error(),
			})
			return
		}

		result := aw.JwtService.ValidateToken(jwtToken)
		if !result {
			c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
				Message: "Invalid access token",
			})
			return
		}

		userId := aw.JwtService.ExtractValue(jwtToken, "id")
		userID, ok := userId.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
				Message: "Invalid token claims",
			})
			return
		}
		c.Set("ID", userID)
		c.Copy().Next()
	}
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}
