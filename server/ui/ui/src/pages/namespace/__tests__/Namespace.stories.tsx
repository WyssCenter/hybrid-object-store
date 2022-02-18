import React from 'react';

import Namespace from '../Namespace';

export default {
  title: 'Pages/namespace/Namespace',
  component: Namespace,
  argTypes: {},
};

const Template = (args) => <Namespace {...args} />;

export const NamespaceComp = Template.bind({});

NamespaceComp.args = {};

NamespaceComp.parameters = {
  jest: ['Namespace.test.js'],
};
