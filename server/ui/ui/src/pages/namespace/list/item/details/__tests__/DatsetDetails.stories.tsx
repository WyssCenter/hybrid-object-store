import React from 'react';

import DatasetDetails from '../DatasetDetails';

export default {
  title: 'Pages/namespace/DatasetDetails',
  component: DatasetDetails,
  argTypes: {
    datasetName: String,
    isExpanded: Boolean
  },
};

const Template = (args) => <DatasetDetails {...args} />;

export const DatasetDetailsComp = Template.bind({});

DatasetDetailsComp.args = {
  datasetName: 'datasetname',
  isExpanded: true
};

DatasetDetailsComp.parameters = {
  jest: ['Namespace.test.js'],
};
