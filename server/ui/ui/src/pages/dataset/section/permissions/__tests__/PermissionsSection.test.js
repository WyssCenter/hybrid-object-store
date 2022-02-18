// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
// context
import AppContext from 'Src/AppContext';
// components
import PermissionsSection from '../PermissionsSection';
// data
import { users, groups } from './PermissionsSectionData';

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

describe('PermissionsSection', () => {
  it('Renders Users', () => {
    render(
      <AppContext.Provider value={{user: adminRole}}>
        <PermissionsSection
          list={users}
          sectionType="user"

        />
      </AppContext.Provider>
    );
    const name = screen.getByText(/admin/i);
    expect(name).toBeInTheDocument();
  });


  it('Renders Groups', () => {
    render(
      <AppContext.Provider value={{user: adminRole}}>

        <PermissionsSection
          list={groups}
          sectionType="group"

        />
      </AppContext.Provider>

    );
    const name = screen.getByText(/alpha/i);
    expect(name).toBeInTheDocument();
  });

  it('Admin user has access to add and remove permissions', async () => {
    render(
      <AppContext.Provider value={{user: adminRole}}>

        <PermissionsSection
          list={groups}
          sectionType="group"

        />
      </AppContext.Provider>
    );


    const buttons = screen.queryAllByRole('button');

    expect(buttons.length).toBe(3);

  });

  it('Privileged user has access to add and remove permissions', async () => {
    render(
      <AppContext.Provider value={{user: privilegedRole}}>

        <PermissionsSection
          list={groups}
          sectionType="group"

        />
      </AppContext.Provider>
    );


    const buttons = screen.queryAllByRole('button');

    expect(buttons.length).toBe(3);

  });

  it('Regular user has no access to add or remove permissions', async () => {
    render(
      <AppContext.Provider value={{user: userRole}}>

        <PermissionsSection
          list={groups}
          sectionType="group"

        />
      </AppContext.Provider>
    );


    const buttons = screen.queryAllByRole('button');

    expect(buttons.length).toBe(0);

  });
})
