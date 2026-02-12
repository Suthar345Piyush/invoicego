// jwt token access , refresh and validating code here
// TOKEN CREATION , TOKEN REFRESH AND TOKEN VALIDATION

package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string
	jwt.RegisteredClaims
}

// function to generate the access token

func GenerateAccessToken(userID uuid.UUID, email, secret string, expiry time.Duration) (string, error) {

	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// creating token

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))

}

// function for generating refresh rokens

func GenerateRefreshToken(userID uuid.UUID, secret string, expiry time.Duration) (string, error) {

	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))

}

// TOKEN VALIDATION FUNCTION

// here we are pointer to allow modification of struct

//INTERNALS of this function -

/*

  - we passed &JWTClaims{}
	- library fills it
	- library stored it inside token.Claims  as interface
	- then we are extracting it back and using it like this - .(*JWTClaims)


	Some efficient use case of pointer are like -

	 -> if we need to modify the struct , then use pointer
	 -> if their is an large struct , then use pointer
	 -> if we are just reading small amount of data , then don't use pointers
	 -> to return data efficiently , use pointers

*/

func ValidateToken(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(secret), nil

	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("Invalid token")

}
