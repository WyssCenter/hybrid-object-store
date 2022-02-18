// vendor
import React from 'react';
// components
import FlatIconTextButton from '../FlatIconTextButton';
import { faCopy } from '@fortawesome/free-solid-svg-icons';

export default {
  title: 'Components/button/special/FlatIconTextButton',
  component: FlatIconTextButtonComp,
  argTypes: {
    click: Function,
    disabled: Boolean,
    color: String,
  },
};



const Template = (args) => (<div>
  <div id="FlatIconTextButton" />
    <FlatIconTextButton {...args} />
</div>);

export const FlatIconTextButtonComp = Template.bind({});

FlatIconTextButtonComp.args = {
  click: () => null,
  disabled: false,
  icon: faCopy,
  color: 'white'
};


FlatIconTextButtonComp.parameters = {
  jest: ['FlatIconTextButton.test.js'],
};
