// vendor
import React from 'react';
import { MemoryRouter } from 'react-router-dom';
// components
import NamespacesTable from '../NamespacesTable';
// context
import AppContext from 'Src/AppContext';
// data
import mockNamespacesTableData from '../__tests__/NamespacesTableData';

const adminRole = {
  profile: {
    name: 'admin',
    role: 'admin'
  }
}

const privilegedRole = {
  profile: {
    name: 'privileged',
    role: 'privileged'
  }
}

const userRole = {
  profile: {
    name: 'user',
    role: 'user'
  }
}

export default {
  title: 'Pages/namespace/NamespacesTable',
  component: NamespacesTable,
  argTypes: {
    dataset: Object,
  },
};

const Template = (args) => (
  <AppContext.Provider value={{user: adminRole}}>
    <MemoryRouter
      initialEntries={['/']}
      initialIndex={0}
    >
      <NamespacesTable
        {...args}
      />
    </MemoryRouter>
  </AppContext.Provider>
)

export const NamespacesTableComp = Template.bind({});

NamespacesTableComp.args = {
  datasets: mockNamespacesTableData
};

NamespacesTableComp.parameters = {
  jest: ['Namespace.test.js'],
};
