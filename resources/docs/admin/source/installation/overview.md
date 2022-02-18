# Overview

## Hoss Summary
The Hoss is currently deployable as a collection of containers running via Docker Compose. The architecture is
designed to allow independent deployment and scaling of components as needed, but future work is still
required to support more complex deployments.

Depending on the configuration, the following services are deployed, each in its own container:

- `db`: A PostgreSQL database
- `core`: A service providing the primary REST API
- `auth`: If enabled, a service providing a REST API for managing groups, PATs, token exchanges, and user accounts
- `dex`: If `auth` is enabled, [Dex](https://dexidp.io/) (a federated OpenID Connect provider) is run to integrate with external auth providers
- `ldap`: If enabled, a local LDAP server is run to manage user accounts internally
- `ldap-admin`: If enabled, a small service provids an admin UI for the internal LDAP server
- `minio`: If enabled, minIO is run in gateway mode to provide access to local storage via the S3 API
- `etcd-0`: If `minio` is enabled, etcd is run to support required minIO features
- `rabbitmq`: If `minio` is enabled, RabbitMQ is run to provide message queues to handle minIO bucket events
- `sync`: A service to manage synchronizing data and indexing metadata for search
- `opensearch`: [Opensearch](https://opensearch.org/) is run to provide metadata search
- `ui`: A service to serve the Hoss web UI

## Choosing The Right Configuration

The most complicated part about installing a Hoss server is likely understanding and determining a desired configuration.
Because the system is so flexible, there are many choices related to not only how individual servers are
deployed, but also how multiple servers can be linked together to enable hybrid cloud workflows.

### Server Options
When deciding on an individual server's configuration there are various options to be considered
and parameters to be set. More details are available throughout the rest of the documentation, but
at a high level you must consider:

- [The external hostname](install-on-prem.md#configure-variables) of the server
  - This cannot easily be changed after a server has been deployed.
- [TLS configuration](../configuration/tls.md)
  - The Hoss can use Let's Encrypt internally, run behind an additional proxy/load balancer that is doing TLS termination for you, or run unencrypted
- [Auth Configuration](../configuration/auth.md)
  - A Hoss server can run its own `auth` service or use an existing one (i.e. in a multi-server configuration)
  - An [internal LDAP provider](../maintenance/internal-ldap.md) can be used to work "out of the box"
  - External LDAP or other authentication providers can be [integrated](../configuration/auth.md)
  - If using the internal LDAP provider, you may want to enable [Google's reCaptcha service](../configuration/captcha.md) on the login page
- Object Store Configuration
  - Currently, you can choose between AWS S3, and externally hosted minIO server, or an internally hosted minIO server
- [Backup location](../maintenance/backup-and-restore.md)
- [Custom UI colors and logos](../configuration/ui.md)

### Multi-Server Configurations
Often, multiple Hoss servers are run on different infrastructures (e.g. one server on-premise and one server in AWS) and
linked together via syncing, auth, or both. These architectural decisions can enable useful hybrid cloud workflows, for example:

- Easy sharing of data generated and managed on-premise with external collaborators
- Off-site data collection and transfer back on-premise
- Portable analytics to leverage both on-premise and cloud compute resources
- Data "delivery" to external users

The first consideration in a multi-server configuration is how auth will be configured. In addition to deciding on what
authentication provider will be used, you must also choose between:

1) (Recommended) One server runs an auth service. Additional servers use the "centralized" auth service.
   * Less steps to configure, and easier to use and manage
   * Auth service must be accessible by all other servers (e.g. runs in the cloud, not on-premise behind a firewall)
   * Users can use a single PAT with any linked server and have only one set of credentials to remember
2) Each server runs its own auth service
   * More complex to configure and manage
   * Depending on the Auth provider configured for each server, different credentials may be needed. 
   * You must be careful to make sure usernames match between servers or there could be issues when synchronizing data, groups, and permissions.


The second consideration is how syncing will be configured. Typically you should run a single sync service that can reach all servers.
For example, if you have one server on-premise and one server in the cloud, you'll likely want to run the sync service on-premise. This sync
service can then be responsible for moving data between object stores around as needed.


## Installation Process
Installing a server requires several manual steps and configuration.
The system is quite flexible and can support various use cases and deployment architectures. Given your
decisions on how to configure both individual servers and if you will be linking multiple servers, the
process at a high level is:

1) Prepare required infrastructure
   1) Create any required cloud resources (i.e. S3 buckets, SQS queues, IAM roles & policies, EC2 instance)
   2) Create any required on-premise resources (e.g. NFS shares, a VM or server)
2) Prepare the server 
   1) Install Docker and Docker Compose
   2) Install additional `make` and `git` dependencies
   3) Configure host user accounts and storage mounts as needed
3) Install, configure, and start the Hoss server software

The following installation documents outline the installation process for common configurations. Details on all 
the available configuration parameters and scenarios is captured in the "Configuration" section.