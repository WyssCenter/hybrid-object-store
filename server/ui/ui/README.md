# Getting Started with Hoss for frontend

##install
1. `brew install make`
2. `pip install docker-compose`


## Build API
1. `make build`
2. `make env`
3. `make up`

May need to run `docker network create web`


## Storybook
Run to build static assets.
`yarn storybook:build`

To serve run.
`yarn storybook`

## API routes
Auth -> `/auth/v1/login`

// TODO write api docs here
v1.GET("ping", api.Ping)

v1.GET("namespace/", api.ListNamespaces)
v1.POST("namespace/", api.CreateNamespace)
v1.GET("namespace/:namespace", api.GetNamespace)
v1.DELETE("namespace/:namespace", api.DeleteNamespace)

v1.POST("namespace/:namespace/dataset/", api.CreateDataset)
v1.DELETE("namespace/:namespace/dataset/:name", api.DeleteDataset)
v1.GET("namespace/:namespace/dataset/:name", api.GetDataset)
v1.GET("namespace/:namespace/dataset/", api.ListDataset)
v1.PUT("namespace/:namespace/dataset/:name/user/:username/access/:accesslevel", api.UpdateUserDatasetPerms)
v1.DELETE("namespace/:namespace/dataset/:name/user/:username", api.RemoveUserDatasetPerms)
v1.PUT("namespace/:namespace/dataset/:name/group/:groupname/access/:accesslevel", api.UpdateGroupDatasetPerms)
v1.DELETE("namespace/:namespace/dataset/:name/group/:groupname", api.RemoveGroupDatasetPerms)

v1.GET("namespace/:namespace/sts", api.GetSTSCredentials)

v1.PUT("user/sync", api.SyncUserGroups)




## Available Scripts

In the project directory, you can run:

### `yarn start`

Runs the app in the development mode.\
Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

The page will reload if you make edits.\
You will also see any lint errors in the console.

### `yarn test`

Launches the test runner in the interactive watch mode.\
See the section about [running tests](https://facebook.github.io/create-react-app/docs/running-tests) for more information.

### `yarn build`

Builds the app for production to the `build` folder.\
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.\
Your app is ready to be deployed!

See the section about [deployment](https://facebook.github.io/create-react-app/docs/deployment) for more information.

## Learn More

You can learn more in the [Create React App documentation](https://facebook.github.io/create-react-app/docs/getting-started).

To learn React, check out the [React documentation](https://reactjs.org/).

### Code Splitting

This section has moved here: [https://facebook.github.io/create-react-app/docs/code-splitting](https://facebook.github.io/create-react-app/docs/code-splitting)

### Analyzing the Bundle Size

This section has moved here: [https://facebook.github.io/create-react-app/docs/analyzing-the-bundle-size](https://facebook.github.io/create-react-app/docs/analyzing-the-bundle-size)

### Making a Progressive Web App

This section has moved here: [https://facebook.github.io/create-react-app/docs/making-a-progressive-web-app](https://facebook.github.io/create-react-app/docs/making-a-progressive-web-app)

### Advanced Configuration

This section has moved here: [https://facebook.github.io/create-react-app/docs/advanced-configuration](https://facebook.github.io/create-react-app/docs/advanced-configuration)

### Deployment

This section has moved here: [https://facebook.github.io/create-react-app/docs/deployment](https://facebook.github.io/create-react-app/docs/deployment)

### `yarn build` fails to minify

This section has moved here: [https://facebook.github.io/create-react-app/docs/troubleshooting#npm-run-build-fails-to-minify](https://facebook.github.io/create-react-app/docs/troubleshooting#npm-run-build-fails-to-minify)




### Architecture

Architecture documentation can be found [here](docs/architecture.md)

### Testing
Testing documentation can be found [here](docs/testing.md)
