// vendor
import React from 'react';
import { act } from 'react-dom/test-utils';
import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
// conte
import AppContext from 'Src/AppContext';
// data
import mockGroupData from './GroupData';
// components
import GroupSection from '../GroupSection';



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

describe('Group', () => {

  it('Admin user can see add user form', async () => {
    render(
      <AppContext.Provider value={{ user: adminRole }}>
        <MemoryRouter
          initialEntries={["/groups"]}
          initialIndex={0}
        >
          <GroupSection group={mockGroupData} />
        </MemoryRouter>
      </AppContext.Provider>
    );

    const textboxElement = screen.queryAllByRole('textbox');
    expect(textboxElement.length).toBe(1);
  });


  it('Privileged user can see add user form', async () => {
    render(
      <AppContext.Provider value={{ user: privilegedRole }}>
        <MemoryRouter
          initialEntries={["/groups"]}
          initialIndex={0}
        >
          <GroupSection group={mockGroupData} />
        </MemoryRouter>
      </AppContext.Provider>
    );

    const textboxElement = screen.queryAllByRole('textbox');
    expect(textboxElement.length).toBe(1);
  });


  it('Regular user can not see add user form', async () => {
    render(
      <AppContext.Provider value={{ user: userRole }}>
        <MemoryRouter
          initialEntries={["/groups"]}
          initialIndex={0}
        >
          <GroupSection group={mockGroupData} />
        </MemoryRouter>
      </AppContext.Provider>
    );

    const textboxElement = screen.queryAllByRole('textbox');
    expect(textboxElement.length).toBe(0);
  })
});
