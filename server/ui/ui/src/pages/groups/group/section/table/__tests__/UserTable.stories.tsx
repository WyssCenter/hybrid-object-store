// vendor
import React from 'react';
import { MemoryRouter } from 'react-router-dom';
// context
import AppContext from 'Src/AppContext';
// component
import UserTable from '../UserTable';
// data
import userData from '../../../__tests__/GroupData';

export default {
  title: 'Pages/groups/group/section/table/UserTable',
  component: UserTable,
  argTypes: {
    user: Object,
  },
};


const adminRole = {
  profile: {
    name: 'admin',
    role: 'admin'
  }
}

const Template = (args) => (
  <AppContext.Provider value={{ user: adminRole }}>
    <MemoryRouter
      initialEntries={["/groups/group-1"]}
      initialIndex={0}
    >
      <UserTable {...args} />
    </MemoryRouter>
  </AppContext.Provider>
);

export const UserTableComp = Template.bind({});

UserTableComp.args = {
  user: userData,
};


UserTableComp.parameters = {
  jest: ['UserTable.test.js'],
};
