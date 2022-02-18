// vendor
import React from 'react';
import { act } from 'react-dom/test-utils';
import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
// components
import DatasetDetails from '../DatasetDetails';
import mockDatsetDetails from './DatasetDetailsData';

// mocks the environment interface
jest.mock('Environment/createEnvironment', () => {
  return {
    get: () => new Promise((resolve) => {
      resolve(
        {
          json: () => new Promise((resolve) => { resolve(mockDatsetDetails)})
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

  it('Renders details view', async () => {
    act(() => {
      render(
        <MemoryRouter
          initialEntries={['/default']}
          initialIndex={0}
        >
          <DatasetDetails
            dataseName="datasetName"
            isExpanded={true}
          />
        </MemoryRouter>
      );
    });

    await waitFor(() => jest.mock);
    const adminElement = screen.getByText(/admin,/i);
    const privilegedElement = screen.getByText(/privileged/i);
    const groupElement = screen.getByText(/group/i);


    expect(adminElement).toBeInTheDocument();
    expect(privilegedElement).toBeInTheDocument();
    expect(groupElement).toBeInTheDocument();
  });
});
