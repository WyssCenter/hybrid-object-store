package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"

	userinfo "github.com/gigantum/hoss-auth/pkg/userinfo"
)

var (
	ErrUnauthorized = errors.New("user is not authorized")
)

func getUser(c *gin.Context) (userinfo.UserInfo, error) {
	user, ok := c.Get("user")
	if ok {
		return user.(userinfo.UserInfo), nil
	} else {
		return userinfo.UserInfo{}, errors.New("User info not available")
	}
}

const (
	ROLE_ADMIN      = "admin"
	ROLE_PRIVILEGED = "privileged"
	ROLE_USER       = "user"
)

// ValidatePrivileged checks if the user has a privileged or admin role
func validatePrivileged(role string) bool {
	return role == ROLE_PRIVILEGED || role == ROLE_ADMIN
}

// ValidateAdmin checks if the user has an admin role
func validateAdmin(role string) bool {
	return role == ROLE_ADMIN
}

func (a *Auth) validateAdminOrPrivilegedMember(role, username, groupName string) bool {
	if admin := validateAdmin(role); !admin {
		// if user is not an admin, they must be both privileged and a member
		privileged := validatePrivileged(role)
		member := a.Database.ValidateMembership(username, groupName)

		if !(privileged && member) {
			return false
		}
	}
	return true
}

func (a *Auth) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		kid := token.Header["kid"].(string)
		if a.KeyID != kid {
			return nil, errors.New("Token's key id does not match signing key")
		}

		return a.JwtKey.Public(), nil
	}

	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return jwt.MapClaims{}, err
	}

	if !token.Valid {
		return jwt.MapClaims{}, errors.New(fmt.Sprintf("Invalid JWT: %v", token.Claims))
	}

	userClaims := token.Claims.(jwt.MapClaims)

	return userClaims, nil
}
