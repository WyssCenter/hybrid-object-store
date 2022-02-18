//vendor
import React from 'react';
// components
import SecondaryButton from '../SecondaryButton';

export default {
  title: 'Components/button/SecondaryButton',
  component: SecondaryButton,
  argTypes: {
    click: Function,
    disabled: Boolean,
    text: String,
  },
};



const Template = (args) => (<div>
  <div id="SecondaryButton" />
    <SecondaryButton {...args} />
</div>);

export const SecondaryButtonComp = Template.bind({});

SecondaryButtonComp.args = {
  click: () => null,
  disabled: false,
  text: "submit",
};


SecondaryButtonComp.parameters = {
  jest: ['SecondaryButton.test.js'],
};
