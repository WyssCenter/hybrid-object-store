import React from 'react';

import Layout from '../Layout';

export default {
  title: 'Hoss',
  component: Layout,
  argTypes: {},
};

const Template = (args) => <Layout {...args} />;

export const LayoutComp = Template.bind({});

LayoutComp.args = {};


LayoutComp.parameters = {
  jest: ['Layout.test.js'],
};
