// vendor
import React from 'react';
// components
import ExpandButton from '../ExpandButton';

export default {
  title: 'Components/button/ExpandButton',
  component: ExpandButton,
  argTypes: {
    click: Function,
    disabled: Boolean,
    text: String,
  },
};



const Template = (args) => (<div>
  <div id="ExpandButton" />
    <ExpandButton {...args} />
</div>);

export const ExpandButtonComp = Template.bind({});

ExpandButtonComp.args = {
  click: () => null,
  disabled: false,
  text: "submit",
};


ExpandButtonComp.parameters = {
  jest: ['ExpandButton.test.js'],
};
