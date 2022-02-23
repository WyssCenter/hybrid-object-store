# Docker Compose

The system uses Docker Compose extensively to deploy the application. Multiple compose files are used, managed independently by each service, and selected automatically by the Makefile.

The `.env` file is passed into docker compose and all vars are available to use in the compose files.

## Organization
Each service contains a primary compose file in its directory inside the `server` directory. These are all loaded and combined together at runtime in the makefile.

If TLS is enabled, there is a `docker-compose-tls.yaml` file for each externally routed service to add additional labels for Traefik.

This pattern of adding additional fields to an existing docker compose service is used to apply additional optional configurations as well. These compose files are in the `server/options` directory.

## Traefik Labels

Traefik is used with the dynamic docker backend. To add services to Traefik, simply add labels to the containers in docker compose. 

It's important to note that we run minIO on the root of the server because it can't run on a prefix. We then redirect a user that goes to the root of minIO to the UI service.