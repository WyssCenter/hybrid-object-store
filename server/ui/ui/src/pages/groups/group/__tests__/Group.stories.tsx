// vendor
import React from 'react';
import { MemoryRouter } from 'react-router-dom';
// context
import AppContext from '../../../../AppContext';
// components
import Group from '../Group';

const user = { profile: { name: 'admin', role: 'admin' }};

export default {
  title: 'Pages/groups/group/Group',
  component: Group,
  argTypes: {},
  parameters: {
   actions: {
     handles: ['mouseover', 'click .button'],
   },
 },
};

const Template = (args) => (
  <AppContext.Provider value={{ user }}>
    <MemoryRouter
      initialEntries={["/groups", "/group"]}
      initialIndex={1}
    >
      <Group />
    </MemoryRouter>
  </AppContext.Provider>
);

export const GroupComp = Template.bind({});

GroupComp.args = {

};

GroupComp.parameters = {
  jest: ['Group.test.js'],
};
