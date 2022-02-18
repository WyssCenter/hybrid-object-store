// vendor
import React from 'react';
import { act } from 'react-dom/test-utils';
import { render, screen, fireEvent } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import moment from 'moment';
// context
import AppContext from 'Src/AppContext';
// components
import NamespaceRow from '../NamespaceRow';
import mockNamespaceDetails from './NamespaceRowData';


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

// mocks the environment interface
jest.mock('Environment/createEnvironment', () => {
  return {
    del: (route) => new Promise((resolve) => {
      resolve(
        {
          json: () => new Promise((resolve) => { resolve({ success: true }) })
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

  it('Renders details view', () => {
    render(
      <AppContext.Provider value={{user: adminRole}}>
        <MemoryRouter
          initialEntries={['/default']}
          initialIndex={0}
        >
          <NamespaceRow
            namespace={mockNamespaceDetails}
          />
        </MemoryRouter>
      </AppContext.Provider>
    );

    const namespaceNameCell = screen.getByText(/default-namespace/i);
    const datasetDescriptionCell = screen.getByText(/Default namespace/i);
    const typeCell = screen.getByText(/minio/i);
    const bucketCell = screen.getByText(/data/i)


    expect(namespaceNameCell).toBeInTheDocument();
    expect(datasetDescriptionCell).toBeInTheDocument();
    expect(typeCell).toBeInTheDocument();
    expect(bucketCell).toBeInTheDocument();
  });


  it('Clicking delete opens tooltip and cancel closes', async () => {
    act(() => {
      render(
        <AppContext.Provider value={{user: adminRole}}>
          <MemoryRouter
            initialEntries={['/default']}
            initialIndex={0}
          >
            <NamespaceRow
              namespace={mockNamespaceDetails}
            />
          </MemoryRouter>
        </AppContext.Provider>
      );
    })

    const deleteButton = screen.queryAllByRole('button')[1];

    fireEvent.click(deleteButton);

    expect(screen.getByText(/are you sure/i)).toBeInTheDocument();

    const cancelButton = screen.getByText('Cancel');

    fireEvent.click(cancelButton);


    expect(screen.queryAllByText(/are you sure/i).length).toBe(0);

  });


  it('Clicking delete opens tooltip, clicking confirm fires delete and closes tooltip', async () => {
    act(() => {
      render(
        <AppContext.Provider value={{user: adminRole}}>
          <MemoryRouter
            initialEntries={['/default']}
            initialIndex={0}
          >
            <NamespaceRow
              namespace={mockNamespaceDetails}
            />
          </MemoryRouter>
        </AppContext.Provider>
      );
    })

    const deleteButton = screen.queryAllByRole('button')[1];

    fireEvent.click(deleteButton);

    expect(screen.getByText(/are you sure/i)).toBeInTheDocument();

    const buttons = screen.queryAllByRole('button');

    fireEvent.click(buttons[2]);

    expect(screen.queryAllByText(/are you sure/i).length).toBe(0);

  });



  it('Admin user has access to delete button', async () => {
    act(() => {
      render(
        <AppContext.Provider value={{user: adminRole}}>
          <MemoryRouter
            initialEntries={['/default']}
            initialIndex={0}
          >
            <NamespaceRow
              namespace={mockNamespaceDetails}
            />
          </MemoryRouter>
        </AppContext.Provider>
      );
    })

    const buttons = screen.queryAllByRole('button');

    expect(buttons.length).toBe(2);

  });


  it('Privileged user does not have access to delete button', async () => {
    act(() => {
      render(
        <AppContext.Provider value={{user: privilegedRole}}>
          <MemoryRouter
            initialEntries={['/default']}
            initialIndex={0}
          >
            <NamespaceRow
              namespace={mockNamespaceDetails}
            />
          </MemoryRouter>
        </AppContext.Provider>
      );
    })

    const deleteButton = screen.queryAllByRole('button');

    expect(deleteButton.length).toBe(1);

  });


  it('Regular user does not have access to delete button', async () => {
    act(() => {
      render(
        <AppContext.Provider value={{user: userRole}}>
          <MemoryRouter
            initialEntries={['/default']}
            initialIndex={0}
          >
            <NamespaceRow
              namespace={mockNamespaceDetails}
            />
          </MemoryRouter>
        </AppContext.Provider>
      );
    })

    const deleteButton = screen.queryAllByRole('button');

    expect(deleteButton.length).toBe(1);

  });
});
