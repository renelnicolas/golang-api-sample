package helpers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// AuthenticationHeader :
func AuthenticationHeader(authorization string) (string, error) {
	if "" == authorization { //Token is missing, returns with error code 403 Unauthorized
		return "", errors.New("Missing header authentication token")
	}

	splitted := strings.Split(authorization, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
	if 2 != len(splitted) {
		return "", errors.New("Invalid/Malformed auth token - Authorization")
	}

	authType := splitted[0]
	jwttoken := splitted[1]

	if "Bearer" != authType {
		return "", errors.New("Invalid/Malformed auth token - Bearer")
	}

	return jwttoken, nil
}

// GenerateToken :
func GenerateToken(claims jwt.Claims, secret []byte) (string, error) {
	/* Create a map to store our claims */
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	/* Sign the token with our secret */
	signedToken, err := token.SignedString(secret)
	if nil != err {
		return "", err
	}

	return signedToken, nil
}

// ValidateToken :
func ValidateToken(sToken string, secret []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(sToken, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if nil != token && token.Valid {
		return token, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, errors.New("jwtoken : ValidateToken - That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return nil, errors.New("jwtoken : ValidateToken - Timing is everything")
		} else {
			return nil, fmt.Errorf("jwtoken : ValidateToken - Couldn't handle this token: %s", err.Error())
		}
	}

	return nil, fmt.Errorf("jwtoken : ValidateToken - Couldn't handle this token : %s", err.Error())
}

// ReadToken :
func ReadToken(sToken string, secret []byte) (jwt.MapClaims, error) {
	vt, err := ValidateToken(sToken, secret)

	if nil != err {
		return jwt.MapClaims{}, err
	}

	if claims, ok := vt.Claims.(jwt.MapClaims); ok && vt.Valid {
		return claims, nil
	}

	return jwt.MapClaims{}, errors.New("ReadToken : Cannot validate token")
}
