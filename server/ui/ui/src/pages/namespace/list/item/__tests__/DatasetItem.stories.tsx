// vendor
import React from 'react';
// components
import DatasetItem from '../DatasetItem';
// data
import mockDatsetDetails from './DatasetItemData';

export default {
  title: 'Pages/namespace/DatasetItem',
  component: DatasetItem,
  argTypes: {
    dataset: Object,
  },
};

const Template = (args) => <DatasetItem {...args} />;

export const DatasetItemComp = Template.bind({});

DatasetItemComp.args = {
  dataset: mockDatsetDetails
};

DatasetItemComp.parameters = {
  jest: ['Namespace.test.js'],
};
