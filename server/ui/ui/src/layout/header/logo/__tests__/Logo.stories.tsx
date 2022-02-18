import React from 'react';
import { Story, Meta } from '@storybook/react/types-6-0';

import Logo from '../Logo';

export default {
  title: 'Layout/Header/Logo',
  component: Logo,
  argTypes: {},
} as Meta;

const Template: Story = (args) => <Logo {...args} />;

export const LogoComp = Template.bind({});

LogoComp.args = {};


LogoComp.parameters = {
  jest: ['Logo.test.js'],
};
