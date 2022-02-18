// vendor
import React from 'react';

import InputText from '../InputText';

export default {
  title: 'Components/form/text/InputText',
  component: InputText,
  argTypes: {},
};



const Template = (args) => <InputText {...args} />;

export const InputTextComp = Template.bind({});

InputTextComp.args = {
  css: 'small',
  inputRef: { current: <div />},
  label: 'Add User',
  placeholder: 'Username',
  updateValue: () => null,
};


InputTextComp.parameters = {
  jest: ['InputText.test.js'],
};
