import React from 'react';

import Dataset from '../Dataset';

jest.mock('Environment/createEnvironment', () => {
  return {
    get: () => new Promise((resolve) => {
      resolve(
        {
          json: () => new Promise((resolve) => { resolve(mockDatasetData)})
        }
      )
    })
  }
})

export default {
  title: 'Pages/dataset/Dataset',
  component: Dataset,
  argTypes: {},
};

const Template = (args) => <Dataset {...args} />;

export const DatasetComp = Template.bind({});

DatasetComp.args = {
};


DatasetComp.parameters = {
  jest: ['Dataset.test.js'],
};
