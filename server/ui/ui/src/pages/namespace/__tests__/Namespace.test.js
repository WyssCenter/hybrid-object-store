// vendor
import React from 'react';
import { act } from 'react-dom/test-utils';
import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
// components
import Namespace from '../Namespace';
// data
import mockNamespaceData  from './NamespaceData';

// mocks the environment interface
jest.mock('Environment/createEnvironment', () => {
  return {
    get: (route) => new Promise((resolve) => {
      const data = (route.indexOf('dataset') > -1)
        ? mockNamespaceData.datasets
        : mockNamespaceData.namespace;
      resolve(
        {
          json: () => new Promise((resolve) => { resolve(data)})
        }
      )
    })
  }
})

describe('Namespace', () => {
  beforeAll(() => {
    const user = JSON.stringify({ id_token: 'id_token' });
    localStorage.setItem(`oidc.user:http://localhost/auth/v1/.well-known/openid-configuration:HossServer`, user);
    window.location.pathname = '/default'
  })

  it('Renders namespace header', () => {
    act(() => {
      render(
        <MemoryRouter
          initialEntries={['/default']}
          initialIndex={0}
        >
          <Namespace />
        </MemoryRouter>
      );
    });
    const headerElement = screen.getByText(/default/i);
    expect(headerElement).toBeInTheDocument();
  })


  it('Loads data from api call', async () => {
    act(() => {
      render(
        <MemoryRouter
          initialEntries={['/default']}
          initialIndex={0}
        >
          <Namespace />
        </MemoryRouter>
      );
    });

    await waitFor(() => jest.mock);

    const datasetNameElement = screen.getByText(/users-dataset/i);

    expect(datasetNameElement).toBeInTheDocument();
  })
});
