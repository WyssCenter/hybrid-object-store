// vendor
import React from 'react';
import { MemoryRouter } from 'react-router-dom';
// context
import AppContext from 'Src/AppContext';
import GroupContext from '../../../../GroupContext';
// component
import UserItem from '../UserItem';



const adminRole = {
  profile: {
    name: 'admin',
    role: 'admin'
  }
}

const membership = {"group":{"group_name":"group-1","description":"this is a description"}}

const send = () => null;

export default {
  title: 'Pages/groups/group/section/table/item/UserItem',
  component: UserItem,
  argTypes: {
    membership: Object,
  },
};

const Template = (args) => (
  <AppContext.Provider value={{ user: adminRole }}>
    <GroupContext.Provider value={{ send, groupname: 'group-1' }}>
      <MemoryRouter>
        <UserItem  {...args} />
      </MemoryRouter>
    </GroupContext.Provider>
  </AppContext.Provider>
);

export const UserItemComp = Template.bind({});

UserItemComp.args = {
  membership,
};


UserItemComp.parameters = {
  jest: ['UserItem.test.js'],
};
