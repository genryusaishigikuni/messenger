package jwt

import (
	"strconv"
	"time"

	"github.com/genryusaishigikuni/messenger/auth-service/pkg/models"
	"github.com/genryusaishigikuni/messenger/auth-service/pkg/utils"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(secret string, userID int, username string) (string, error) {
	utils.Info("Generating token...")

	claims := &models.TokenClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		utils.Error("Failed to sign token: " + err.Error())
		return "", err
	}

	utils.Info("Token generated successfully for user ID: " + strconv.Itoa(userID))
	return signedToken, nil
}

func ValidateToken(secret, tokenString string) (*models.TokenClaims, error) {
	utils.Info("Validating token...")

	token, err := jwt.ParseWithClaims(tokenString, &models.TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		utils.Error("Failed to parse token: " + err.Error())
		return nil, err
	}

	if claims, ok := token.Claims.(*models.TokenClaims); ok && token.Valid {
		utils.Info("Token validated successfully for user ID: " + strconv.Itoa(claims.UserID))
		return claims, nil
	}

	utils.Error("Invalid token signature")
	return nil, jwt.ErrSignatureInvalid
}
