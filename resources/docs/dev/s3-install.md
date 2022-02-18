

## Create a bucket 

1. Go to the S3 console
2. Click on the "Create bucket" button
3. Give the bucket a name and use default settings 


## Create a IAM Policy the role that is used to generate temporary credentials

1. Go to the IAM console and selec the "Policies"
2. Click on the "Create policy" button
3. Enter the following JSON policy, updating <BUCKET_NAME> with the bucket created above

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "HOSSUserRolePolicy",
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

4. Create a tag with the Key `HOSS-Resource` and Value `user-assume-role-policy`
5. Give the policy a name (e.g. `hoss-user-assume-role-policy`) and description


## Create a IAM Role for generating temporary credentials via AssumeRole

1. Go to the IAM console and select the "Roles" section.
2. Click on the "Create role" button
3. Select Another AWS Account
4. Enter the account ID for your current account
5. Select the policy created in the previous step (e.g. `hoss-user-assume-role-policy`)
6. Create a tag with the Key `HOSS-Resource` and Value `hoss-user-assume-role`
7. Give the role a name (e.g. `hoss-user-assume-role`) and description
8. After creation, set the "Maximum session duration" to 12 hours.



## Create a IAM Policy for the HOSS Service Acount IAM User

1. Go to the IAM console and selec the "Policies"
2. Click on the "Create policy" button
3. Enter the following JSON policy, updating <ACCOUNT_ID> with the account id and <BUCKET> with the bucket

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
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
                "arn:aws:iam::<ACCOUNT_ID>:policy/*",
                "arn:aws:s3:::<BUCKET>",
                "arn:aws:s3:::<BUCKET>/*"
                "arn:aws:sqs:*:<ACCOUNT_ID>:*",
            ]
        },
        {
            "Sid": "VisualEditor1",
            "Effect": "Allow",
            "Action": [
                "sts:AssumeRole"
            ],
            "Resource": [
                "arn:aws:iam::<ACCOUNT_ID>:role/hoss-user-assume-role"
            ]
        },
        {
            "Sid": "VisualEditor2",
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
4. Create a tag with the Key `HOSS-Resource` and Value `service-account-policy`
5. Give the policy a name (e.g. `hoss-service-account-policy`) and description


## Create a IAM User for the HOSS Service Account

1. Go to the IAM console
2. Click on the "Add user" button
3. Set a username (e.g. `hoss-service-account`) and check "Programmatic access"
4. Click "Next: Permissions"
5. Click "Attach existing policies directly" and select the policy created in the previous step (e.g. `hoss-service-account-policy`)
6. Create a tag with the Key `HOSS-Resource` and Value `service-account`
7. Download the creds file and use this to add a profile to the `~/.hoss/core/aws_credentials` file (and `~/.hoss/sync/aws_credentials` file if using the sync service). This is the profile (i.e. `hoss-demo`) that you need to set when creating an object store that connects to S3 using this account. e.g.:

```
[hoss-demo]
aws_access_key_id = <ACCESS KEY ID>
aws_secret_access_key = <SECRET ACCESS KEY>
region=us-east-1
```

## (Optional) Create SQS Queues for Sync
If the deployment is going to be a source for syncing then two SQS queues need to be created, one for S3 bucket notifications and one for HOSS API notifications. Both queues are created the same, except for the queue name and queue access policy.

Run the following instructions to create the two queues (e.g. hoss-bucket-notifications and hoss-api-notifications.fifo)
1. Go to the SQS console
2. Click on the "Create Queue" button
3. Select the Fifo type queue (for API notifications) or leave it as Standard type (for bucket notifications)
4. Set the queue name (making sure the API notification queue name ends in `.fifo`)
5. For API notifications only, enable "Content based deduplication"
6. Set the queue access policy. For HOSS API notifications use the default values. For S3 bucket notifications use the following advanced policy, updating `<REGION>`, `<ACCOUNT_ID>`, `<QUEUE_NAME>`, and `<BUCKET_NAME>` with the appropriate values.
```
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
7. Create a tag with the Key `HOSS-Resource` and Value `sync-notification-queue`
8. Click the "Create Queue" button

Once the two queues have been created the core and sync config.yaml files need to be updated to reference the queues.
