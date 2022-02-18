// vendor
import React from 'react';
import { MemoryRouter } from 'react-router-dom';
// components
import NamespaceDetails from '../NamespaceDetails';

export default {
  title: 'Pages/namespaces/NamespaceDetails',
  component: NamespaceDetails,
  argTypes: {
    namespaceName: String,
    isExpanded: Boolean
  },
};

const Template = (args) => (

  <MemoryRouter
    initialEntries={['/default']}
    initialIndex={0}
  >
    <NamespaceDetails
      namespaceName="namespaceName"
      isExpanded={true}
    />
  </MemoryRouter>
)

export const NamespaceDetailsComp = Template.bind({});

NamespaceDetailsComp.args = {
  namespaceName: 'namespaceName',
  isExpanded: true
};

NamespaceDetailsComp.parameters = {
  jest: ['Namespace.test.js'],
};
