module github.com/gigantum/hoss-auth

go 1.14

replace (
	github.com/gigantum/hoss-error => ../libs/hoss-error
	github.com/gigantum/hoss-service => ../libs/hoss-service
)

require (
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/gigantum/hoss-error v0.0.0-00010101000000-000000000000
	github.com/gigantum/hoss-service v0.0.0-00010101000000-000000000000
	github.com/gin-gonic/gin v1.7.7
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/go-pg/migrations/v8 v8.0.1
	github.com/go-pg/pg/v10 v10.7.7
	github.com/go-playground/validator/v10 v10.9.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.1.0
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/pkg/errors v0.9.1
	github.com/pquerna/cachecontrol v0.0.0-20201205024021-ac21108117ac // indirect
	github.com/rs/cors v1.7.0
	github.com/sirupsen/logrus v1.8.1
	github.com/swaggo/files v0.0.0-20210815190702-a29dd2bc99b2
	github.com/swaggo/gin-swagger v1.4.0
	github.com/swaggo/swag v1.7.8
	github.com/ugorji/go v1.2.6 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/oauth2 v0.0.0-20210313182246-cd4f82c27b84
	golang.org/x/sys v0.0.0-20220128215802-99c3d69c2c27 // indirect
	golang.org/x/tools v0.1.9 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
	gopkg.in/yaml.v2 v2.4.0
)
