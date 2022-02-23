This directory contains tests and scripts for running integration level tests using the Python client library and pytest.


## Setup
To get started, first get a Python virtualenv setup to run the library and test code.

- create a virtualenv
- Install the client library. If you are just running the tests you can typically install the latest version from PYPI. If you are developing a new feature that requires and up-to-date version that has yet to be release, you may have to install from the git repo directly. Typically:

```
pip install -U hoss-client
```

- Install testing deps from requirements.txt

```
cd test
pip install -r requirements.txt
```

## Running Base Tests
These tests assume you are running the server locally. Follow the server README to get the server up and running
locally. For all tests except for S3, using the default configuration should work fine. Typically it is best to reset
your server before running tests to ensure the search index is empty. 

To get started:

- login as all three test accounts (admin, privileged, and user)
- log in as the admin user
- Generate a PAT
- In your terminal with the virtual env activated, export and set the env var `HOSS_PAT` equal to your PAT 
- Run tests via pytest


## Setting up a Hoss server to run S3 tests
A Hoss server must be configured for S3 to run the S3 tests and running. This means you must create an S3 bucket and associated
AWS resources (SQS, IAM, etc) as outlined in the [Single AWS Server](https://hybrid-object-store.readthedocs.io/en/latest/installation/install-aws.html#single-aws-server) admin documentation. 

You can still run the tests locally, but you must configure your local system to have an S3 object store and namespace.

Make sure to set your core and sync service credential files:

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
    role_arn: arn:aws:iam::1234567890:role/hoss-user-assume-role
    notification_arn: "arn:aws:sqs:us-east-1:1234567890:hoss-test"
namespaces:
  - name: default
    description: Default namespace
    bucket: data
    object_store: default
  - name: s3-int-test
    description: Namespace for S3 integration tests
    bucket: my-example-int-test-bucket
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
  elasticsearch_endpoint: http://opensearch:9200
  sync_frequency_minutes: 5
  dataset_delete_delay_minutes: 0
  dataset_delete_period_seconds: 2
```

Finally, when running integration test via pytest, add the `--s3` flag to enable tests that require S3. Note, these
tests take some time to run due to frequently reconfiguring the sync service.

## Setting up VSCode (Recommended for Developmet)
While you can run the test via pytest, you can easily run and deubg the tests in VSCode.

- Open the `test` dir in vscode
- Open the command palette and run the `Python: Select Interpreter` and select your virualenv from the Setup section
- Open the command palette and run the `Python: Configure Tests` and select pytest
- Copy the `.vscode/launch.json.template` file to `.vscode/launch.json`
- Get a PAT as the admin user from the server and set it in the launch.json file
- Run and deub tests right from vscode!
