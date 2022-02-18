import React from 'react';

import Error from '../Error';

export default {
  title: 'Auth/Pages/Error',
  component: Error,
  argTypes: {},
};

const Template = (args) => <Error {...args} />;

export const ErrorComp = Template.bind({});

ErrorComp.args = {};


ErrorComp.parameters = {
  jest: ['Error.test.js'],
};
