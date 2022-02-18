// vendor
import React from 'react';
import { MemoryRouter } from 'react-router-dom';
// context
import AppContext from 'Src/AppContext';
import NamespaceListingContext from '../NamespaceListingContext'
// components
import NamespaceListing from '../NamespaceListing';


const adminRole = {
  profile: {
    name: 'admin',
    role: 'admin'
  }
}

export default {
  argTypes: {},
  component: NamespaceListing,
  title: 'Pages/namespaces/NamespaceListing',
};

const Template = (args) => (
  <AppContext.Provider value={{user: adminRole}}>
    <NamespaceListingContext.Provider value={{send: jest.fn()}}>
      <MemoryRouter
        initialEntries={['/default']}
        initialIndex={0}
      >
        <NamespaceListing />
      </MemoryRouter>
    </NamespaceListingContext.Provider>
  </AppContext.Provider>
);

export const NamespaceListingComp = Template.bind({});

NamespaceListingComp.args = {};


NamespaceListingComp.parameters = {
  jest: ['NamespaceListing.test.js'],
};
