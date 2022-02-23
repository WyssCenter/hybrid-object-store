# Frontend Development

## Setting Up The Hoss Server

In order to begin frontend development you will want to have a Hoss server set up and running locally.

### Dependencies
To get started, first install the required dependencies to run and build the code base.

1.  Install Docker, likely via [Docker Desktop](https://www.docker.com/products/docker-desktop)
2.  Install [Docker Compose](https://docs.docker.com/compose/install/). This is included in Docker Desktop, but depending on your host OS you may need to follow additional steps.
3.  Install `make` and `git`


### Set Up a development Hoss server
To develop the UI you need a functioning Hoss server running locally so API requests will work. The default configuration of the server is set up for development.

1. Clone the code repository
2. cd into `hybrid-object-store/server`
3. run `make setup`
4. run `make env`
5. run `make config`
6. run `make build`
7. run `make up`
	> It will take a few moments for the API to be ready

Note: Run `make down` to ensure the server is gracefully turned off.

Note: To refresh your install back to a clean starting point, run:
1. `make reset`
2. `rm -rf ~/.hoss` (you can skip this step if you wish to keep your configuration)
3. `make env`
4. `make config`
5. `make build`
6. `make up`

## Running The Frontend Locally

### Requirements

You must have your local development environment configured. When running in frontend development mode, you'll be serving the UI locally instead of from the UI container.

1. Install [Node](https://nodejs.org/en/) 14
2. Install [Yarn](https://classic.yarnpkg.com/en/)


### Starting the Frontend Dev Server

The UI project directory is `hybrid-object-store/server/ui/ui`. In this directory,

1. run `yarn` (this is the same as `yarn install` command)
2. run `yarn start`
   * Starts the development server and will automatically open http://localhost:3000 in the browser. Hot loading has been set up for this project and the page will reload if edits are made to the frontend.

Alternatively to create a production build for deployment use `yarn build`
