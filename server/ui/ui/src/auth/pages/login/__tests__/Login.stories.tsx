import React from 'react';

import Login from '../Login';

export default {
  title: 'Auth/Pages/Login',
  component: Login,
  argTypes: {},
};

const Template = (args) => <Login {...args} />;

export const LoginComp = Template.bind({});

LoginComp.args = {};


LoginComp.parameters = {
  jest: ['Login.test.js'],
};
