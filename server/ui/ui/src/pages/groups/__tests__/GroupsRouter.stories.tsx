// vendor
import React from 'react';
// components
import GroupsRouter from '../GroupsRouter';

export default {
  title: 'Pages/groups/GroupsRouter',
  component: GroupsRouter,
  argTypes: {},
};

const Template = (args) => <GroupsRouter {...args} />;

export const GroupsRouterComp = Template.bind({});

GroupsRouterComp.args = {};

GroupsRouterComp.parameters = {
  jest: ['GroupsRouter.test.js'],
};
