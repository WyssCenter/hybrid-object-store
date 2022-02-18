# Dual Server Hybrid Cloud Installation

This section describes a common deployment scenario where one Hoss server is running on-premise
and one Hoss server is running AWS. The servers are then linked together via auth and sync
services. 

This enabled hybrid cloud workflows that can empower
individual users to optimize their time and costs while delivering new capabilities that
are typically out of reach for data scientists. Some examples of these capabilities are:

* Generating data at an off-site location and having it automatically available on-premise for storage and analysis
* Collaborating with external partners on data generated on-premise
* Delivering data generated on-premise to external parters
* Seamlessly switching between local, on-premise and cloud compute resources for analysis
* Discover data and track where data is located

The main thing to understand is the choices made here are to:

* Use S3 for object storage in the cloud
* Use minIO for object storage on-premise
* Use internal LDAP provider
* Use a shared auth configuration, with a single auth service running in the cloud
* Use a single sync service running on-premise
  * This is chosen because it removes the requirement for the cloud server to be able to reach the on-premise network. You could
    flip this so sync runs in the cloud if the local server is reachable from the cloud.
* Use Let's Encrypt
  * If you plan to run your server behind an ALB and use AWS provisioned certificates, review the [TLS Configuration](../configuration/tls.md) section.
* Enable reCaptcha
* Running an on-premise server with the FQDN `hoss-local.mycompany.com`
* Running a cloud server with the FQDN `hoss-cloud.mycompany.com`

To complete this deployment, you will follow the [Single On-premise Server](install-on-prem.md) and [Single AWS Server](install-aws.md) instructions
with modifications outlined below.

## Deploy AWS Server
You must deploy the AWS server first since it will be running the shared auth service. Follow the [Single AWS Server](install-aws.md) instructions to complete the deployment.

### Configure Variables
In the [Configure Variables](install-aws.md#configure-variables) section, you must modify the `SERVICES` variable to not include the `sync` service, since that will only be running on-premise. For example:

```SERVICES=opensearch ldap db dex reverse-proxy auth ldap-admin core ui```

You may also skip the [Configure the Sync Service](install-aws.md#configure-the-sync-service) section, since that service is now disabled.

### Configure Auth Service
In the [Configure the Auth Service](install-on-prem.md#configure-the-auth-service) section, when editing
`~/.hoss/auth/config.json` you must also set `additional_allowed_servers` to include your on-premise
server. For the example where the local server is hosted at http://hoss-local.mycompany.com, then you would set:

```
additional_allowed_servers:
  - http://hoss-local.mycompany.com
```

This will allow redirecting back to the on-premise server after successfully logging in.

## Deploy On-Premise Server
You must deploy the on-premise server second since it will be using the shared auth service that is
running in the cloud. Follow the [Single On-premise Server](install-on-prem.md) instructions to complete the deployment with the modifications outlined below.

### Configure Variables
In the [Configure Variables](install-on-prem.md#configure-variables) section, you must modify the `AUTH_SERVICE_ENDPOINT` variable to point to your cloud server. So assuming the cloud server's FQDN
is `hoss-cloud.mycompany.com` and TLS is enabled via Let's Encrypt, then you would set `AUTH_SERVICE_ENDPOINT=https://hoss-cloud.mycompany.com/auth/v1`.

You must also set the `SERVICE_AUTH_SECRET` to the same value as the cloud server. View the `.env` file in
your cloud server and make note of the `SERVICE_AUTH_SECRET` value. Then set `SERVICE_AUTH_SECRET` equal
to this value in your on-premise server. This is the PAT that the service account will use to get credentials when making API calls to the servers.

Finally, you must set the `SERVICE` variable to not run auth resources, for example:

```shell
SERVICES=opensearch rabbitmq db reverse-proxy etcd-0 minio core ui sync
```

### Configure Core Service
In the [Configure the Core Service](install-on-prem.md#configure-the-core-service) section, when editing
`~/.hoss/core/config.json` you must also set `server.auth_service` to the same value that you set `AUTH_SERVICE_ENDPOINT` to in the `.env` file. So assuming the cloud server's FQDN
is `hoss-cloud.mycompany.com` and TLS is enabled via Let's Encrypt, then you would set it to `https://hoss-cloud.mycompany.com/auth/v1`.

You must add the object store and SQS queue used in the cloud server so that it is available to the sync service. To do this, add an entries the `object_stores` and `queues` items in `~/.hoss/core/config.json`. For example:

```yaml
object_stores:
  - name: default
    description: Default object store
    type: minio
    endpoint: http://hoss-local.mycompany.com
    region: null
    profile: null
    role_arn: null
    notification_arn: null
  - name: s3
    description: S3 object store
    type: s3
    endpoint: https://s3.amazonaws.com
    region: <REGION>
    profile: <PROFILE>
    role_arn: "arn:aws:iam::<ACCOUNT_ID>:role/<HOSS_ASSUME_USER_ROLE>"
    notification_arn: "arn:aws:sqs:us-east-1:<ACCOUNT_ID>:<QUEUE_NAME>"
```

and

```yaml
queues:
  - type: amqp
    settings:
      url: amqp://${RABBITMQ_USER}:${RABBITMQ_PASS}@rabbitmq:5672
    object_store: default
  - type: sqs
    settings:
      queue_name: <FIFO_QUEUE_NAME>
      region: <REGION>
      profile: <PROFILE>
    object_store: s3
```

Where the variables shown are the same as those in the [Configure the Core Service](install-aws.md#configure-the-core-service) section of the AWS server install instructions.

### Configure Sync Service
In the [Configure the Sync Service](install-on-prem.md#configure-the-sync-service) section, first make sure both core services are listed. For example:

```yaml
core_services:
  - https://hoss-cloud.mycompany.com/core/v1
  - http://hoss-local.mycompany.com/core/v1
```

Set the `auth_endpoint` to the same value that you set `AUTH_SERVICE_ENDPOINT` to in the `.env` file. So assuming the cloud server's FQDN is `hoss-cloud.mycompany.com` and TLS is enabled via Let's Encrypt, then you would set it to `https://hoss-cloud.mycompany.com/auth/v1`.

Since the on-premise sync service will need to make requests to the cloud server's S3 bucket, make sure `sqs_profile` is equal to the service account profile used in the `~/.hoss/sync/aws_credentials` file below. In this example, it would be set to `hoss-service-account`

Finally, edit `~/.hoss/sync/aws_credentials` and set a profile for the service account you either created manually or via Terraform. E.g.

```
[hoss-service-account]
aws_access_key_id = SDHFVMWJSD343ANSADNa
aws_secret_access_key = SDjfdsjSYwnd8*56$7s2hdsjdASF
region=us-east-1
```

