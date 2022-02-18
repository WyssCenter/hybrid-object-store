// vendor
import React from 'react';
// components
import WarningButton from '../WarningButton';

export default {
  title: 'Components/button/WarningButton',
  component: WarningButton,
  argTypes: {
    click: Function,
    disabled: Boolean,
    text: String,
  },
};



const Template = (args) => (<div>
  <div id="WarningButton" />
    <WarningButton {...args} />
</div>);

export const WarningButtonComp = Template.bind({});

WarningButtonComp.args = {
  click: () => null,
  disabled: false,
  text: "submit",
};


WarningButtonComp.parameters = {
  jest: ['WarningButton.test.js'],
};
