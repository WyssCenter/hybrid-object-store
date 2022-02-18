// vendor
import React from 'react';
import { act } from 'react-dom/test-utils';
import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
// conte
import AppContext from '../../../../AppContext';
// data
import mockGroupsData from './GroupsData';
// components
import Groups from '../Groups';


jest.mock('Environment/createEnvironment', () => {
  return {
    get: () => new Promise((resolve) => {
      resolve(
        {
          json: () => new Promise((resolve) => { resolve(mockGroupsData)})
        }
      )
    })
  }
});

const user = { profile: { name: 'admin' }};

describe('Groups', () => {

  it('Renders Groups', async () => {
    act(() => {
      render(
        <AppContext.Provider value={{ user }}>
          <MemoryRouter
            initialEntries={["/groups"]}
            initialIndex={0}
          >
            <Groups />
          </MemoryRouter>
        </AppContext.Provider>
      );
    });

    const headerElement = screen.getByText(/Groups/i);
    const loadingElement = screen.getByText(/Fetching/i);
    expect(headerElement).toBeInTheDocument();
    expect(loadingElement).toBeInTheDocument();
  });


  it('Renders groups table', async () => {
    act(() => {
      render(
        <AppContext.Provider value={{ user }}>
          <MemoryRouter
            initialEntries={["/groups"]}
            initialIndex={0}
          >
            <Groups />
          </MemoryRouter>
        </AppContext.Provider>
      );
    });

    await waitFor(() => jest.mock);

    const listElement = screen.getByText(/group-1/i);
    expect(listElement).toBeInTheDocument();
  })
});
