# API Documentation

API docs are maintained via [Swag](https://github.com/swaggo/swag) which converts Go annotations to Swagger Documentation 2.0.
The swagger docs are then served when running in development mode.

## Viewing the docs
Docs are available for the `core` and `auth` REST APIs. To view the docs, run these services in development mode. 

For the core service, in `~/.hoss/core/config.yaml` set `server.dev` to `true`. After starting the server open `/core/v1/swagger/index.html`.

For the auth service, in `~/.hoss/auth/config.yaml` set `dev_server` to `true`. After starting the server open `/auth/v1/swagger/index.html`.

## Developers
Swaggo v1.7.8 is currently in use. If you upgrade the package in the core and auth go.mod files, you must also update the install step
in both core and auth Dockerfiles and `.github/workflows/main.yaml`.

While creating or modifying API endpoints, you should update the associated Go annotations. This will ensure that documentation is always correct
once released.

You can simply build and run the server in development mode to view your changes. If you wish to build the docs locally you'll have to follow the
instructions in the [Swag](https://github.com/swaggo/swag) readme to install the `swag` tool.
