// vendor
import React from 'react';
import { MemoryRouter } from 'react-router-dom';
// context
import GroupsContext from '../../../../GroupsContext';
// component
import GroupItem from '../GroupItem';


const membership = {"group":{"group_name":"group-1","description":"this is a description"}}

const send = () => null;

export default {
  title: 'Pages/groups/groups/table/item/GroupItem',
  component: GroupItem,
  argTypes: {
    membership: Object,
  },
};

const Template = (args) => (
  <GroupsContext.Provider value={{ send }}>
    <MemoryRouter>
      <GroupItem  {...args} />
    </MemoryRouter>
  </GroupsContext.Provider>
);

export const GroupItemComp = Template.bind({});

GroupItemComp.args = {
  membership,
};


GroupItemComp.parameters = {
  jest: ['GroupItem.test.js'],
};
