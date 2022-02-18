// vendor
import React from 'react';
import { act } from 'react-dom/test-utils';
import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
// components
import NamespaceDetails from '../NamespaceDetails';
import mockDatsetDetails from './NamespaceDetailsData';

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
          <NamespaceDetails
            dataseName="datasetName"
            isExpanded={true}
          />
        </MemoryRouter>
      );
    });

    await waitFor(() => jest.mock);
    const datasetName1 = screen.getByText(/hanks-dataset,/i);
    const datasetName2 = screen.getByText(/my-datase/i);
    const datasetName3 = screen.getByText(/my-next-dataset/i);


    expect(datasetName1).toBeInTheDocument();
    expect(datasetName2).toBeInTheDocument();
    expect(datasetName3).toBeInTheDocument();
  });
});
