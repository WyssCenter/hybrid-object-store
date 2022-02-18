import React from 'react';

import AddSection from '../AddSection';

export default {
  title: 'Pages/dataset/AddSection',
  component: AddSection,
  argTypes: {},
};

const Template = (args) => <AddSection {...args} />;

export const AddSectionComp = Template.bind({});

AddSectionComp.args = {
  sectionType: "user"
};


AddSectionComp.parameters = {
  jest: ['AddSection.test.js'],
};
