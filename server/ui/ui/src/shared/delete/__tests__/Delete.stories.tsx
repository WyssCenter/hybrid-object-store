// vendor
import React from 'react';
// components
import Delete from '../Delete';

export default {
  title: 'Pages/dataset/Delete',
  component: Delete,
  argTypes: {
    datasetName: string,
    name: string,
    namespace: string,
    sectionType: string,
  },
};

const Template = (args) => <Delete {...args} />;

export const DeleteComp = Template.bind({});

DeleteComp.args = {
  datasetName: "my-dataset",
  name: "user",
  namespace: "my-namespace",
  sectionType: "user"
};


DeleteComp.parameters = {
  jest: ['Delete.test.js'],
};
