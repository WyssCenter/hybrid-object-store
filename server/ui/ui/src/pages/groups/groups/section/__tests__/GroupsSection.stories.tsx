// vendor
import React from 'react';
import { MemoryRouter } from 'react-router-dom';
// context
import AppContext from '../../../../../AppContext';
// data
import mockGroupsSectionData from '../../__tests__/GroupsData';
// components
import GroupsSection from '../GroupsSection';

export default {
  title: 'Pages/groups/groups/section/GroupsSection',
  component: GroupsSection,
  argTypes: {},
};

const user = { profile: { name: 'admin' }};

const Template = (args) => (
  <AppContext.Provider value={{ user }}>
    <MemoryRouter
      initialEntries={["/groups"]}
      initialIndex={0}
    >
      <GroupsSection {...args} />
    </MemoryRouter>
  </AppContext.Provider>
);

export const GroupsSectionComp = Template.bind({});

GroupsSectionComp.args = {
  user: mockGroupsSectionData
};

GroupsSectionComp.parameters = {
  jest: ['GroupsSection.test.js'],
};
