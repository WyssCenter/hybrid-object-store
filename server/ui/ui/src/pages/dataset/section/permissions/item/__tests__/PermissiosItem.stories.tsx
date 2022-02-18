// vendor
import React from 'react';
// components
import PermissionsItem from '../PermissionsItem';

export default {
  title: 'Pages/dataset/PermissionsItem',
  component: PermissionsItem,
  argTypes: {
    item: Object,
    sectionType: String,
  },
};

const Template = (args) => <PermissionsItem {...args} />;

export const PermissionsItemComp = Template.bind({});

PermissionsItemComp.args = {
  item: {
    permission: 'r',
    group: {
      group_name: 'user',
    }
  },
  sectionType: 'user'
};


PermissionsItemComp.parameters = {
  jest: ['PermissionsItem.test.js'],
};
