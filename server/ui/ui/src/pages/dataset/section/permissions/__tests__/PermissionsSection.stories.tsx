import React from 'react';

import PermissionsSection from '../PermissionsSection';

export default {
  title: 'Pages/dataset/PermissionsSection',
  component: PermissionsSection,
  argTypes: {},
};

const Template = (args) => <PermissionsSection {...args} />;

export const PermissionsSectionComp = Template.bind({});

PermissionsSectionComp.args = {};


PermissionsSectionComp.parameters = {
  jest: ['PermissionsSection.test.js'],
};
