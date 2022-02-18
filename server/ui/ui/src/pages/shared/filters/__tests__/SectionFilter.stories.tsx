// vendor
import React from 'react';
import SectionFilter from '../SectionFilter';
// data

export default {
  title: 'Pages/namespace/SectionFilter',
  component: SectionFilter,
  argTypes: {
    dataset: Object,
  },
};

const Template = (args) => (
  <SectionFilter
    {...args}
  />
)

export const SectionFilterComp = Template.bind({});

SectionFilterComp.args = {
  list: [],
  modalClose: () =>  null,
  modalVisible: false,
  openModal: () => null,
  permissions: true,
  postRoute: "",
  formattedSection: "Dataset",
  section: "dataset",
  sendRefetch: () => null,
  updateList: () => null,
};

SectionFilterComp.parameters = {
  jest: ['Namespace.test.js'],
};
