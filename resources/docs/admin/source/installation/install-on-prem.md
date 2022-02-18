# Single On-premise Server

This section describes setting up a single server in a typical configuration for use "on-premise", 
but ultimately you could run this configuration in the cloud. The main thing to understand is the 
choices made here are to:

* Run all services, including minIO
* Use internal LDAP provider
* Running a server with the FQDN `hoss.mycompany.com`
* An NFS mount at `/mnt/storage` provides a large storage volume



## Set Up Repository
The Hoss system is installed and managed directly from the source repository. The version of this repository will determine the
version of the server that is run. You can checkout tags to deploy specific releases (e.g. `git checkout 0.2.3`).

The system also uses a working directory `~/.hoss` contains all configuration information and data for the running server. This
location is created and pre-populated via the Makefile.

To get started, clone the repository:

```shell
git clone https://github.com/gigantum/hybrid-object-store.git
```

This will install the latest version of the server. To install a known release, review the [Releases](https://github.com/gigantum/hybrid-object-store/releases) 
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
EXTERNAL_HOSTNAME=http://hoss.mycompany.com
DOMAIN=hoss.mycompany.com
NAS_ROOT=/mnt/storage
LDAP_ORGANISATION=My Company Inc.
LDAP_DOMAIN=hoss.mycompany.com
```

Finally, run `make config` to finish configuring all services with base settings.

## Configure the Core Service
There is additional configuration to the core service possible depending on your deployment scenario. You can review 
all possible items in the [Core Service](../configuration/core.md) configuration section. 

For this scenario, edit `~/.hoss/core/config.json`. 

Set the `endpoint` value for the `default` object store to be equal to what was set for the `EXTERNAL_HOSTNAME` variable.

Set `server.dev` to `false`

## Configure the Sync Service
There is additional configuration to the sync service possible depending on your deployment scenario. You can review all 
possible items in the [Sync Service](../configuration/sync.md) configuration section. 

For this scenario, edit `~/.hoss/sync/config.json`. Set `core_services` to include this server's
external hostname. For example, assuming the server is running at `hoss.mycompany.com`:

```yaml
core_services:
  - http://hoss.mycompany.com/core/v1
```

## Configure the Auth Service
There is additional configuration to the auth service possible depending on your deployment scenario. You can review all 
possible items in the [Auth Service](../configuration/auth.md) configuration section. 

For this scenario, edit `~/.hoss/auth/config.json`. 

Set `server.dev` to `false`

Optionally modify the `password_policy` as desired. These critera are enforced on users when they change their password via the Hoss UI.


## Configure the UI Service
There is additional configuration to the UI service possible depending on your deployment scenario. You can review all possible items in the [UI Service](../configuration/ui.md) configuration section. 

For this deployment, set the "name" of the server that is presented to users in the menu bar. Edit `~/.hoss/ui/config.json` and set the `server_name` value to what ever you wish, e.g. "My Company Hoss Server - On-premise"

## Configure the Ingres Proxy
Since this example configuration is not using TLS termination for the on-premise install, an additional
modification is recommended to disable development features in the proxy. Note, it is recommended to use
TLS when possible.

In `~/.hoss/traefik.yaml`, set `log.level=WARNING` and `api.insecure=false`

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