// vendor
import React from 'react';
import { act } from 'react-dom/test-utils';
import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
// conte
import AppContext from '../../../../AppContext';
// data
import mockGroupData from './GroupData';
// components
import Group from '../Group';


jest.mock('Environment/createEnvironment', () => {
  return {
    get: () => new Promise((resolve) => {
      resolve(
        {
          json: () => new Promise((resolve) => { resolve(mockGroupData)})
        }
      )
    })
  }
});

const user = { profile: { name: 'admin' }};

describe('Group', () => {

  it('Renders Group', async () => {
    act(() => {
      render(
        <AppContext.Provider value={{ user }}>
          <MemoryRouter
            initialEntries={["/groups"]}
            initialIndex={0}
          >
            <Group />
          </MemoryRouter>
        </AppContext.Provider>
      );
    });

    const headerElement = screen.queryAllByText(/Group/i);
    const loadingElement = screen.getByText(/Fetching/i);
    expect(headerElement[0]).toBeInTheDocument();
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
            <Group />
          </MemoryRouter>
        </AppContext.Provider>
      );
    });

    await waitFor(() => jest.mock);

    const listElement = screen.queryAllByText(/admin/i);
    expect(listElement.length).toBe(3);
  })
});
