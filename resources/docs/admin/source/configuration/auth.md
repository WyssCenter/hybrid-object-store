# Auth Service
The auth service is responsible for managing groups, personal access tokens, and token exchanges. By default a Hoss server will deploy an auth service, but only one service is required for a collection of linked servers.

## Auth Service Configuration
Configuring the auth service is done via `~/.hoss/auth/config.yaml`. It's values are described below:

- `dev_server`: if `true`, CORS will be enabled for local frontend development and the API will run in development mode
- `client_id`: This is the OAuth2 client id
- `client_secret`: This is the OAuth2 secret. In a production environment, you should modify this to a strong, random value. If changed, you must also update your Dex config. 
- `open_id_config_file`: this is the path to the OIDC file that will be served as `.well-known` info
- `issuer`: Is the token issuer. If you are running with the internal LDAP provider, you should leave the default value of `http://dex:5556/dex`. If you are using an external auth provider that will need to make callbacks to Dex, this should instead point to the external route to Dex in the format `<EXTERNAL_HOSTNAME>/dex`, where `<EXTERNAL_HOSTNAME>` is the externally accessible hostname for your Hoss server, including the scheme (e.g. https://hoss.myserver.com/dex).
- `admin_group`: Group name to expect from the Auth provider for administrators. If integrating with an external Auth provider you may need to adjust this. If using the internal LDAP provider, you do not need to change this.
- `privileged_group`: Group name to expect from the Auth provider for privileged users. If integrating with an external Auth provider you may need to adjust this. If using the internal LDAP provider, you do not need to change this.
- `token_expiration_hours`
  - `access`: Number of hours until an access token expires. Typically should be the same as the `id` token. 
  - `id`: Number of hours until the id token expires. Typically should be the same as the `access` token. 
  - `refresh`: Number of hours until the refresh token expires. Note, refresh tokens are not fully supported, so this can be ignored.
- `additional_allowed_servers`: A list of servers (e.g. https://my.server.com) that are allowed to use this auth service. If a server is not included here, an attempt to redirect to that server will be blocked.
- `password_policy`: This section defines the requirements when a user changes their password via the Hoss UI. This is only in effect when using the internal LDAP provider.
  - `min_length`: Minimum number of characters
  - `require_uppercase`: The password must contain uppercase characters
  - `require_special`: The password must contain special characters
- `username_claim`: This defines what claim the Hoss should use to determine a username. The default value of `nil` will replicate the behavior prior to this config value being added. The Hoss will first look for a `nickname` claim, then a `name` claim, and finally an `email` claim to generate the username. You may also now explicitly set this to `nickname`, `name`, or `email`. The `email` option will take the first part of the email address before the `@`, remove any `+` suffix, and replace all characters other than `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_.` with `-`.


## Dex Configuration
Dex is used to to federate OpenID Connect providers. The auth service uses Dex to handle login, and then generates tokens as needed from there. Dex is primarily configured via the `~/.hoss/auth/config-dex.yaml` file. By default, you don't need to touch this file, but there are a few cases where you would.

If you set a better `client_secret`, you should edit the `secret` value under the `staticClients` section.

If you are using an external LDAP provider, you'll need to provide details to connect with that LDAP server as discussed in the External LDAP section below.

The Dex UI is customized via the `~/.hoss/auth/web` directory. If needed you can make changes there, but default/automatic configuration should work in most cases.

## Auth Providers

### Internal LDAP
By default an internal LDAP server is deployed, along with a managment UI. There is little additional configuration needed to make this work. You should remove the test accounts before using in a real environment. You can learn more about managing the LDAP server in the [Internal LDAP Server](../maintenance/internal-ldap.md) section.

### External LDAP
To use an external LDAP provider, you would disable the `ldap` and `ldap-admin` services. Then
edit the `~/.hoss/auth/config-dex.yaml` file to connect to your LDAP provider instead of the internal
server, which is automatically populated. This configuration has not been extensively tested, but should work.

If you must use a custom certificate authority (CA), you can place the certificate in the `~/.hoss/auth/certificates` directory. At runtime this
folder will be mounted into the Dex container at `/opt/certificates`. You can then reference the cert from your `~/.hoss/auth/config-dex.yaml` file
as needed.

### Microsoft Azure Active Directory

A common configuration is using Microsoft Azure Active Directory. This is a good option if you are already using Azure AD at your organization because you'll get additional the benefits like MFA and single sign on with minimal additional management overhead. Follow the instructions [here for more details](auth-azure-ad.md).

### Additional External Auth Providers

Dex supports various external auth providers, and includes documentation on how to configure them: [https://dexidp.io/docs/connectors/](https://dexidp.io/docs/connectors/)

## Running Multiple Servers
When you run multiple Hoss servers that are configured to work together (i.e. via syncing), you must decide if the servers will share a single auth service, or if each server will run its own. 

Running individual auth services is a more complex to configure and maintain. This configuration is likely only useful in scenarios where network partitioning does not allow a single auth service deployment or you wish to have different users be able to access different servers. This can effect the ability to apply permissions when syncing if usernames change between servers.

Running a shared auth service is more simple for the end user. It also ensures that permissions will always work when syncing because each user has a single account. It also can have less maintenance overhead.

### (Recommended) Shared Auth Service
In a shared auth service configuration, one auth service is run in one Hoss server, with other Hoss servers deferring to that auth service for authentication. 

When configured properly, each Hoss server will broadcast this configuration via the `/core/v1/discover.json` endpoint, allowing clients to automatically contact the correct auth service for PAT exchange.

Review the [Dual Server Hybrid Cloud Installation](../installation/install-aws.md) section for an example of how you would configure this type of deployment.

You must:

* Properly set `AUTH_SERVICE_ENDPOINT` in the `~/.hoss/.env file
* Copy the `SERVICE_AUTH_SECRET` value from the server running the auth service to all other servers
* Properly set `server.auth_service` in `~/.hoss/core/config.json`
* Properly set `auth_endpoint` in `~/.hoss/sync/config.json`

### Individual Auth Services
In this scenario each deployment will run its own auth service instance. The way this scenario works is by selecting an auth service as the primary one and updating the other auth services to mirror its configuration.

After completing the initial deployment of both servers, on the "primary" auth service:
1. Extract the `/secrets/private.pem` key from the running service container (e.g. `docker cp (id):/secrets/private.pem ./auth_service_private.pem`)

For each of the other deployments:
1. After running `make config`, edit the `~/.hoss/auth/openid-config.json` file and replace the `issuer` value with the URL of the primary Auth Service (all other values are untouched). Note that if you ever run `make config` again (e.g. during an update) you will have to redo this step!
2. Run `make up`
4. Replace the `/secrets/private.pem` key in the running Auth Service container (e.g. `docker cp ./auth_service_private.pem (id):/secrets/private.pem`)
5. Restart the Auth Service so that it uses the new private key (`make restart SERVICE_NAME=auth`)

Once this is all finished each non-primary auth service will be issuing JWTs that match those of the primary auth service, allowing a JWT generated by one auth service to be use on the services in a different deployment.
