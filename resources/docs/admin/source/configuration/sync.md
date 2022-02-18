# Sync Service
The sync service is responsible for synchronizing data between Hoss servers and indexing metadata for search.

You typically should run with a single sync service that is configured to sync data between servers. The sync service should be run in a location that has access to all servers (e.g. on-premise).

If you only have a single server, the sync service should still be enabled to index metadata.

## Sync Service Configuration
Configuring the sync service is done via `~/.hoss/sync/config.yaml`. It's values are described below:

* `refresh_intervals`
  * `core_service`: Rate at which the core service should be checked for new sync configurations
  * `auth_token`: Period between refreshing a worker's JWT. This should be less than (and ideally half) the JWT timeout set in the auth service
  * `sts_creds`: Period between refreshing a worker's STS credentials. This must be less than the max STS session duration.
* `core_services`: A list of core services to monitor. You must include any core service that you want the sync service to index or sync.
* `auth_endpoint`: The auth service endpoint. By default the internal Docker route is used. If using an auth service running in a different server, you must update this value.
* `elasticsearch_endpoint`: The endpoint where the Opensearch API is accessible. By default the internal Docker route is used. You should not have to modify this value.
* `sqs_profile`: The profile name in the `~/.hoss/sync/aws_credentials` file. If not needed (because you aren't using S3), just leave the default value.
* `worker_buffer_size`: The channel size for the worker channel. The larger the buffer the more messages can be queued for the worker(s) without the demuxer blocking. 
* `worker_instance_count`: The number of workers that should be started per core service. Typically this is fine to set at 1, but if you have lots of activity or data to sync, more workers could help. Setting this value too high may result in workers running out of bandwidth and sync operations timing out.


## Setting AWS Credentials
AWS credentials are provided to the sync service via the `~/.hoss/sync/aws_credentials`, which is bind mount into the service container. You should set the Hoss service account credentials in this file as shown below. You can use any profile name as long as you are sure to set it in all required config files.

```
[hoss-service-account]
aws_access_key_id = SDHFVMWJSD343ANSADNa
aws_secret_access_key = SDjfdsjSYwnd8*56$7s2hdsjdASF
region=us-east-1
```
