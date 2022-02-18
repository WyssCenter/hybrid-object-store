// vendor
import React from 'react';
import { act } from 'react-dom/test-utils';
import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
// context
import AppContext from '../../../../../AppContext';
// data
import mockGroupsSectionData from '../../__tests__/GroupsData';
// components
import GroupsSection from '../GroupsSection';


const user = { profile: { name: 'admin' }};

const privilegedData = Object.assign({}, mockGroupsSectionData, { role: 'privileged' });

const userData = Object.assign({}, mockGroupsSectionData, { role: 'user' });

describe('GroupsSection', () => {

  it('Renders mockGroupsSection', async () => {
    act(() => {
      render(
        <AppContext.Provider value={{ user }}>
          <MemoryRouter
            initialEntries={["/groups"]}
            initialIndex={0}
          >
            <GroupsSection user={mockGroupsSectionData} />
          </MemoryRouter>
        </AppContext.Provider>
      );
    });


    // await waitFor(() => jest.mock);

    const tableItemElement = screen.getByText(/group-1/i);
    expect(tableItemElement).toBeInTheDocument();
  });


  it('Admin is not in table', async () => {
    act(() => {
      render(
        <AppContext.Provider value={{ user }}>
          <MemoryRouter
            initialEntries={["/groups"]}
            initialIndex={0}
          >
            <GroupsSection user={mockGroupsSectionData}  />
          </MemoryRouter>
        </AppContext.Provider>
      );
    });

    await waitFor(() => jest.mock);

    const listElement = screen.queryAllByText(/admin/i);
    expect(listElement.length).toBe(0);
  })


  it('Renders create section', async () => {
    act(() => {
      render(
        <AppContext.Provider value={{ user }}>
          <MemoryRouter
            initialEntries={["/groups"]}
            initialIndex={0}
          >
            <GroupsSection user={mockGroupsSectionData}  />
          </MemoryRouter>
        </AppContext.Provider>
      );
    });

    await waitFor(() => jest.mock);

    const headerDividerElement = screen.getByText(/create/i);
    expect(headerDividerElement).toBeInTheDocument();
  })

  it('Admin user has access to Add user and Remove user buttons', async () => {
    act(() => {
      render(
        <AppContext.Provider value={{ user }}>
          <MemoryRouter
            initialEntries={["/groups"]}
            initialIndex={0}
          >
            <GroupsSection user={mockGroupsSectionData}  />
          </MemoryRouter>
        </AppContext.Provider>
      );
    });

    const buttons = screen.queryAllByRole('button');
    expect(buttons.length).toBe(2);
  })

  it('privileged user has access to Add user and Remove user buttons', async () => {
    act(() => {
      render(
        <AppContext.Provider value={{ user }}>
          <MemoryRouter
            initialEntries={["/groups"]}
            initialIndex={0}
          >
            <GroupsSection user={privilegedData}  />
          </MemoryRouter>
        </AppContext.Provider>
      );
    });

    const buttons = screen.queryAllByRole('button');
    expect(buttons.length).toBe(2);
  })

  it('Regular user has no access to Add user or Remove user buttons', async () => {
    act(() => {
      render(
        <AppContext.Provider value={{ user }}>
          <MemoryRouter
            initialEntries={["/groups"]}
            initialIndex={0}
          >
            <GroupsSection user={userData}  />
          </MemoryRouter>
        </AppContext.Provider>
      );
    });

    const buttons = screen.queryAllByRole('button');
    expect(buttons.length).toBe(0);
  })
});
