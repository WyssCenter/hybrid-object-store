import React from 'react';

import UserInput from '../UserInput';

export default {
  title: 'Pages/dataset/UserInput',
  component: UserInput,
  argTypes: {},
};

const Template = (args) => <UserInput {...args} />;

export const UserInputComp = Template.bind({});

UserInputComp.args = {
  inputRef: { current: <input /> },
  permissionType: "user",
  updateName: () => null
};


UserInputComp.parameters = {
  jest: ['UserInput.test.js'],
};
