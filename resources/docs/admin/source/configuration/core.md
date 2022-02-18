# Core Service
The core service provides the primary REST API for a Hoss server. 

## Core Service Configuration

Configuring the core service is done via `~/.hoss/auth/config.yaml`. It's values are described below:

* `object_stores`: A list of `ObjectStore` items that describe an object store the Hoss can interface with. On boot, the server will create a database entry for any `ObjectStore` item that does not yet exist.
* `namespaces`: A list of `Namespace` items that describe namespaces that should be created on first boot.
* `queues`: A list of `Queue` items that are used by the sync service to monitor the related `ObjectStore`'s bucket events.
* `server`
  * dev: if `true`, CORS will be enabled for local frontend development and the API will run in development mode. If `false`, CORS will be disabled.
  * `auth_service`: The auth service endpoint. By default the internal Docker route is used. If using an auth service running in a different server, you must update this value.
  * `elasticsearch_endpoint`: The endpoint wher the Opensearch API is accessible. By default the internal Docker route is used. You should not have to modify this value.
  * `sync_frequency_minutes`: The rate at which the core service will query the auth service to syncronize user group information.


`ObjectStore` items contain the following fields:
* `name`: The name of the object store. This is how the store is referenced by other parts of the system.
* `description`: A description of this object store
* `type`: The type of object store. Currently this can be `minio` or `s3`
* `endpoint`: The endpoint where the object store API is available. When running minIO, this will be the root of the server (i.e. the `EXTERNAL_HOSTNAME` value in the `.env` file). When running S3, this should be `https://s3.amazonaws.com`
* `region`: (Optional) The region where the server connecting to the object store is running. This can be `null` when using minIO.
* `profile`: (Optional) The profile name in the `~/.hoss/core/aws_credentials` file. This can be `null` when using minIO.
* `role_arn`: (Optional) The ARN for the service account role that is used to assume users via STS. This can be `null` when using minIO.
* `notification_arn`: (Optional) The ARN for the SQS queue where bucket events will be sent. This can be `null` when using minIO.

`Namespace` items contain the following fields:
* `name`: The name of the namespace. This is how the namespace is referenced by other parts of the system and is visible to users in the Hoss UI.
* `description`: A description of this namespace
* `bucket`: The bucket name that this namespace uses
* `object_store`: The `ObjectStore` name that contains the bucket that this namespace uses


`Queue` items contain the following fields:
* `type`: The type of queue. Currently `amqp` and `sqs` are supported, with `amqp` being used by minIO and `sqs` being used by S3.
* `settings`: Settings are dependent on the type
  * If using an `amqp` queue
    * `url`: The URL used to connect to the amqp service
  * If using `sqs`
    * `queue_name`: The name of the FIFO queue used for API notifications
    * `region`: The region the queues are in
    * `profile`: The profile name in the `~/.hoss/core/aws_credentials` file used to connect to the queues.
* `object_store`: The `ObjectStore` name that this queue is used with


## Setting AWS Credentials
AWS credentials are provided to the core service via the `~/.hoss/core/aws_credentials`, which is bind mount into the service container. You should set the Hoss service account credentials in this file as shown below. You can use any profile name as long as you are sure to set it in all required config files.

```
[hoss-service-account]
aws_access_key_id = SDHFVMWJSD343ANSADNa
aws_secret_access_key = SDjfdsjSYwnd8*56$7s2hdsjdASF
region=us-east-1
```
