// vendor
import React from 'react';
// components
import IconButton from '../IconButton';
import { faCopy } from '@fortawesome/free-solid-svg-icons';

export default {
  title: 'Components/button/IconButton',
  component: IconButton,
  argTypes: {
    click: Function,
    disabled: Boolean,
    color: String,
  },
};



const Template = (args) => (<div>
  <div id="IconButton" />
    <IconButton {...args} />
</div>);

export const IconButtonComp = Template.bind({});

IconButtonComp.args = {
  click: () => null,
  disabled: false,
  icon: faCopy,
  color: 'white'
};


IconButtonComp.parameters = {
  jest: ['IconButton.test.js'],
};
