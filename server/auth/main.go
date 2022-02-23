package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gigantum/hoss-auth/pkg/api"
	"github.com/gigantum/hoss-auth/pkg/userinfo"
	hosserror "github.com/gigantum/hoss-error"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/sirupsen/logrus"

	"github.com/gigantum/hoss-auth/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @Title Hoss - Auth Service API
// @description This is the auth service API used by the Hoss system to manage users, groups, and tokens
// @securityDefinitions.apikey BearerToken
// @in header
// @name Authorization
func main() {
	docs.SwaggerInfo.BasePath = "/auth/v1"

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery()) // handle panics

	auth := api.LoadConfig()
	if auth.Settings.DevServer {
		log.Println("Dev mode - Allowing all origins")
		r.Use(cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "DELETE", "HEAD", "PUT"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
			Debug:            true,
		}))
	} else if len(auth.Settings.AdditionalAllowedServers) > 0 {
		log.Println("Additional Servers allowed. Enabling CORS for additional servers")
		r.Use(cors.New(cors.Options{
			AllowedOrigins:   auth.Settings.AdditionalAllowedServers,
			AllowedMethods:   []string{"GET", "POST", "DELETE", "HEAD", "PUT"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
			Debug:            false,
		}))
	}

	// Endpoints in the v1_public group do not enforce the AuthorizeJWT
	v1_public := r.Group("v1")
	{
		v1_public.GET("ping", auth.Ping)

		// OIDC Provider Endpoints
		v1_public.GET(".well-known/openid-configuration", auth.OpenID)
		v1_public.GET("keys", auth.JWKS)
		v1_public.GET("authorize", auth.Authorize)
		v1_public.GET("callback", auth.Callback)
		v1_public.POST("token", auth.Tokens)

		// If running in dev mode, publicly host the swagger docs.
		if auth.Settings.DevServer {
			v1_public.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		}
	}

	v1 := r.Group("v1")
	v1.Use(AuthorizationMiddleware(auth))
	{
		// HOSS Token Endpoints
		v1.POST("pat/exchange/jwt", auth.PATforJWT)
		v1.POST("pat", auth.JWTForPAT)
		v1.GET("pat/", auth.ListPAT)
		v1.DELETE("pat/:id", auth.DeletePAT)

		v1.POST("sat/exchange/jwt", auth.SATforJWT)

		// Group Endpoints
		v1.POST("group/", auth.CreateGroup)
		v1.DELETE("group/:groupname", auth.DeleteGroup)
		v1.GET("group/:groupname", auth.GetGroup)
		v1.GET("user/:username", auth.GetUser)
		v1.PUT("group/:groupname/user/:username", auth.UpdateGroupUser)
		v1.DELETE("group/:groupname/user/:username", auth.RemoveGroupUser)
		v1.GET("usernames", auth.GetUsernames)

		// OIDC endpoint that requires a JWT
		v1.GET("userinfo", auth.UserInfo)

		// Endpoints to change user passwords if running internal auth provider
		v1.GET("password", auth.ChangePasswordSupported)
		v1.PUT("password", auth.ChangePassword)

		// Endpoint to remove user tokens and permissions
		v1.DELETE("user/:username", auth.DeleteUser)
	}

	r.Run() // listen and serve on 0.0.0.0:8080

	log.Println("Skipping cert validation: ", auth.Settings.DevServer)
	log.Println("Client is running at port 8080")
	//log.Fatal(http.ListenAndServe(":8080", nil))
}

// AuthorizationMiddleware enforces the use of a JWT Bearer token or personal access token
// Th JWT must contain a `nickname` claim.
// If there is no bearer token, not PAT, or the token is invalid, this returns a 401 Unauthorized response
func AuthorizationMiddleware(config *api.Auth) gin.HandlerFunc {

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// The Authorization header is not set
			logrus.Warning("No Authorization header provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			logrus.Warning("Malformed Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Note: the prefix `hp_` is used to identify the token in the Authorization header and provide support for adding
		// future capabilities to tokens (i.e. the system would be able to infer the version of the token format)
		// This prefix is set in the database migrations and tokens are automatically generated on record insertion.
		if strings.HasPrefix(parts[1], "hp_") {
			// This is a PAT
			userInfo, err := authorizePAT(parts[1], config)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				return
			}

			c.Set("pat", parts[1])
			c.Set("user", userInfo)
		} else if strings.HasPrefix(parts[1], "hsvc_") {
			// This is a service token
			if parts[1] != os.Getenv("SERVICE_AUTH_SECRET") {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				return
			}

			// The API verifies a request was made with a service token by checking for this key to exist.
			c.Set("sat", parts[1])
		} else {
			// This is a likely a JWT
			userInfo, err := authorizeJWT(parts[1], config)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				return
			}

			c.Set("jwt", parts[1])
			c.Set("user", userInfo)
		}
	}
}

// authorizeJWT enforces the use of a JWT Bearer token for auth
// Requires the JWT to contain a `nickname` claim
// If there is no bearer token or the token is invalid the this returns a 401 Unauthorized response to the user
func authorizeJWT(token string, config *api.Auth) (userinfo.UserInfo, error) {
	userClaims, err := config.ValidateToken(token)
	if err != nil {
		logrus.Errorf("Invalid JWT: %v", err)
		return userinfo.UserInfo{}, hosserror.ErrUnauthorized
	}

	// NOTE: by default Parse() will validate the exp claim
	iss := userClaims["iss"].(string)
	if iss != config.JWTIssuer {
		logrus.Errorf("Invalid JWT Issuer: %v", iss)
		return userinfo.UserInfo{}, hosserror.ErrUnauthorized
	}

	aud := userClaims["aud"].(string)
	if aud != "HossServer" {
		logrus.Errorf("Invalid JWT Audience: %v", iss)
		return userinfo.UserInfo{}, hosserror.ErrUnauthorized
	}

	// fill out user info from claims
	userInfo, errClaims := userinfo.ParseUserClaims(userClaims,
		config.Settings.AdminGroup, config.Settings.PrivilegedGroup, config.Settings.UsernameClaim)
	if errClaims != nil {
		logrus.Error("Unable to extract user info claims")
		return userinfo.UserInfo{}, hosserror.ErrUnauthorized
	}

	return userInfo, nil
}

// authorizePAT enforces the use of a PAT to make requests
func authorizePAT(token string, config *api.Auth) (userinfo.UserInfo, error) {
	patRecord, err := config.Database.GetPAT(token)
	if err != nil {
		logrus.Errorf("Error verifying PAT: %s", err)
		return userinfo.UserInfo{}, hosserror.ErrUnauthorized
	}

	userInfo := userinfo.UserInfo{
		Subject:       patRecord.Owner.Subject,
		FullName:      patRecord.Owner.FullName,
		GivenName:     patRecord.Owner.GivenName,
		FamilyName:    patRecord.Owner.FamilyName,
		Username:      patRecord.Owner.Username,
		Email:         patRecord.Owner.Email,
		EmailVerified: &patRecord.Owner.EmailVerified,
		Role:          patRecord.Owner.Role,
	}

	return userInfo, nil
}
