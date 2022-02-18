module github.com/gigantum/hoss-sync

go 1.14

replace (
	github.com/gigantum/hoss-error => ../libs/hoss-error
	github.com/gigantum/hoss-service => ../libs/hoss-service
)

require (
	github.com/aws/aws-sdk-go-v2 v1.7.1
	github.com/aws/aws-sdk-go-v2/config v1.3.0
	github.com/aws/aws-sdk-go-v2/credentials v1.2.1
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.2.3
	github.com/aws/aws-sdk-go-v2/service/s3 v1.10.0
	github.com/aws/aws-sdk-go-v2/service/sqs v1.7.0
	github.com/ghodss/yaml v1.0.0
	github.com/gigantum/hoss-error v0.0.0-00010101000000-000000000000
	github.com/gigantum/hoss-service v0.0.0-00010101000000-000000000000
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/streadway/amqp v1.0.0
)
