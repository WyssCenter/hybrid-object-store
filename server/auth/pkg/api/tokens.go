package api

import (
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"

	"github.com/gigantum/hoss-auth/pkg/state"
	"github.com/gigantum/hoss-auth/pkg/userinfo"
)

// GetProviderTokens exchanges a code for JWT tokens and verifies them
func (a *Auth) GetProviderTokens(c *gin.Context, code string) (state.Tokens, jwt.MapClaims, error) {
	tokens := state.Tokens{}
	userClaims := jwt.MapClaims{}

	token, err := a.GetAuthConfig(c).Exchange(a.Ctx, code)
	if err != nil {
		return tokens, userClaims, err
	}

	tokens.AccessToken = token.AccessToken

	// Get id_token from the response
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return tokens, userClaims, errors.New("Could not get id_token")
	}

	// Verify using the OIDC Library, as it handles the JWK information
	_, err = a.Verifier.Verify(a.Ctx, idToken)
	if err != nil {
		log.Println("Error verifying ID Token: " + err.Error())
		return tokens, userClaims, err
	}

	tokens.IDToken = idToken

	// Extract the data from the ID token
	userClaims, err = a.GetTokenClaims(idToken)
	if err != nil {
		return tokens, userClaims, err
	}

	// Query userinfo service
	userInfo, err := a.Provider.UserInfo(a.Ctx, a.Config.TokenSource(a.Ctx, token))
	if err != nil {
		return tokens, userClaims, err
	}

	if err = userInfo.Claims(&userClaims); err != nil {
		return tokens, userClaims, err
	}

	return tokens, userClaims, nil
}

// GetTokenClaims extracts the user claims from an id token
func (a *Auth) GetTokenClaims(token string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, _, err := (&jwt.Parser{}).ParseUnverified(token, claims)
	return claims, err
}

// CreateGigantumTokens creates new gigantum tokens based on the user info from the original dex tokens
func (a *Auth) CreateGigantumTokens(userInfo userinfo.UserInfo, scopes state.Scopes, nonce string, serviceTokens bool) (state.Tokens, error) {
	tokens := state.Tokens{}

	// Create Access JWT
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": a.JWTIssuer,
		"aud": "HossServer",
		"sub": userInfo.Subject,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * time.Duration(a.Settings.TokenExpirationHours.Access)).Unix(),
		"azp": a.Config.ClientID,
	})

	jwtToken.Header["kid"] = a.KeyID

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := jwtToken.SignedString(a.JwtKey)
	if err != nil {
		return tokens, err
	}

	tokens.AccessToken = tokenString

	// Create Refresh JWT
	scopeList := []string{}
	if scopes.ScopeOpenID {
		scopeList = append(scopeList, "openid")
	}
	if scopes.ScopeProfile {
		scopeList = append(scopeList, "profile")
	}
	if scopes.ScopeEmail {
		scopeList = append(scopeList, "email")
	}
	if scopes.ScopeHOSS {
		scopeList = append(scopeList, "hoss")
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":      a.JWTIssuer,
		"aud":      "HossServer",
		"sub":      userInfo.Subject,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour * time.Duration(a.Settings.TokenExpirationHours.Refresh)).Unix(),
		"azp":      a.Config.ClientID,
		"nickname": strings.ToLower(userInfo.Username),
		"scopes":   strings.Join(scopeList, ","),
	})

	refreshToken.Header["kid"] = a.KeyID

	// Sign and get the complete encoded token as a string using the secret
	refreshTokenString, err := refreshToken.SignedString(a.JwtKey)
	if err != nil {
		return tokens, err
	}

	tokens.RefreshToken = refreshTokenString

	if scopes.ScopeOpenID {
		// Create ID JWT
		claims := jwt.MapClaims{
			"iss": a.JWTIssuer,
			"aud": "HossServer",
			"sub": userInfo.Subject,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Hour * time.Duration(a.Settings.TokenExpirationHours.Id)).Unix(),
		}
		if scopes.ScopeProfile {
			claims["name"] = userInfo.FullName
			claims["given_name"] = userInfo.GivenName
			claims["family_name"] = userInfo.FamilyName
			claims["nickname"] = strings.ToLower(userInfo.Username)
		}
		if scopes.ScopeEmail {
			claims["email"] = userInfo.Email
			claims["email_verified"] = userInfo.EmailVerified
		}
		if scopes.ScopeHOSS {
			claims["role"] = userInfo.Role
			groups, err := a.Database.ListUserGroupNames(userInfo.Username)
			if err != nil {
				return tokens, err
			}
			claims["groups"] = strings.Join(groups[:], ",")
		}
		if nonce != "" {
			claims["nonce"] = nonce
		}
		if serviceTokens {
			claims["nickname"] = userInfo.Subject // claim required for STS generation
			claims["role"] = "service"
			claims["groups"] = ""
			claims["service"] = true // DP ??? Just use the `service` role instead of this?
		}
		jwtToken = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

		jwtToken.Header["kid"] = a.KeyID

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err = jwtToken.SignedString(a.JwtKey)
		if err != nil {
			return tokens, err
		}

		tokens.IDToken = tokenString
	}

	return tokens, nil
}
