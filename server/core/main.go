package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/gigantum/hoss-core/pkg/api"
	"github.com/gigantum/hoss-core/pkg/config"
	"github.com/gigantum/hoss-core/pkg/database"
	"github.com/gigantum/hoss-core/pkg/opensearch"
	"github.com/gigantum/hoss-core/pkg/store"
	"github.com/gigantum/hoss-core/pkg/sync"
	"github.com/gigantum/hoss-core/pkg/worker"
	"github.com/golang-jwt/jwt/v4"

	auth "github.com/gigantum/hoss-service"

	"github.com/gigantum/hoss-core/docs"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @Title Hoss - Core Service API
// @description This is the primary API used by the Hoss system.
// @BasePath /core/v1
// @securityDefinitions.apikey BearerToken
// @in header
// @name Authorization
func main() {
	docs.SwaggerInfo.BasePath = "/core/v1"

	// Load config file from default location
	c := config.Load("")

	r := gin.Default() // Default config includes Logger and Recovery middlewares
	r.Use(ConfigMiddleware(c))

	if c.Server.Dev {
		r.Use(cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "DELETE", "HEAD", "PUT"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
			Debug:            true,
		}))
	}

	// Endpoints in the v1_public group do not enforce the AuthorizeJWT
	v1_public := r.Group("v1")
	{
		v1_public.GET("discover", api.Discover)

		// If running in dev mode, publicly host the swagger docs.
		if c.Server.Dev {
			v1_public.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		}
	}

	v1 := r.Group("v1")
	v1.Use(AuthorizeJWT(c))
	{
		// object stores
		v1.GET("object_store/", api.ListObjectStores)
		v1.GET("object_store/:object_store", api.GetObjectStore)

		// namespaces
		v1.GET("namespace/", api.ListNamespaces)
		v1.POST("namespace/", api.CreateNamespace)
		v1.GET("namespace/:namespace", api.GetNamespace)
		v1.DELETE("namespace/:namespace", api.DeleteNamespace)

		// namespace syncing
		v1.PUT("namespace/:namespace/sync", api.EnableSyncNamespace)
		v1.GET("namespace/:namespace/sync", api.GetSyncNamespace)
		v1.DELETE("namespace/:namespace/sync", api.DisableSyncNamespace)

		// datasets
		v1.POST("namespace/:namespace/dataset/", api.CreateDataset)
		v1.DELETE("namespace/:namespace/dataset/:name", api.DeleteDataset)
		v1.PUT("namespace/:namespace/dataset/:name/restore", api.RestoreDataset)
		v1.GET("namespace/:namespace/dataset/:name", api.GetDataset)
		v1.GET("namespace/:namespace/dataset/", api.ListDataset)

		// dataset syncing
		v1.GET("namespace/:namespace/dataset/:name/sync", api.GetSyncDataset)
		v1.PUT("namespace/:namespace/dataset/:name/sync", api.EnableSyncDataset)
		v1.DELETE("namespace/:namespace/dataset/:name/sync", api.DisableSyncDataset)

		// dataset permissions
		v1.PUT("namespace/:namespace/dataset/:name/user/:username/access/:accesslevel", api.UpdateUserDatasetPerms)
		v1.DELETE("namespace/:namespace/dataset/:name/user/:username", api.RemoveUserDatasetPerms)
		v1.PUT("namespace/:namespace/dataset/:name/group/:groupname/access/:accesslevel", api.UpdateGroupDatasetPerms)
		v1.DELETE("namespace/:namespace/dataset/:name/group/:groupname", api.RemoveGroupDatasetPerms)

		// metadata search
		v1.GET("search", api.SearchMetadata)
		v1.GET("search/namespace/:namespace/dataset/:name/key", api.SuggestKeys)
		v1.GET("search/namespace/:namespace/dataset/:name/key/:key/value", api.SuggestValues)
		v1.GET("search/namespace/:namespace/dataset/:name/metadata", api.GetMetadata)

		// credentials
		v1.GET("namespace/:namespace/sts", api.GetUserSTSCredentials)

		// core-auth group syncing
		v1.PUT("user/sync", api.SyncCurrentUserGroups)

		// Service Account only endpoints
		v1.GET("configuration/sync", api.GetSyncConfiguration)
		v1.GET("configuration/queue", api.GetNotificationQueues)
		v1.GET("object_store/:object_store/sts", api.GetServiceSTSCredentials)

		v1.PUT("search/document/metadata", api.CreateOrUpdateMetadataDocument)
		v1.DELETE("search/document/metadata", api.DeleteMetadataDocument)
	}

	r.Run() // listen and serve on 0.0.0.0:8080
}

// ConfigMiddleware is middleware to load config and store instances
func ConfigMiddleware(config *config.Configuration) gin.HandlerFunc {

	// Load the application configuration
	db := database.Load()

	// Bootstrap the default namespace if needed. noop if the namespace exists
	err := database.BootstrapDefaults(config, db)
	if err != nil {
		logrus.Errorf("Failed to bootstrap namespaces: %v", err)
	}

	// Bootstrap all object stores
	// Since we currently can only configure Object Stores via the config
	// file, it's safe to just list all the object stores in the database at this
	// point (they will be added in the call above if new), then load
	// the instances for each object store. The object store name must be the
	// same and unique between linked servers for syncing to work properly.
	objectStores, err := db.ListObjectStores(1000, 0)
	if err != nil {
		logrus.Errorf("Failed to list available object stores: %v", err)
	}
	s := store.LoadObjectStores(config, objectStores)

	// start syncing goroutine to keep core group memberships in sync with auth
	go api.SyncAllUserGroups(config, db, s)

	// Load Exchange for sending API sync messages
	ase, err := sync.LoadApiSyncExchange(config)
	if err != nil {
		logrus.Errorf("Failed to load API sync exchange")
	}

	// Patch all minio events
	// This function will list all datasets in a namespace that is backed by a minio
	// object store. Then it will disable/enable bucket events on the dataset
	// For some unknown reason this is required when running minio in gateway
	// mode on a container restart. This issue is tracked here:
	// https://github.com/minio/minio/issues/13816
	// If this is fixed upstream, we can remove this functionality.
	err = store.PatchMinioEvents(db, s)
	if err != nil {
		logrus.Errorf("Failed to patch minio events: %v", err.Error())
	}

	// Start the background dataset delete worker
	// This function will loop infinitely, waiting for datasets to be
	// ready for delete.
	exitCh := make(chan bool)
	go worker.DeleteDatasetWorker(config, db, s, exitCh)

	// Wait for opensearch to be ready
	for i := 0; i < 30; i++ {
		_, err = http.Get(config.Server.ElasticsearchEndpoint)
		if err == nil {
			break
		}

		logrus.Info("Opensearch service is not ready, sleeping")
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		logrus.WithField("error", err.Error()).Fatal("Couldn't connect to Opensearch service after 60 seconds")
	}

	// Create metadata search index if it has yet to be created
	err = opensearch.CreateMetadataSearchIndex(config.Server.ElasticsearchEndpoint)
	if err != nil {
		logrus.Errorf("Failed to initialize metadata search index: %v", err.Error())
	}

	return func(c *gin.Context) {
		c.Set("config", config)
		c.Set("stores", s)
		c.Set("db", db)
		c.Set("apiSyncExchange", ase)
		c.Next()
	}
}

// AuthorizeJWT enforces the use of a JWT Bearer token for auth
// Requires the JWT to contain a `nickname` claim
// If there is no bearer token or the token is invalid the this returns a 401 Unauthorized response to the user
func AuthorizeJWT(config *config.Configuration) gin.HandlerFunc {
	var oidcConfig auth.OpenIDConfiguration
	var err error
	for i := 0; i < 12; i++ {
		oidcConfig, err = auth.GetOpenIDConfiguration(config.Server.AuthService)
		if err == nil {
			break
		}

		logrus.Info("Auth service is not ready, sleeping")
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		logrus.WithField("error", err.Error()).Fatal("Couldn't get OpenID configuration after 60 seconds")
	}

	keys, err := auth.GetJWKS(oidcConfig.JwksUri)
	if err != nil {
		logrus.WithField("error", err.Error()).Fatal("Couldn't get JWKS information")
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		kid := token.Header["kid"].(string)
		key := keys.GetSigningKey(kid)
		if key == nil {
			return nil, errors.Wrap(err, "Couldn't locate key id "+kid)
		}

		public, err := key.GetKey()
		if err != nil {
			return nil, errors.Wrap(err, "Couldn't get public key")
		}

		return public, nil
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(parts[1], keyFunc)

		if err != nil {
			logrus.Errorf("Parse Error: %s", err.Error())
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			logrus.Errorf("Invalid JWT: %v", token.Claims)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userClaims := token.Claims.(jwt.MapClaims)

		// DP NOTE: by default Parse() will validate the exp claim

		iss := userClaims["iss"].(string)
		if iss != oidcConfig.Issuer {
			logrus.Errorf("Invalid JWT Issuer: %v", iss)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		aud := userClaims["aud"].(string)
		if aud != "HossServer" {
			logrus.Errorf("Invalid JWT Audience: %v", aud)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		isService := false
		service, ok := userClaims["service"]
		if ok && service.(bool) {
			isService = true
		}

		userInfo := api.UserInfo{
			JWT:       parts[1],
			Claims:    userClaims,
			Username:  userClaims["nickname"].(string),
			Role:      userClaims["role"].(string),
			Groups:    strings.Split(userClaims["groups"].(string), ","),
			IsService: isService,
		}

		c.Set("user", userInfo)
		c.Next()
	}
}
