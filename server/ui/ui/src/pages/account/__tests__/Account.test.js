// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
// context
import AppContext from 'Src/AppContext';
// components
import Account from '../Account';

const user = {
  profile: {
    family_name: '*',
    given_name: '*',
    name: 'admin',
    email: 'admin@example.com',
    role: 'admin',
  }
};

const namedUser = {
  profile: {
    family_name: 'Test',
    given_name: 'Test',
    name: 'admin',
    email: 'admin@example.com',
    role: 'admin',
  }
};

describe('Account', () => {
  it('Renders Account', async () => {
    render(
      <AppContext.Provider value={{ user }}>
        <Account
        />
      </AppContext.Provider>
    );
    const AccountHeader = screen.getByText('Account');

    expect(AccountHeader).toBeInTheDocument();

  });

  it('Account section will ignore first/last name if value is *', async () => {
    render(
      <AppContext.Provider value={{ user }}>
        <Account
        />
      </AppContext.Provider>
    );
    const firstName = screen.queryAllByText('First Name:');
    expect(firstName.length).toBe(0);
    const lastName = screen.queryAllByText('Last Name:');
    expect(lastName.length).toBe(0);

  });

  it('First/last name will appear if set', async () => {
    render(
      <AppContext.Provider value={{ user: namedUser }}>
        <Account
        />
      </AppContext.Provider>
    );
    const firstName = screen.queryAllByText('First Name:');
    expect(firstName.length).toBe(1);
    const lastName = screen.queryAllByText('Last Name:');
    expect(lastName.length).toBe(1);

  });
});
