// vendor
import React from 'react';
// components
import TernaryButton from '../TernaryButton';

export default {
  title: 'Components/button/TernaryButton',
  component: TernaryButton,
  argTypes: {
    click: Function,
    disabled: Boolean,
    text: String,
  },

};



const Template = (args) => (<div>
  <div id="TernaryButton" />
    <TernaryButton {...args} />
</div>);

export const TernaryButtonComp = Template.bind({});

TernaryButtonComp.args = {
  click: () => null,
  disabled: false,
  text: "submit",
};


TernaryButtonComp.parameters = {
  jest: ['TernaryButton.test.js'],
};
