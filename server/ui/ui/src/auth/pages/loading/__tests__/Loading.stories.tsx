import React from 'react';

import Loading from '../Loading';

export default {
  title: 'Auth/Pages/Loading',
  component: Loading,
  argTypes: {},
};

const Template = (args) => <Loading {...args} />;

export const LoadingComp = Template.bind({});

LoadingComp.args = {};


LoadingComp.parameters = {
  jest: ['Loading.test.js'],
};
