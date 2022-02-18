// vendor
import React from 'react';
import { MemoryRouter } from 'react-router-dom';
// context
import AppContext from 'Src/AppContext';
// components
import NamespaceRow from '../NamespaceRow';
// data
import mockNamespacesTableData from '../../__tests__/NamespacesTableData';



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
  title: 'Pages/namespaces/row/NamespaceRow',
  component: NamespaceRow,
  argTypes: {
    namespace: Object,
  },
};

const Template = (args) => (
  <AppContext.Provider value={{user: adminRole}}>
    <MemoryRouter
      initialEntries={['/default']}
      initialIndex={0}
    >
      <NamespaceRow
        {...args}
      />
    </MemoryRouter>
  </AppContext.Provider>
);

export const DatasetItemComp = Template.bind({});

DatasetItemComp.args = {
  namespace: mockNamespacesTableData[0]
};

DatasetItemComp.parameters = {
  jest: ['Namespace.test.js'],
};
