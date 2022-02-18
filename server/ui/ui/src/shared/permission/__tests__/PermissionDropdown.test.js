// vendor
import fetchMock, { enableFetchMocks } from 'jest-fetch-mock';
import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
// context
import AppContext from 'Src/AppContext';
// components
import PermissionDropdown from '../PermissionDropdown';

enableFetchMocks();

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


describe('PermissionDropdown', () => {
  it('Renders dropdown', () => {
    render(
      <AppContext.Provider value={{user: adminRole}}>
        <PermissionDropdown
          name="user"
          permission={{name: "r"}}
          updatePermissions={() => null}
        />
      </AppContext.Provider>
    );
    const dropdownText = screen.getByText(/r/i);
    expect(dropdownText).toBeInTheDocument();
  });


  it('Updates dropdown to rw', () => {
    const updatePermissions = jest.fn();
    render(
      <AppContext.Provider value={{user: adminRole}}>
        <PermissionDropdown
          name="user"
          permission={{name: "r"}}
          updatePermissions={updatePermissions}
        />
      </AppContext.Provider>

    );
    const dropdown = screen.getByText('r');

    fireEvent.click(dropdown);

    const menuItem = screen.getByText('rw');

    fireEvent.click(menuItem);


    const newDropdown = screen.getByText('rw');


    expect(newDropdown).toBeInTheDocument();
  });

  it('Admin users have access to dropdown', () => {
    const updatePermissions = jest.fn();
    render(
      <AppContext.Provider value={{user: adminRole}}>
        <PermissionDropdown
          name="user"
          permission="r"
          updatePermissions={updatePermissions}
        />
      </AppContext.Provider>

    );
    const dropdownElement = screen.getByRole('presentation');
    expect(dropdownElement).toBeInTheDocument();
  });

  it('Privileged users have access to dropdown', () => {
    const updatePermissions = jest.fn();
    render(
      <AppContext.Provider value={{user: privilegedRole}}>
        <PermissionDropdown
          name="user"
          permissions={{name: "r"}}
          updatePermissions={updatePermissions}
        />
      </AppContext.Provider>

    );
    const dropdownElement = screen.getByRole('presentation');
    expect(dropdownElement).toBeInTheDocument();
  });

  it('Regular users do not have access to dropdown', () => {
    const updatePermissions = jest.fn();
    render(
      <AppContext.Provider value={{user: userRole}}>
        <PermissionDropdown
          name="user"
          permissions={{name: "r"}}
          updatePermissions={updatePermissions}
        />
      </AppContext.Provider>

    );
    expect(screen.queryAllByRole('presentation').length).toBe(0)

  });

})
