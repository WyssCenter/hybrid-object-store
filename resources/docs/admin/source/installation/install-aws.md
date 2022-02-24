# Single AWS Server

This section describes setting up a single server in a typical configuration for an AWS EC2 instance.
The main thing to understand is the choices made here are to:

* Use S3 for object storage
* Use internal LDAP provider
* Use Let's Encrypt
  * If you plan to run your server behind an ALB and use AWS provisioned certificates, review the [TLS Configuration](../configuration/tls.md) section.
* Enable reCaptcha
* Running a server with the FQDN `hoss.mycompany.com`

## Create AWS Resources
Since we are using S3 for object storage, in addition to the instance there are various AWS resources that must be created.

There is a [Terraform](https://www.terraform.io/) module available to help manage this - [https://github.com/WyssCenter/terraform-hoss-aws](https://github.com/WyssCenter/terraform-hoss-aws).
This is the recommended way to do things since it is much easier and less prone to error. Simply configure and deploy this module and all the variables
you need to finish configuration of the system will be output.

Alternatively, you can manually configure resources. If you wish to instead manually create all the resources, you must follow the following steps.

### Create A Bucket 

1. Go to the S3 console
2. Click on the "Create bucket" button
3. Give the bucket a name and use default settings
4. To support the web UI file browser, you must enable CORS on the bucket for your server. 
   1. Navigate to the bucket in the S3 console
   2. Open the "Permissions" tab
   3. In the Cross-origin resource sharing (CORS) section, enter the following policy, swapping in your
      servers FQDN. If running from localhost (e.g. for development) you must include that in the "AllowedOrigins"

      ```json
        [
            {
                "AllowedHeaders": [
                    "*"
                ],
                "AllowedMethods": [
                    "HEAD",
                    "GET",
                    "POST",
                    "PUT",
                    "DELETE"
                ],
                "AllowedOrigins": [
                    "https://hoss.mycompany.com"
                ],
                "ExposeHeaders": [
                    "ETag",
                    "x-amz-meta-custom-header",
                    "x-amz-server-side-encryption",
                    "x-amz-request-id",
                    "x-amz-id-2",
                    "date"
                ],
                "MaxAgeSeconds": 3000
            }
        ]
      ```

### Create an IAM Policy for the role that is used to generate temporary credentials for users

1. Go to the IAM console and select the "Policies"
2. Click on the "Create policy" button
3. Enter the following JSON policy, updating <BUCKET_NAME> with the bucket created above

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "HossUserRolePolicy",
            "Effect": "Allow",
            "Action": [
                "s3:AbortMultipartUpload",
                "s3:DeleteObject",
                "s3:DeleteObjectVersion",
                "s3:GetBucketAcl",
                "s3:GetBucketCORS",
                "s3:GetBucketNotification",
                "s3:GetBucketPolicy",
                "s3:GetEncryptionConfiguration",
                "s3:GetObject",
                "s3:GetObjectAcl",
                "s3:GetObjectVersion"
                "s3:GetObjectVersionAcl",
                "s3:ListBucket",
                "s3:ListBucketMultipartUploads",
                "s3:ListBucketVersions",
                "s3:ListMultipartUploadParts",
                "s3:PutObject",
            ],
            "Resource": [
                "arn:aws:s3:::<BUCKET_NAME>/*",
                "arn:aws:s3:::<BUCKET_NAME>"
            ]
        }
    ]
}
```

4. Give the policy a name (e.g. `hoss-user-assume-role-policy`) and description


### Create a IAM Role for generating temporary credentials via AssumeRole

1. Go to the IAM console and select the "Roles" section.
2. Click on the "Create role" button
3. Select Another AWS Account
4. Enter the account ID for your current account
5. Select the policy created in the previous step (e.g. `hoss-user-assume-role-policy`)
6. Give the role a name (e.g. `hoss-user-assume-role`) and description
7. After creation, set the "Maximum session duration" to 12 hours.

### Create a IAM Policy for the HOSS Service Acount IAM User

1. Go to the IAM console and selec the "Policies"
2. Click on the "Create policy" button
3. Enter the following JSON policy, updating <ACCOUNT_ID> with the account id and <BUCKET> with the bucket

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "ServiceAccountCore",
            "Effect": "Allow",
            "Action": [
                "iam:CreatePolicy",
                "iam:CreatePolicyVersion",
                "iam:DeletePolicy",
                "iam:DeletePolicyVersion"
                "iam:GetPolicy",
                "iam:GetPolicyVersion",
                "iam:ListPolicyVersions",
                "iam:TagPolicy",
                "s3:DeleteObject",
                "s3:GetBucketNotification",
                "s3:GetObject",
                "s3:ListBucket",
                "s3:PutBucketNotification",
                "s3:PutObject",
                "sqs:DeleteMessage",
                "sqs:GetQueueAttributes",
                "sqs:GetQueueUrl",
                "sqs:ListDeadLetterSourceQueues",
                "sqs:ListQueueTags",
                "sqs:ReceiveMessage",
                "sqs:SendMessage",
            ],
            "Resource": [
                "arn:aws:s3:::<BUCKET>",
                "arn:aws:s3:::<BUCKET>/*"
                "arn:aws:iam::<ACCOUNT_ID>:policy/hoss-user-policy-*",
                "arn:aws:sqs:*:<ACCOUNT_ID>:hoss-*",
            ]
        },
        {
            "Sid": "ServiceAccountCoreAssumeRole",
            "Effect": "Allow",
            "Action": [
                "sts:AssumeRole"
            ],
            "Resource": [
                "arn:aws:iam::<ACCOUNT_ID>:role/hoss-user-assume-role"
            ]
        },
        {
            "Sid": "ServiceAccountSyncListResources",
            "Effect": "Allow",
            "Action": [
                "iam:ListPolicies",
                "sqs:ListQueues",
                "sts:GetCallerIdentity"
            ],
            "Resource": "*"
        }
    ]
}
```
4. Give the policy a name (e.g. `hoss-service-account-policy`) and description


### Create a IAM User for the HOSS Service Account

1. Go to the IAM console
2. Click on the "Add user" button
3. Set a username (e.g. `hoss-service-account`) and check "Programmatic access"
4. Click "Next: Permissions"
5. Click "Attach existing policies directly" and select the policy created in the previous step (e.g. `hoss-service-account-policy`)
6. Download the creds file and use this to add a profile to the `~/.hoss/core/aws_credentials` file (and `~/.hoss/sync/aws_credentials` file if using the sync service). This is the profile (i.e. `hoss-service-account`) that you need to set when creating an object store that connects to S3 using this account. e.g.:

```
[hoss-service-account]
aws_access_key_id = <ACCESS KEY ID>
aws_secret_access_key = <SECRET ACCESS KEY>
region=us-east-1
```

### Create SQS Queues for Sync
Two SQS queues need to be created, one for S3 bucket notifications and one for Hoss API notifications. 

Run the following instructions to create the two queues (e.g. hoss-bucket-notifications and hoss-api-notifications.fifo)
1. Go to the SQS console
2. Click on the "Create Queue" button
3. Select the Fifo type queue (for API notifications) or leave it as Standard type (for bucket notifications)
4. Set the queue name (making sure the API notification queue name ends in `.fifo`)
5. Set the queue access policy. For Hoss API notifications (the fifo queue) use the default values. For S3 bucket notifications use the following advanced policy, updating `<REGION>`, `<ACCOUNT_ID>`, `<QUEUE_NAME>`, and `<BUCKET_NAME>` with the appropriate values.

    ```json
    {
        "Version": "2012-10-17",
        "Statement": [
        {
            "Sid": "S3-notifications",
            "Effect": "Allow",
            "Principal": {
                "Service": "s3.amazonaws.com"
            },
            "Action": [
                "SQS:SendMessage"
            ],
            "Resource": "arn:aws:sqs:<REGION>:<ACCOUNT_ID>:<QUEUE_NAME>",
            "Condition": {
                "ArnLike": { "aws:SourceArn": "arn:aws:s3:*:*:<BUCKET_NAME>" },
                "StringEquals": { "aws:SourceAccount": "<ACCOUNT_ID>" }
            }
        }
        ]
    }
    ```
6. Click the "Create Queue" button

Once the two queues have been created the [core](../configuration/core.md) and [sync](../configuration/sync.md) config.yaml files need to be updated to reference the queues
and AWS profile.


## Set Up Repository
The Hoss system is installed and managed directly from the source repository. The version of this repository will determine the
version of the server that is run. You can checkout tags to deploy specific releases (e.g. `git checkout 0.2.5`).

The system also uses a working directory `~/.hoss` contains all configuration information and data for the running server. This
location is created and pre-populated via the Makefile.

To get started, clone the repository:

```shell
git clone https://github.com/WyssCenter/hybrid-object-store.git
```

This will install the latest version of the server. To install a known release, review the [Releases](https://github.com/WyssCenter/hybrid-object-store/releases) 
page and make note of the desired version (e.g. 0.2.5). Then check out that version via `git checkout <version>`.

## Configure Variables
The Hoss uses an environment variables file to maintain common and sensitive information. These data are used to configure various
parts of the system, and in the future in more advanced deployment scenarios could be moved to other secure means of managing environment variables.

You should review the complete description of what can be configured in the [Environment Variables section](../configuration/env-vars.md).

For this scenario we'll use mostly defaults and set a few items. Assuming we are logged into an ubuntu server with the username `ubuntu`,
we first initialize the working directory and `.env` file.

```
cd server
make env
```

Next, we edit the file `/home/ubuntu/.hoss/.env`, changing only the following items:

```
SERVICES=opensearch ldap db dex reverse-proxy auth ldap-admin core ui sync
LETS_ENCRYPT_ENABLED=true
EXTERNAL_HOSTNAME=https://hoss.mycompany.com
DOMAIN=hoss.mycompany.com
ADMIN_EMAIL=admin@mycompany.com
LDAP_ORGANISATION=My Company Inc.
LDAP_DOMAIN=hoss.mycompany.com
RECAPTCHA_SITE_KEY=examplekey1234
RECAPTCHA_SECRET_KEY=examplekey1234
```

Finally, run `make config` to finish configuring all services with base settings.

```{note}
When using Let's Encrypt, your server must be reachable on port 80 for the ACME challenge to succeed. The server
will automatically redirect all requests (except for the ACME challenge) to port 443/https.
```

## Configure the Core Service
There is additional configuration to the core service possible depending on your deployment scenario. You can review 
all possible items in the [Core Service](../configuration/core.md) configuration section. 

For this scenario, edit `~/.hoss/core/config.json`. 

You need to configure the `object_store`, `namespace`, and `queue` entries for the S3 bucket you either
created manually or via Terraform. A typical configuration when using Terraform is shown below, replacing
the variables with your values.

- `<ACCOUNT_ID>`: The AWS account ID where the Hoss is deployed
- `<BUCKET_NAME>`: The name of the S3 bucket you created via Terraform or manually
- `<REGION>`: The region your server is deployed
- `<PROFILE>`: The name of the profile in the `~/.hoss/core/aws_credentials` file. In this example, it would be set to `hoss-service-account`
- `<HOSS_ASSUME_USER_ROLE>`: The name of the role used to generate STS credentials. If you used the Terraform, this would be set to `hoss-user-assume-role`
- `<QUEUE_NAME>`: This is the name of the bucket event SQS queue. If using the terraform, this is likely `hoss-bucket-notifications-<BUCKET_NAME>`
- `<FIFO_QUEUE_NAME>`: The name of the API notification FIFO SQS queue. If using the terraform this is likely `hoss-api-notification.fifo`


```
object_stores:
  - name: s3
    description: S3 object store
    type: s3
    endpoint: https://s3.amazonaws.com
    region: <REGION>
    profile: <PROFILE>
    role_arn: "arn:aws:iam::<ACCOUNT_ID>:role/<HOSS_ASSUME_USER_ROLE>"
    notification_arn: "arn:aws:sqs:us-east-1:<ACCOUNT_ID>:<QUEUE_NAME>"
namespaces:
  - name: default
    description: Default namespace
    bucket: <BUCKET_NAME>
    object_store: s3
queues:
  - type: sqs
    settings:
      queue_name: <FIFO_QUEUE_NAME>
      region: <REGION>
      profile: <PROFILE>
    object_store: s3
```

Set `server.dev` to `false`

Finally, edit `~/.hoss/core/aws_credentials` and set a profile for the service account you either created manually or via Terraform. E.g.

```
[hoss-service-account]
aws_access_key_id = SDHFVMWJSD343ANSADNa
aws_secret_access_key = SDjfdsjSYwnd8*56$7s2hdsjdASF
region=us-east-1
```

## Configure the Sync Service
There is additional configuration to the sync service possible depending on your deployment scenario. You can review all 
possible items in the [Sync Service](../configuration/sync.md) configuration section. 

For this scenario, edit `~/.hoss/sync/config.json`. Set `core_services` to include this server's
external hostname. For example, assuming the server is running at `hoss.mycompany.com`:

```yaml
core_services:
  - https://hoss.mycompany.com/core/v1
```

Also, make sure `sqs_profile` is equal to the service account profile used in the `~/.hoss/sync/aws_credentials` file below. In this example, it would be set to `hoss-service-account`

Finally, edit `~/.hoss/sync/aws_credentials` and set a profile for the service account you either created manually or via Terraform. E.g.

```
[hoss-service-account]
aws_access_key_id = SDHFVMWJSD343ANSADNa
aws_secret_access_key = SDjfdsjSYwnd8*56$7s2hdsjdASF
region=us-east-1
```

## Configure the Auth Service
There is additional configuration to the auth service possible depending on your deployment scenario. You can review all 
possible items in the [Auth Service](../configuration/auth.md) configuration section. 

For this scenario, edit `~/.hoss/auth/config.json`. 

Set `server.dev` to `false`

Optionally modify the `password_policy` as desired. These critera are enforced on users when they change their password via the Hoss UI.

## Configure the UI Service
There is additional configuration to the UI service possible depending on your deployment scenario. You can review all possible items in the [UI Service](../configuration/ui.md) configuration section. 

For this deployment, set the "name" of the server that is presented to users in the menu bar. Edit `~/.hoss/ui/config.json` and set the `server_name` value to what ever you wish, e.g. "My Company Hoss Server - Cloud"

## Build Images
Hoss managed containers are simply built from the repository. Run `make build` to run any necessary build process.

## Start the Server
Finally, run `make up` to start the server. This will run the server in the foreground, showing all log output. This can
be useful when first setting up a server. 

Typically, you should instead run `make up DETACH=true`, which will run the server in the background. If you then
want to view log outputs, you can as described in the [monitoring logs section](../maintenance/monitor-logs.md).

```{warning}
When using the internal auth provider, default test accounts are created. **YOU MUST REMOVE THEM TO SECURE THE SERVER.** Review the Internal LDAP Server section for details.
```

## Stop the Server
To stop the server, run `make down`. 

If you are trying to reset a server (e.g. during development or to do a restore), make sure
to run `make down` first. If containers are running while you try to reset, you'll likely end up in a broken state where services
expect different credentials.