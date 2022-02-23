# Project Architecture

## Development Stack

### Typescript
The project uses [Typescript](https://www.typescriptlang.org/) in order to minimize run-time errors that are prevalent in native JavaScript.

### React/React Hooks
The components for the front end are built out using the [React](https://reactjs.org/) library. It utilizes [React Hooks](https://reactjs.org/docs/hooks-intro.html), a relatively new addition to the library that uses hooks to replace a lot of the previous functionality.

### xState
[xState](https://xstate.js.org/docs/) is being used throughout the application to create robust state machines. This proves to be extra useful building out pages with different views.


## Testing
The project uses [Storybook](https://storybook.js.org/docs/react/get-started/introduction) and [Jest](https://jestjs.io/docs/getting-started) for testing. Storybook allows you to render and view and interact with individual UI components. It will also run and display the results for written Jest tests.

### Writing Tests

In order for test files to be detected, the test files must be placed in a `__tests__` directory. To match our current project structure, you will want to create this directory in the same directory as the tested Component. For example, the test directory for the 'Dataset' component that exists in `/pages/dataset/Dataset.tsx` would be `/pages/dataset/__tests__`

Each testing directory should contain at least two files, a storybook file and a jest test file. Our current naming convention is `component-name.stories.tsx` and `component-name.test.tsx`.  You can find an example of this [here](https://github.com/gigantum/hybrid-object-store/tree/main/server/ui/ui/src/pages/dataset/__tests__).


### Running Tests
1. run `yarn storybook:build`
2. run `yarn storybook`

Note: If you only want to run unit tests and not individual components using story book, use `yarn test`


## Project Structure
This repository uses a strict project structure for the UI. It also utilizes aliasing to easily reference commonly referenced directories, i.e. `Images/icons/...` and `Components/button...`.

### Conventions
- Each component should be in its own directory with a corresponding testing folder and `.scss` file for styling purposes if needed
- Commonly reused components should be defined in `src/components/componentName/variationName` with a `index.tsx` file in the above folder with exports to all the different variations.
- Custom hooks should be defined in `src/hooks/hookType/hookName`
- All assets should be placed in `src/images/imageType/imageName`
- Individual routes should be placed in their own directory within `src/pages/`


## Server API

The Hoss server uses a REST API. There are two different API services that the UI has to interact with; the core service (`core/v1`) and the auth service (`auth/v1`).

When in developer mode, server will automatically serve up documentation for the services at `http://localhost/core/v1/swagger/index.html` and `http://localhost/auth/v1/swagger/index.html`. 

To easily make calls to the services it is highly recommended to import the appropriate function from `src/environment/createEnvironment.ts` as the utility has already been set up to.
