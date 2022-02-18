// vendor
import React from 'react';
// machine
import authMachine from '../auth/machine/AuthStateMachine';
// components
import App from '../App';

export default {
  title: 'Hoss',
  component: App,
  argTypes: {},
};

const Template = (args) => <App {...args} />;

export const Application = Template.bind({});

Application.args = {
  machine: authMachine,
};


Application.parameters = {
  jest: ['App.test.js'],
};
