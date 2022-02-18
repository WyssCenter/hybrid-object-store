import React from 'react';

import Error from '../Error';

export default {
  title: 'Shared/modal/create/Error',
  component: Error,
  argTypes: {
    errorMessage: 'This is an error',
    name: 'namespace-12',
    send: () => null,
  },
};



const Template = (args) => (<div>
  <div id="Error" />
    <Error
      {...args}
    >
      <div>
        Error JSX Context
      </div>
    </Error>
</div>);

export const ErrorComp = Template.bind({});

ErrorComp.args = {
  errorMessage: 'This is an error',
  name: 'namespace-12',
  send: () => null,
};


ErrorComp.parameters = {
  jest: ['Error.test.js'],
};
