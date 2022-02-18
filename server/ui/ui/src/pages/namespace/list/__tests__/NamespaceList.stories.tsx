// vendor
import React from 'react';
// components
import NamespaceList from '../NamespaceList';
// data
import mockDatsetDetails from './NamespaceListData';

export default {
  title: 'Pages/namespace/NamespaceList',
  component: NamespaceList,
  argTypes: {
    dataset: Object,
  },
};

const Template = (args) => <NamespaceList {...args} />;

export const NamespaceListComp = Template.bind({});

NamespaceListComp.args = {
  datasets: mockDatsetDetails
};

NamespaceListComp.parameters = {
  jest: ['Namespace.test.js'],
};
