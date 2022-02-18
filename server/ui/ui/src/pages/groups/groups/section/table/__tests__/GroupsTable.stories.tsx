// vendor
import React from 'react';
import { MemoryRouter } from 'react-router-dom';
// component
import GroupsTable from '../GroupsTable';
// data
import userData from '../../../__tests__/GroupsData';

export default {
  title: 'Pages/groups/groups/table/GroupsTable',
  component: GroupsTable,
  argTypes: {
    user: Object,
  },
};

const Template = (args) => (
  <MemoryRouter>
    <GroupsTable {...args} />
  </MemoryRouter>
);

export const GroupsTableComp = Template.bind({});

GroupsTableComp.args = {
  user: userData,
};


GroupsTableComp.parameters = {
  jest: ['GroupsTable.test.js'],
};
