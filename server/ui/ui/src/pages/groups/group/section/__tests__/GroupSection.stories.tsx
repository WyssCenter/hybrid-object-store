// vendor
import React from 'react';
import { MemoryRouter } from 'react-router-dom';
// context
import AppContext from 'Src/AppContext';
// components
import GroupSection from '../GroupSection';

// data
import mockGroupData from './GroupData';

const user = { profile: { name: 'admin', role: 'admin' }};

export default {
  title: 'Pages/groups/group/GroupSection',
  component: GroupSection,
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
      <GroupSection {...args} />
    </MemoryRouter>
  </AppContext.Provider>
);

export const GroupSectionComp = Template.bind({});

GroupSectionComp.args = {
  group: mockGroupData
};

GroupSectionComp.parameters = {
  jest: ['GroupSection.test.js'],
};
