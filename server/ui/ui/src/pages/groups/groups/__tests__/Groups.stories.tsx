// vendor
import React from 'react';
import { MemoryRouter } from 'react-router-dom';
// context
import AppContext from '../../../../AppContext';
// components
import Groups from '../Groups';

const user = { profile: { name: 'admin' }};

export default {
  title: 'Pages/groups/Groups',
  component: Groups,
  argTypes: {},
};

const Template = (args) => (
  <AppContext.Provider value={{ user }}>
    <MemoryRouter
      initialEntries={["/groups", "/group"]}
      initialIndex={0}
    >
      <Groups />
    </MemoryRouter>
  </AppContext.Provider>
);

export const GroupsComp = Template.bind({});

GroupsComp.args = {

};

GroupsComp.parameters = {
  jest: ['Groups.test.js'],
};
