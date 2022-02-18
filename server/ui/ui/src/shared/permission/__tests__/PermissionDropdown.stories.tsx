import React from 'react';

import PermissionDropdown from '../PermissionDropdown';

export default {
  title: 'Pages/dataset/PermissionDropdown',
  component: PermissionDropdown,
  argTypes: {},
};

const Template = (args) => <PermissionDropdown {...args} />;

export const PermissionDropdownComp = Template.bind({});

PermissionDropdownComp.args = {
  inputRef: { current: <input /> },
  permissionType: "user",
  updateName: () => null
};


PermissionDropdownComp.parameters = {
  jest: ['PermissionDropdown.test.js'],
};
