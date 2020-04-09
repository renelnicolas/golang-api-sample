package models

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"ohmytech.io/platform/helpers"
)

// ContextKey :
type ContextKey string

// ContextUserKey :
const ContextUserKey ContextKey = "user"

// UserMain :v
type UserMain struct {
	Firstname string   `json:"firstName"`
	Lastname  string   `json:"lastName"`
	Email     string   `json:"email"`
	Country   *Country `json:"country"`
}

// User :
type User struct {
	UserMain
	ID                int64                   `json:"id"`
	Password          NullToEmptyString       `json:"-"`
	Phone             NullToEmptyString       `json:"phone"`
	ExternalID        NullToEmptyString       `json:"externalId"`
	Company           *Company                `json:"company"`
	Roles             Uint8ArrayToArrayString `json:"roles"`
	Enabled           bool                    `json:"enabled"`
	CompanyRegistered NullToEmptyString       `json:"company_registered"`
	RegistredIP       string                  `json:"registred_ip"`
	RegistredAt       string                  `json:"registred_at"`
}

// UserSignIn :
type UserSignIn struct {
	UserMain
	ExternalID NullToEmptyString       `json:"id"`
	Token      string                  `json:"token"`
	Roles      Uint8ArrayToArrayString `json:"roles"`
	Company    *Company                `json:"company"`
}

// UserSignUp :
type UserSignUp struct {
	UserMain
	Phone             string `json:"phone"`
	Password          string `json:"password"`
	CompanyRegistered string `json:"company_registered"`
	RegistredIP       string `json:"registred_ip"`
	RegistredAt       string `json:"registred_at"`
}

// UserConnection :
type UserConnection struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	RequestedAt string `json:"requested_at"`
	RequestedIP string `json:"requested_ip"`
}

// UserToken :
type UserToken struct {
	User
	ExpirationTime time.Time `json:"exp"`
}

// UserClaims : JWT claims struct
type UserClaims struct {
	UserToken UserToken
	jwt.StandardClaims
}

// Users :
type Users []User

// UserAlias :
type UserAlias User

// NewUserConnection : Returns a new UserConnection struct with default infos
func NewUserConnection(remoteAddr string, headers http.Header) UserConnection {
	uc := UserConnection{
		RequestedAt: time.Now().Format("2006-01-02 15:04:05"),
		RequestedIP: helpers.GetIP(remoteAddr, headers),
	}

	return uc
}

// ToClaims :
func (user User) ToClaims() UserClaims {
	expT := time.Now().Add(24 * time.Hour) // 24h

	usrMain := UserMain{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Country:   user.Country,
	}

	usr := User{
		UserMain:   usrMain,
		ID:         user.ID,
		ExternalID: user.ExternalID,
		Company:    user.Company,
		Roles:      user.Roles,
	}

	return UserClaims{
		UserToken: UserToken{
			User:           usr,
			ExpirationTime: expT,
		},
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expT.Unix(),
		},
	}
}

// ToUser :
func (userc UserConnection) ToUser() User {
	usrMain := UserMain{
		Email: userc.Email,
	}

	usr := User{
		UserMain: usrMain,
		Password: NullToEmptyString(userc.Password),
	}

	return usr
}

// ToUserSignIn :
func (usert UserToken) ToUserSignIn(token string) UserSignIn {
	usrsi := UserSignIn{
		UserMain:   usert.User.UserMain,
		ExternalID: usert.ExternalID,
		Company:    usert.Company,
		Roles:      usert.Roles,
		Token:      token,
	}

	return usrsi
}

// ToUser :
func (usert UserToken) ToUser() User {
	usr := User{
		UserMain:   usert.User.UserMain,
		ExternalID: usert.ExternalID,
		Company:    usert.Company,
		Roles:      usert.Roles,
	}

	return usr
}
