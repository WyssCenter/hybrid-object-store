// vendor
import React from 'react';
import { act } from 'react-dom/test-utils';
import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import moment from 'moment';
// context
import AppContext from 'Src/AppContext';
// components
import DatasetItem from '../DatasetItem';
import mockDatsetDetails from './DatasetItemData';


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


describe('Namespace', () => {
  beforeAll(() => {
    const user = JSON.stringify({ id_token: 'id_token' });
    localStorage.setItem(`oidc.user:http://localhost/auth/v1/.well-known/openid-configuration:HossServer`, user);
    window.location.pathname = '/default'
  })

  it('Renders details view', () => {
    render(
      <AppContext.Provider value={{user: adminRole}}>
        <MemoryRouter
          initialEntries={['/default']}
          initialIndex={0}
        >
          <DatasetItem
            dataset={mockDatsetDetails}
          />
        </MemoryRouter>
      </AppContext.Provider>
    );

    const datasetNameElement = screen.getByText(/my-dataset-name/i);
    const datasetDescription = screen.getByText(/This is my dataset/i);
    const createdElement = screen.getByText(moment(mockDatsetDetails.created).fromNow());
    const directoryElement = screen.getByText(/my-dataset-directory/i)


    expect(datasetNameElement).toBeInTheDocument();
    expect(datasetDescription).toBeInTheDocument();
    expect(createdElement).toBeInTheDocument();
    expect(directoryElement).toBeInTheDocument();
  });

  it('Admin user has access to delete button', async () => {
    act(() => {
      render(
        <AppContext.Provider value={{user: adminRole}}>
          <MemoryRouter
            initialEntries={['/default']}
            initialIndex={0}
          >
            <DatasetItem
              dataset={mockDatsetDetails}
            />
          </MemoryRouter>
        </AppContext.Provider>
      );
    })

    const buttons = screen.queryAllByRole('button');

    expect(buttons.length).toBe(2);

  });

  it('Privileged user has access to delete button', async () => {
    act(() => {
      render(
        <AppContext.Provider value={{user: privilegedRole}}>
          <MemoryRouter
            initialEntries={['/default']}
            initialIndex={0}
          >
            <DatasetItem
              dataset={mockDatsetDetails}
            />
          </MemoryRouter>
        </AppContext.Provider>
      );
    })

    const buttons = screen.queryAllByRole('button');

    expect(buttons.length).toBe(2);

  });

  it('Regular user does not have access to delete button', async () => {
    act(() => {
      render(
        <AppContext.Provider value={{user: userRole}}>
          <MemoryRouter
            initialEntries={['/default']}
            initialIndex={0}
          >
            <DatasetItem
              dataset={mockDatsetDetails}
            />
          </MemoryRouter>
        </AppContext.Provider>
      );
    })

    const buttons = screen.queryAllByRole('button');

    expect(buttons.length).toBe(1);

  });
});
