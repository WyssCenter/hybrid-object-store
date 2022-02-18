// vendor

import React from 'react';
import { act } from 'react-dom/test-utils';
import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
// context
import AppContext from 'Src/AppContext';
import NamespaceListingContext from '../NamespaceListingContext'
// components
import NamespaceListing from '../NamespaceListing';
// data
import mockNamespaceListingData from './NamepsaceListingData';


// mocks the environment interface
jest.mock('Environment/createEnvironment', () => {
  return {
    get: () => new Promise((resolve) => {
      resolve(
        {
          json: () => new Promise((resolve) => { resolve(mockNamespaceListingData)})
        }
      )
    })
  }
});


const adminRole = {
  profile: {
    name: 'admin',
    role: 'admin'
  }
}

describe('NamespaceListing', () => {

  test('Renders api data', async () => {
    act(() => {
      render(
        <AppContext.Provider value={{user: adminRole}}>
          <NamespaceListingContext.Provider value={{send: jest.fn()}}>
            <MemoryRouter
              initialEntries={['/default']}
              initialIndex={0}
            >
              <NamespaceListing />
            </MemoryRouter>
          </NamespaceListingContext.Provider>
        </AppContext.Provider>
      );
    });

    await waitFor(() => jest.mock);

    const namespaceNameHeader = screen.getByText(/default-namespace/i);
    expect(namespaceNameHeader).toBeInTheDocument();
  });
});
