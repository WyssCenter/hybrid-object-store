// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
// context
import AppContext from '../../../AppContext';
// components
import GroupsRouter from '../GroupsRouter';

const user = { profile: { name: 'admin' }};

describe('Groups Router', () => {
  beforeAll(() => {
    const user = JSON.stringify({ id_token: 'id_token' });
    localStorage.setItem(`oidc.user:http://localhost/auth/v1/.well-known/openid-configuration:HossServer`, user);
  })

  it('Renders groups page', () => {
    render(
      <AppContext.Provider value={{ user }}>
        <MemoryRouter
          initialEntries={["/groups", "/groups/group-name"]}
          initialIndex={0}
        >
          <GroupsRouter />
        </MemoryRouter>
      </AppContext.Provider>
    );
    const linkElement = screen.queryAllByText(/Groups/i);
    expect(linkElement[0]).toBeInTheDocument();
  })


  it('Renders Group', () => {
    render(
      <AppContext.Provider value={{ user }}>
        <MemoryRouter
          initialEntries={["/groups", "/groups/group-name"]}
          initialIndex={1}
        >
          <GroupsRouter />
        </MemoryRouter>
      </AppContext.Provider>
    );
    const linkElement = screen.getByText(/groupname/i);
    expect(linkElement).toBeInTheDocument();
  })
});
