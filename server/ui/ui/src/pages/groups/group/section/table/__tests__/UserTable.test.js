// vendor
import React from 'react';
import { MemoryRouter } from 'react-router-dom';
import { render, screen } from '@testing-library/react';
// context
import AppContext from 'Src/AppContext';
// data
import userData from '../../../__tests__/GroupData'
// components
import UserTable from '../UserTable';


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


describe('UserTable', () => {
  it('Renders table header', () => {
    render(
      <AppContext.Provider value={{ user: adminRole }}>
        <MemoryRouter
          initialEntries={["/groups/group-1"]}
          initialIndex={0}
        >
          <UserTable group={userData} />
        </MemoryRouter>
      </AppContext.Provider>
    );

    expect(screen.getByText(/Username/i)).toBeInTheDocument();
    expect(screen.getByText(/Name/)).toBeInTheDocument();
    expect(screen.getByText(/Role/i)).toBeInTheDocument();
  })

  it('Row has all elements', () => {
    render(
      <AppContext.Provider value={{ user: adminRole }}>
        <MemoryRouter
          initialEntries={["/groups/group-1"]}
          initialIndex={0}
        >
          <UserTable group={userData} />
        </MemoryRouter>
      </AppContext.Provider>
    );

    expect(screen.queryAllByRole('button').length).toBe(2);
    expect(screen.queryAllByText('admin').length).toBe(3);
    expect(screen.queryAllByText('privileged').length).toBe(3);
  });


  it('Admin user role actions appear in table header', () => {
    render(
      <AppContext.Provider value={{ user: adminRole }}>
        <MemoryRouter
          initialEntries={["/groups/group-1"]}
          initialIndex={0}
        >
          <UserTable group={userData} />
        </MemoryRouter>
      </AppContext.Provider>
    );

    expect(screen.queryAllByText(/actions/i).length).toBe(1);
  });


  it('Privileged user role actions appear in table header', () => {
    render(
      <AppContext.Provider value={{ user: privilegedRole }}>
        <MemoryRouter
          initialEntries={["/groups/group-1"]}
          initialIndex={0}
        >
          <UserTable group={userData} />
        </MemoryRouter>
      </AppContext.Provider>
    );

    expect(screen.queryAllByText(/actions/i).length).toBe(1);
  });


  it('Regular user role actions does not appear in table header', () => {
    render(
      <AppContext.Provider value={{ user: userRole }}>
        <MemoryRouter
          initialEntries={["/groups/group-1"]}
          initialIndex={0}
        >
          <UserTable group={userData} />
        </MemoryRouter>
      </AppContext.Provider>
    );

    expect(screen.queryAllByText(/actions/i).length).toBe(0);
  });

})
