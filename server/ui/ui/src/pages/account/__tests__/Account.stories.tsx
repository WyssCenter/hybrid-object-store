import React from 'react';
// context
import AppContext from 'Src/AppContext';
// component
import Account from '../Account';

export default {
  title: 'Pages/dataset/Account',
  component: Account,
  argTypes: {},
};

const user = {
  profile: {
    family_name: 'Test',
    given_name: 'Test',
    name: 'admin',
    email: 'admin@example.com',
    role: 'admin',
  }
};

const Template = (args) => (
  <AppContext.Provider value={{ user }}>
    <Account {...args} />
  </AppContext.Provider>
);

export const AccountComp = Template.bind({});

AccountComp.args = {};


AccountComp.parameters = {
  jest: ['Account.test.js'],
};
