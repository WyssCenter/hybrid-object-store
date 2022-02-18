// vendor
import React from 'react';
import { render, screen,  fireEvent } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
// context
import AppContext from '../../../../../AppContext';
// components
import User from '../User';
// data
import context from './UserData';

const UserWrapper = () => {
  return (
    <MemoryRouter>
      <AppContext.Provider value={context}>
        <User
          send={(stateValue) => null}
        />
      </AppContext.Provider>
    </MemoryRouter>
  )
}

describe('User', () => {
  it('Renders username', async () => {
    await render(
      <UserWrapper />
    );
    const username = screen.getByText(/admin/i);
    expect(username).toBeInTheDocument();
  });

  it('Menu opens and renders all buttons for admin users', async () => {
    await render(
      <UserWrapper />
    );
    const menuButton = screen.getByRole('button');
    fireEvent.click(menuButton, { target: { value: 'name' }});

    const accountButton = screen.getByText(/Account/i);
    const groupButton = screen.getByText(/Group/i);
    const tokenButton = screen.getByText(/Token/i);
    const logoutButton = screen.getByText(/Logout/i);


    expect(accountButton).toBeInTheDocument();
    expect(groupButton).toBeInTheDocument();
    expect(tokenButton).toBeInTheDocument();
    expect(logoutButton).toBeInTheDocument();
  });



  it('Menu opens and clicking item closes menu', async () => {
    await render(
      <UserWrapper />
    );
    const menuButton = screen.getByRole('button');
    fireEvent.click(menuButton, { target: { value: 'name' }});

    const accountButton = screen.getByText(/Account/i);


    expect(accountButton).toBeInTheDocument();


    fireEvent.click(accountButton, { target: { value: 'name' }});

    expect(accountButton).not.toBeInTheDocument();
  });

})
