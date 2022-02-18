## Docker Compose Options

This directory contains docker compose files that are used to enable various optional configurations.

The makefile will include these automatically based on the server's configuration.


## Available Options

- `auth-minio-dependency.yaml`: File included if both the `auth` and `minio` services are enabled. This adds a dependency on the `auth` service to the `minio` service
- `no-minio-redirect.yaml`: File included if `minio` service is NOT enabled. Provides a redirect from `/` to `/ui`.
- `health-check.yaml`: File provides a route on a different host to make sure an internal health check route would succeed.