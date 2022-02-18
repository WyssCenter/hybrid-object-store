package userinfo

import (
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const ValidCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_."

const (
	ROLE_ADMIN      = "admin"
	ROLE_PRIVILEGED = "privileged"
	ROLE_USER       = "user"
)

// const (
// 	GROUP_ADMIN      = "admins"
// 	GROUP_PRIVILEGED = "developers"
// )

type UserInfo struct {
	Subject       string
	FullName      string
	GivenName     string
	FamilyName    string
	Username      string
	Email         string
	EmailVerified *bool
	Role          string
	Groups        []string
}

func ParseUserClaims(userClaims jwt.MapClaims, adminGroup string, privilegedGroup string, nicknameClaim string) (UserInfo, error) {
	userInfo := UserInfo{}

	isService, ok := userClaims["service"]
	if ok && isService != nil && isService.(bool) {
		role, ok := userClaims["role"]
		if ok && role != nil {
			userInfo.Role = role.(string)
		}

		nick, ok := userClaims["nickname"]
		if ok && nick != nil {
			userInfo.Username = nick.(string)
		}

		sub, ok := userClaims["sub"]
		if ok && sub != nil {
			userInfo.Subject = sub.(string)
		}

		return userInfo, nil
	}

	email, ok := userClaims["email"]
	if ok && email != nil {
		userInfo.Email = email.(string)
	} else {
		return userInfo, errors.New("No email claim")
	}

	emailVerified, ok := userClaims["email_verified"]
	if ok && emailVerified != nil {
		emailVerifiedCast := emailVerified.(bool)
		userInfo.EmailVerified = &emailVerifiedCast
	} else {
		userInfo.EmailVerified = nil
	}

	// if claims are from gigantum JWT the role claim should already be set
	role, ok := userClaims["role"]
	if ok && role != nil {
		userInfo.Role = role.(string)
	} else {
		// if claims come from a dex JWT, the role needs to be converted from the ldap groups
		groups, ok := userClaims["groups"]
		if ok && groups != nil {
			for _, group := range groups.([]interface{}) {
				if strings.Contains(group.(string), adminGroup) {
					userInfo.Role = ROLE_ADMIN
					break
				} else if strings.Contains(group.(string), privilegedGroup) {
					userInfo.Role = ROLE_PRIVILEGED
				}
			}
		}
	}
	// if neither role or groups claims were set, user gets lowest privilege role
	if userInfo.Role == "" {
		userInfo.Role = ROLE_USER
	}

	name, ok := userClaims["name"]
	if ok && name != nil {
		userInfo.FullName = name.(string)
	} else {
		given, ok1 := userClaims["given_name"]
		family, ok2 := userClaims["family_name"]
		if ok1 && ok2 && given != nil && family != nil {
			userInfo.FullName = given.(string) + " " + family.(string)
		} else {
			userInfo.FullName = "* *" // null value, as understood by Gigantum
		}
	}

	given, ok := userClaims["given_name"]
	if ok && given != nil {
		userInfo.GivenName = given.(string)
	} else {
		userInfo.GivenName = "*" // null value, as understood by Gigantum
	}

	family, ok := userClaims["family_name"]
	if ok && family != nil {
		userInfo.FamilyName = family.(string)
	} else {
		userInfo.FamilyName = "*" // null value, as understood by Gigantum
	}

	// Load the "nickname", which is used as the username in the system
	// If nicknameClaim is set use the indicated claim. If not, try to find a
	// reasonable username
	if nicknameClaim == "nickname" {
		logrus.Info("NICNAME USERNAME")
		nick, ok := userClaims["nickname"]
		if ok && nick != nil {
			userInfo.Username = nick.(string)
		} else {
			return userInfo, errors.New("Failed to load nickname from nickname claim")
		}
	} else if nicknameClaim == "name" {
		logrus.Info("NAME USERNAME")
		name, ok := userClaims["name"]
		if ok && name != nil {
			userInfo.Username = name.(string)
		} else {
			return userInfo, errors.New("Failed to load nickname from name claim")
		}
	} else if nicknameClaim == "email" {
		logrus.Info("EMAIL USERNAME")
		email, err := emailToNickname(userClaims)
		if err != nil {
			return userInfo, errors.New("Failed to load nickname from email claim")
		}
		userInfo.Username = email
	} else {
		nickname, err := findNickname(userClaims)
		if err != nil {
			return userInfo, errors.New("Failed to find nickname from claims")
		}
		userInfo.Username = nickname
	}

	// Verify that the username only contains valid characters
	mut := []rune(userInfo.Username)
	for i, c := range userInfo.Username {
		if !strings.ContainsAny(ValidCharacters, string(c)) {
			mut[i] = '-'
		}
	}
	userInfo.Username = string(mut)

	sub, ok := userClaims["sub"]
	if ok && sub != nil {
		userInfo.Subject = sub.(string)
	} else {
		userInfo.Subject = userInfo.Username
	}

	return userInfo, nil
}

// emailToNickname extracts a username from an email address by stripping off the domain
func emailToNickname(userClaims jwt.MapClaims) (string, error) {
	email, ok := userClaims["email"]
	if ok && email != nil {
		nickname := strings.Split(email.(string), "@")[0] // strip the domain
		nickname = strings.Split(nickname, "+")[0]        // strip of any + suffix if present
		return nickname, nil
	}

	return "", errors.New("No email claim to parse")
}

// findNickname attempts to discover a username from the nickname, name, and email claims
func findNickname(userClaims jwt.MapClaims) (string, error) {
	nick, ok := userClaims["nickname"]
	if ok && nick != nil {
		return nick.(string), nil
	} else {
		name, ok := userClaims["name"]
		if ok && name != nil {
			return name.(string), nil
		} else {
			email, err := emailToNickname(userClaims)
			if err == nil {
				return email, nil
			} else {
				return "", errors.New("Failed to find nickname")
			}
		}
	}
}
