module github.com/gigantum/hoss-core

go 1.14

replace (
	github.com/gigantum/hoss-error => ../libs/hoss-error
	github.com/gigantum/hoss-service => ../libs/hoss-service
)

require (
	github.com/aws/aws-sdk-go-v2 v1.7.1
	github.com/aws/aws-sdk-go-v2/config v1.3.0
	github.com/aws/aws-sdk-go-v2/service/iam v1.5.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.8.0
	github.com/aws/aws-sdk-go-v2/service/sqs v1.7.0
	github.com/aws/aws-sdk-go-v2/service/sts v1.4.1
	github.com/gigantum/hoss-service v0.0.0-00010101000000-000000000000
	github.com/gin-gonic/gin v1.7.7
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/go-pg/migrations/v8 v8.0.1
	github.com/go-pg/pg/v10 v10.7.7
	github.com/go-playground/validator/v10 v10.10.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.1.0
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/minio/minio-go/v7 v7.0.10
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.7.0
	github.com/sirupsen/logrus v1.8.1
	github.com/streadway/amqp v1.0.0
	github.com/swaggo/files v0.0.0-20210815190702-a29dd2bc99b2
	github.com/swaggo/gin-swagger v1.4.0
	github.com/swaggo/swag v1.7.8
	github.com/ugorji/go v1.2.6 // indirect
	golang.org/x/crypto v0.0.0-20220131195533-30dcbda58838 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20220128215802-99c3d69c2c27 // indirect
	golang.org/x/tools v0.1.9 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v2 v2.4.0
)
