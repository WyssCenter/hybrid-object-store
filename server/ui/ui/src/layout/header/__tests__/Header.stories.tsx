import React from 'react';

import Header from '../Header';

import { Story, Meta } from '@storybook/react/types-6-0';


export default {
  title: 'Layout/Header',
  component: Header,
  argTypes: {},
} as Meta;

const Template:Story = (args) => <Header {...args} />;

export const HeaderComp = Template.bind({});

HeaderComp.args = {
  location: {
    pathname: '/',
  }
};


HeaderComp.parameters = {
  jest: ['Header.test.js'],
};
