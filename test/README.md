This directory contains tests and scripts for running integration level tests using the Python client library.


## Setup
To get started, first get a Python virtualenv setup to run the library and test code.

- create a virtualenv
- Install the client library from the source repo

```
cd client
pip install -U .
```

- Install testing deps from requirements.txt
- 
```
cd test
pip install -r requirements.txt
```

## Setting up a Hoo server to run S3 tests
A Hoss server must be configured for S3 to run the S3 tests and running. 

The int tests will use the aws profile `hoss-int-test`, so you must add creds with that profile in `~/.hoss/core/aws_credentials`
with the credentials for the Hoss service account.

```
[hoss-int-test]
aws_access_key_id = <ACCESS KEY ID>
aws_secret_access_key = <SECRET ACCESS KEY>
region=us-east-1
```

Next, you must add the S3 object store and namespace to your `~/.hoss/core/config.yaml` file. The full file
including the default minio configuration should look something like:

```
object_stores:
  - name: default
    description: Default object store
    type: minio
    endpoint: http://localhost
    region: null
    profile: null
    role_arn: null
  - name: s3
    description: S3 Object Store
    type: s3
    endpoint: https://s3.amazonaws.com
    region: us-east-1
    profile: hoss-int-test
    role_arn: arn:aws:iam::363182369550:role/hoss-user-assume-role
    notification_arn: "arn:aws:sqs:us-east-1:363182369550:hoss-test"
namespaces:
  - name: default
    description: Default namespace
    bucket: data
    object_store: default
  - name: s3-int-test
    description: Namespace for S3 integration tests
    bucket: com.gigantum.hoss-test
    object_store: s3
queues:
  - type: amqp
    settings:
      url: amqp://${RABBITMQ_USER}:${RABBITMQ_PASS}@rabbitmq:5672
    object_store: default
  - type: sqs
    settings:
      queue_name: hoss-test-api.fifo
      region: us-east-1
      profile: hoss-int-test
    object_store: s3
server:
  dev: true
  auth_service: http://auth:8080/v1
  sync_frequency_minutes: 5
```

## Setting up VSCode (Recommended for Developmet)
While you can run the test via pytest, you can easily run and deubg the tests in VSCode.

- Open the `test` dir in vscode
- Open the command palette and run the `Python: Select Interpreter` and select your virualenv from the Setup section
- Open the command palette and run the `Python: Configure Tests` and select pytest
- Copy the `.vscode/launch.json.template` file to `.vscode/launch.json`
- Get a PAT as the admin user from the server and set it in the launch.json file
- Run and deub tests right from vscode!


## Running Tests
If you prefer just running the tests from pytest:

- login as admin
- Get a PAT
- Export the env var `HOSS_PAT`
- Login as privileged and user as well to create all accounts
- Run tests via pytest
