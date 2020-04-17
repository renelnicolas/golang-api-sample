package validators

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"ohmytech.io/platform/helpers"
	"ohmytech.io/platform/models"
	"ohmytech.io/platform/repositories"
)

var (
	signingKey = []byte(os.Getenv("SECRET_KEY_JWT"))
)

// JwtAuthentication :
func JwtAuthentication(authenticationHeader string, userClaims *models.UserClaims) error {
	jwttoken, err := helpers.AuthenticationHeader(authenticationHeader)
	if nil != err {
		return err
	}

	claims, err := helpers.ReadToken(jwttoken, signingKey)
	if nil != err {
		return fmt.Errorf("ReadToken error: %s", err.Error())
	}

	muc, err := json.Marshal(claims)
	if nil != err {
		return fmt.Errorf("Conversion error: %s", err.Error())
	}

	err = json.Unmarshal(muc, &userClaims)
	if nil != err {
		return fmt.Errorf("Conversion error: %s", err.Error())
	}

	return nil
}

// Authentication :
func Authentication(uc models.UserConnection) (*models.UserSignIn, error) {
	err := validateValue(uc)
	if nil != err {
		return nil, err
	}

	user := uc.ToUser()

	err = validateUserExist(&user)
	if nil != err {
		return nil, err
	}

	err = validatePassword(user, uc)
	if nil != err {
		return nil, err
	}

	jwtClaims := user.ToClaims()

	token, err := helpers.GenerateToken(jwtClaims, signingKey)
	if nil != err {
		return nil, err
	}

	userSignIn := jwtClaims.UserToken.ToUserSignIn(token)

	return &userSignIn, nil
}

// validateValue :
func validateValue(userConection models.UserConnection) error {
	if "" == strings.TrimFunc(userConection.Email, helpers.TrimWhitespaceFn) {
		return errors.New("Email cannot be empty")
	}

	if "" == strings.TrimFunc(userConection.Password, helpers.TrimWhitespaceFn) {
		return errors.New("Password cannot be empty")
	}

	return nil
}

// validateUserExist :
func validateUserExist(user *models.User) error {
	repo := repositories.UserRepository{}

	err := repo.FindOneByEmail(user)
	if nil != err {
		return err
	}

	return nil
}

// validatePassword :
func validatePassword(user models.User, userConection models.UserConnection) error {
	return helpers.ComparePasswords(string(user.Password), userConection.Password)
}
