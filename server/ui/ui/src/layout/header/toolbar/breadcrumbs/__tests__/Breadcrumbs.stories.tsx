import React from 'react';

import Breadcrumbs from '../Breadcrumbs';
import { Story, Meta } from '@storybook/react/types-6-0';

export default {
  title: 'Layout/Header/Toolbar/Breadcrumbs',
  component: Breadcrumbs,
  argTypes: {},
} as Meta;

const Template: Story = (args) => <Breadcrumbs {...args} />;

export const BreadcrumbsComp = Template.bind({});

BreadcrumbsComp.args = {
  location: {
    pathname: '/namespace/datasetName',
  }
};


BreadcrumbsComp.parameters = {
  jest: ['Breadcrumbs.test.js'],
};
