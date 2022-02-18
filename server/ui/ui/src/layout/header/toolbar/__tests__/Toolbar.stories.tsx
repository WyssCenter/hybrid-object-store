import React from 'react';
import { Story, Meta } from '@storybook/react/types-6-0';
import Toolbar from '../Toolbar';

export default {
  title: 'Layout/Header/Toolbar',
  component: Toolbar,
  argTypes: {},
} as Meta;

const Template: Story = (args) => <Toolbar {...args} />;

export const ToolbarComp = Template.bind({});

ToolbarComp.args = {
  location: {
    pathname: '/namespace',
  }
};


ToolbarComp.parameters = {
  jest: ['Toolbar.test.js'],
};
