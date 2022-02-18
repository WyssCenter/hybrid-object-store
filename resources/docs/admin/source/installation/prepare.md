# Preparing for Installation
You must prepare for all servers you plan to deploy. This includes preparing cloud and on-premise
infrastructure, along with configuring the instance that will run the Hoss software.

In the Hoss, a `namespace` is an abstraction around an individual bucket. Each server requires
at least one `namespace` and associated bucket to be configured. 

## Prepare Required Infrastructure
Depending on your configuration, there may be additional infrastructure that needs to be deployed. When thinking
about what is required, you need to determine what object store you will use and where the server will run.

If you plan to use minIO for object storage locally, you might want to provision external storage via something like NFS to deal
with large data that may not fit directly on the instance. The env var `NAS_ROOT` can be modified during installation 
to use whatever mount location is chosen. Also, having a location to store backups off the server is 
recommended and another possible use for an external file mount. The [on-premise single server](install-on-prem.md) installation 
instructions provide more details on this configuration.

If using S3 for object storage, you must provision and configure a bucket, SQS queues, and IAM roles, policies, and user. There is
a [Terraform](https://www.terraform.io/) module available to help manage this - [https://github.com/gigantum/terraform-hoss-aws](https://github.com/gigantum/terraform-hoss-aws).
Alternatively, you can manually configure resources if needed. The [AWS single server](install-aws.md) installation instructions provide more details on
this configuration.

If you enable bucket versioning in S3 (recommended) you should also configure bucket lifecycle rules. Read more about this process in the
[Bucket Versioning section.](../configuration/versioning.md)

You will also need to provision a VM or server to run the Hoss software itself. The Hoss has been primarily tested on Ubuntu 18.04 and 20.04, 
but likely works well on any linux distribution that properly supports running Docker and Docker Compose. If running all of the services, a minimum recommendation is:
 * 2 cores
 * 8GB of RAM
 * At least 48GB of disk space 

Depending on your configuration and load, you may benefit from more RAM or need less. Depending on your configuration and load, you may benefit from more 
cores, especially if there will be a lot of data syncing and getting indexed. If you are running minIO locally and do not have external storage attached,
it is likely you'll need more disk space. The current disk space recommendation is for ensuring there will be plenty of space for Docker images, updates,
database and search index storage, and local backups if needed. Again, depending on your use case you may need more storage capacity.

## Prepare Server

To prepare for installation of the software, you must install Docker, Docker Compose, make, and git. While the details may vary depending on your
host's distribution, for Ubuntu you would:

1) Install Docker following the instructions [here](https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository)
2) Install Docker Compose following the instructions [here](https://docs.docker.com/compose/install/#install-compose-on-linux-systems)
3) Install make and git via `sudo apt-get install make git`

Note, the Hoss will attempt to run and map file permissions to the user who runs the `make` commands. It's possible you may want to
create a service account to complete the installation. 

You can optionally configure Docker's logging settings to prevent local logs from filling storage over time. Review the
 [Monitoring Logs](../maintenance/monitor-logs.md#configure-docker-logging-driver) section for more detail.
