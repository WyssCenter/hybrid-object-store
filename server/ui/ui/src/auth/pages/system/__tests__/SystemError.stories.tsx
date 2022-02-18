// vendor
import React from 'react';
// components
import SystemError from '../SystemError';

export default {
  title: 'Auth/Pages/SystemError',
  component: SystemError,
  argTypes: {
    errorMessage: String,
    send: Function,
  },
};

const Template = (args) => <SystemError {...args} />;

export const SystemErrorComp = Template.bind({});

SystemErrorComp.args = {
  errorMessage: 'Cannot load auth server',
  send: () => null,
};

SystemErrorComp.parameters = {
  jest: ['SystemError.test.js'],
};
