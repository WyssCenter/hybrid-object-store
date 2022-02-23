# Makefile

The Hoss deployment and development is driven by a Makefile.

## setup
This target creates the additional docker network `web`. If it's the first time you run the Hoss on a system you'll need to run this.

## env
The env build target initializes the Hoss working directory and `.env` file. It is always the first thing that is run. Most other targets require an `.env` file to exist before they will run.

If new settings that are sensitive or shared across services are needed, they should likely be added here. When modifying the .env file, you must make sure the release notes indicate changes admins must make when upgrading.

## config
The config build target takes values in the .env file and write out additional config files:

1) Configures Traefik for https via Let's Encrypt or http
2) Configures Dex
3) Initializes additional directories in the Hoss working directory
4) Initializes the UI customization configuration
5) Generates and updates the `UI_REDIRECT_REGEX` value in the `.env` file
6) Checks if your UID is 1000, and if not warns about permission issues for Opensearch

Effort should be made to make it safe to run `make config` multiple times. If you run `make reset` and then `make env` while developing, it's likely that you will have to delete `~/.hoss/auth/config.yaml` and `~/.hoss/auth/config-dex.yaml` before running `make config` to ensure the OAuth2 client secret is properly set everywhere.

## build
The build build target builds all the containers defined by Docker Compose based on the settings in the `.env` file. Depending on your configuration different containers may be built or pulled.

## up
The up build target starts the server via docker compose. If you include a `DETACH=true` variable the system will launch in the background instead of foreground. For development purposes it's typically preferable to run the the foreground to monitor log output.

## down
The down build target stops a running server and removes the stopped containers.

## restart
The restart build target is used to restart a single service, e.g.:

```
make restart SERVICE_NAME=core
```

This is useful if you need to recreate a service without taking the whole stack down. 

## reset
The reset build target removes Hoss configuration and resources used by the running server. It is a way to get a "fresh" server that will boot up as if it has never been used. The target:

1) Prunes stopped containers
2) Removes the etd volume
3) Removes the auth volume (which contains the private PEM used to sign JWTs)
4) Removes the LDAP server volumes
5) Removes rabbitmq persisted data
6) Removes the data directory - Note, this will erase and object data that was stored by minIO if running with the default configuration!
7) removes the `.env` file

While developing, if using this command it's likely that you will have to delete `~/.hoss/auth/config.yaml` and `~/.hoss/auth/config-dex.yaml` before running `make config` to ensure the OAuth2 client secret is properly set everywhere.

## watch-logs
The watch-logs build target displays the last 10 lines from all of the service log files and then streams logs from the Hoss services to the console. This can be very useful to view logs and debug when running with `make up DETACH=true`.

You can use the `SERVICE_NAME` variable to specify which service to watch, e.g:

```
make watch-logs SERVICE_NAME=core
```

## get-logs
The get-logs build target is similar to watch-logs except it gets all available logs and returns them without color ANSI escape codes included. This can be useful when you want to collect logs for analysis.

You can use the `SERVICE_NAME` variable to specify which service to inspect.

## api-docs
The api-docs build target builds the swagger API docs. You only need to run this if you are trying to build and run core or auth services locally (not in a container). This scenario usually comes up when you are trying to run unit tests.


## up-testing
The up-testing build target starts a subset of the Hoss required for running unit tests
