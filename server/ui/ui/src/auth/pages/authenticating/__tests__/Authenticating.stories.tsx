import React from 'react';

import Authenticating from '../Authenticating';

export default {
  title: 'Auth/Pages/Authenticating',
  component: Authenticating,
  argTypes: {},
};

const Template = (args) => <Authenticating {...args} />;

export const AuthenticatingComp = Template.bind({});

AuthenticatingComp.args = {};


AuthenticatingComp.parameters = {
  jest: ['Authenticating.test.js'],
};
