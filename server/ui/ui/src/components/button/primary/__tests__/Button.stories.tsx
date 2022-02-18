// vendor
import React from 'react';
// components
import Button from '../Button';

export default {
  title: 'Components/button/Button',
  component: Button,
  argTypes: {
    click: Function,
    disabled: Boolean,
    text: String,
  },
};



const Template = (args) => (<div>
  <div id="Button" />
    <Button {...args} />
</div>);

export const ButtonComp = Template.bind({});

ButtonComp.args = {
  click: () => null,
  disabled: false,
  text: "submit",
};


ButtonComp.parameters = {
  jest: ['Button.test.js'],
};
