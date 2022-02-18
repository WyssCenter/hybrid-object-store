// vendor
import React from 'react';
// components
import DatasetSection from '../DatasetSection';
// data
import dataset from './DatasetSectionData';

export default {
  title: 'Pages/dataset/DatasetSection',
  component: DatasetSection,
  argTypes: {},
};

const Template = (args) => <DatasetSection {...args} />;

export const DatasetSectionComp = Template.bind({});

DatasetSectionComp.args = {
  dataset,
  datasetName:"my-dataset",
  namespace:"my-namespace"
};


DatasetSectionComp.parameters = {
  jest: ['Dataset.test.js'],
};
