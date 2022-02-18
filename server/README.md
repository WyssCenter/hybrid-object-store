# Hoss Server

## Running via Docker Compose
You can run the entire system locally via docker compose. By default it currently will run everything on localhost in dev mode.

> To expose the services externally you need to edit `.env` and change the `EXTERNAL_HOSTNAME=` value to be the external DNS name that clients will be using to access the HOSS services. Afterward run `make config` to regenerate files with the new hostname before running `make up`

> To run in production mode (e.g. disabling cors), set `dev_server` to `false` in the `auth` and `core` service `config.yaml` files.

Currently, the core service config file will bootstrap only the local minio object store at start. To enable the S3 backed store,
you must configure the `~/.hoss/core/aws_credentials` file per the instructions at the end of `docs/s3-install.md` and add another
object store and namespace to the core config file (`~/.hoss/core/config.yaml`). Then you must restart the services to create
and load the new object store and namespace. 


The following are the safest steps to get running.

1. run `make reset`
2. run `make build`
3. run `make env`
4. run `make up`

In this configuration you should be able to reach minio directly on localhost and load the minio browser when
visitin that in your browser.

## Obtaining an auth token
Once the services are running you need to get an auth token so that you can query the HOSS services and get credentials for minio.

To obtain a JWT token:
1. Browse to http://localhost/auth/v1/login
2. Log into Dex using one of the accounts below to gain access with the desired role
3. Copy the ID Token that is returned by the Auth service and export it as the `HOSS_JWT` environmental variable


| username               | password | role       |
|------------------------|----------|------------|
| admin@example.com      | foo      | admin      |
| privileged@example.com | bar      | privileged |
| user@example.com       | password | base user  |


> The Auth service JWTs have an expiration of 24hrs. After the JWT expires you will need to login again and re-export the new ID token JWT

Once you've obtained an auth token you can use the client library (in the `client/` directory in the root of the repository) to access the Auth and Dataset HOSS services and connect to minio to work with dataset data.


## Managing users
By default, an LDAP server is run with the phpLDAPadmin console. Dex then uses this LDAP provider to authenticate users.

After running `make env` you can modify the LDAP related env vars before running `make config`.

Additionally, if you want a more complex auth configuration, you should make those changes before running `make up`

If running things using the defaults, 3 test users will be created that you should delete in a "real" deployment.

phpLDAPadmin is running on localhost:6443 and is not routed through the proxy. This adds a layer of complexity for managing, but increases security.
Future work could route through the proxy with additional auth if desired.

To add/remove/modify users:

- ssh into the server with something like `ssh -L 6443:localhost:6443`
- Open `http://localhost:6443` in your browser
- Log in using the `LDAP_ADMIN_PASSWORD` env var found in the `.env` file.


## Multiple Instance Deployments
When there are multiple deployments of the HOSS (e.g. on-prem and in the cloud) and there is a requirement to sync data between the deployments, a couple of tweaks are needed for the Sync Service to work correctly. There are two deployment scenarios: a single central Auth Service and multiple mirrored Auth Services.

### Single Auth Service
In this scenario a single Auth Service is run in a location where all users and services can access it for all of the auth needs of the whole system.

In the deployment where the Auth Service is being run:
1. Update `~/.hoss/auth/config.yaml` `additional_allowed_servers` list to include the external hostnames of the additional servers (e.g. https://hoss2.mydomain.com).

For each deployment where the Auth Service is not being run:
1. Update the `~/.hoss/.env` file so that the `auth`, `ldap`, and `ldap-admin` services are not listed in the `SERVICES` variable.
2. Update the Core (`~/.hoss/core/config.yaml`) and Sync (`~/.hoss/sync/config.yaml`) Service configuration files to update the Auth Service references (e.g. https://hoss1.mydomain.com/auth/v1)
3. Copy the service token value from the `.env` file on the server running the auth service and set it in the other server(s)


### Multiple Auth Services
In this scenario each deployment will run its own Auth Service instance. There are multiple reasons for this, including having different users for each deployment. The way this scenario works is by selecting on Auth Service as the primary one and updating the other Auth Services to mirror its configuration.

Note: These instructions don't cover changing the hostname of the deployment from `localhost` nor updating the service config files to reflect this change in deployment configuration.

On the primary Auth Service:
1. Extract the `/secrets/private.pem` key from the running service container (`docker cp (id):/secrets/private.pem ./auth_service_private.pem`)

For each of the other deployments:
1. Edit the `auth/openid-configuration.json.tmpl` file and replace the `issuer` value with the URL of the primary Auth Service (all other values are untouched)
2. Run `make config`
3. Run `make up`
4. Replace the `/secrets/private.pem` key in the running Auth Service container (`docker cp ./auth_service_private.pem (id):/secrets/private.pem`)
5. Restart the Auth Service so that it uses the new private key

Once this is all finished each non-primary Auth Service will be issuing JWTs that match those of the primary Auth Service, allowing a JWT generated by one Auth Service to be use on the services in a different deployment.
