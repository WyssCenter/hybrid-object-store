## Architecture

This application was built using React hooks and react context libraries and xstate.


### Auth and it's state machine

### Pages and state machines

Components are organized into pages to correspond with the routes in the Routes.tsx file.

React router is being used to manage the applications routing.

Each route should have a machine to manage to the pages loading state. In the pages directory there is a generic loading machine with the following states.

##### idle
This state is for the initial load of the page, this will display return a null and kick off the pages fetch/loading state.

##### loading
Loading sate is the state that will persist until data is returned and from the initial fetch.

##### refetching
When data has been manipulated the refetching state will be fired to get the updated data from the backend.

##### success
Success state will have a successful query returned and will display the pages main content with the data that has been fetched.

##### error
Error state is to handle api query errors and will diplay an error message to the user to communicate system issues.


### Generic components

Generic components are an iternal library of components that are to be used within the Hoss application.


Generic components should be structured as follows.

- Componets should be placed in a logical file system and be exported via a group index file.

Examples of importing buttons.


```js
import Button, { IconButton } from `Components/button/index`;


// from the following tree structure

- Components
  - button
    - Button
      -Button.tsx
    - icon
      - IconButton.tsx
    - index


// Components/button/index
import Button from './button/Button'
import IconButton from './icon/IconButton'

export {
  Button,
  IconButton,
}

export default Button;

```


Generic components should be agnostic to the pages using them.


*Good*
- all states can be manipulated by the parent component invoking the button
- states that are not needed have a default so they do not need to be included when initiating the component.

```js
// vendor
import React, { FC, MouseEvent } from 'react';
// css
import './Button.scss';


interface Props {
  click: (event: MouseEvent) => void;
  disabled?: boolean;
  text: string,
}

const Button: FC<Props> = ({
  click,
  disabled = false,
  text,
}: Props) => {

  return (
    <button
      className="Button"
      disabled={disabled}
      onClick={click}
    >
      {text}
    </button>
  )
}


export default Button;

```


*Bad*
- text cannot be changed
- the disabled state is using page specific data to make a decision if the button is disabled. This logic should be handled outside a generic component

```js
// vendor
import React, { FC, MouseEvent } from 'react';
// css
import './Button.scss';


interface Props {
  click: (event: MouseEvent) => void;
  users: Array<Users>;
}

const Button: FC<Props> = ({
  click,
  data,
}: Props) => {

  const isDisabled = users.filter((user => {
    return user.name !== 'admin';
  })).length > 0;

  return (
    <button
      className="Button"
      disabled={isDisabled}
      onClick={click}
    >
      Add Human
    </button>
  )
}


export default Button;

```
