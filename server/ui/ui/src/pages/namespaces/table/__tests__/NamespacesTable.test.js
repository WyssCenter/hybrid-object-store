// vendor
import React from 'react';
import { act } from 'react-dom/test-utils';
import { render, screen, fireEvent } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
// context
import AppContext from 'Src/AppContext';
// components
import NamespacesTable from '../NamespacesTable';
import mockNamepacesData from './NamespacesTableData';


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


describe('NamespacesTable', () => {
  beforeAll(() => {
    let modalElement = window.document.createElement('div');
    modalElement.setAttribute('id', 'modal');

    document.body.appendChild(modalElement)
  })

  it('Renders all rows', () => {

    render(
      <AppContext.Provider value={{user: adminRole}}>
        <MemoryRouter
          initialEntries={['/']}
          initialIndex={0}
        >
          <NamespacesTable
            namespaces={mockNamepacesData}
          />
        </MemoryRouter>
      </AppContext.Provider>
    );

    const namespaceName1Element = screen.getByText(/default-namespace/i);
    const namespaceName2Element = screen.getByText(/my-namespace/i);
    const namespaceName3Element = screen.getByText(/my-next-namespace/i);
    const namespaceName4Element = screen.getByText(/new-namespace/i);

    expect(namespaceName1Element).toBeInTheDocument();
    expect(namespaceName2Element).toBeInTheDocument();
    expect(namespaceName3Element).toBeInTheDocument();
    expect(namespaceName4Element).toBeInTheDocument();

  });



  it('Search filters correctly', async () => {
    render(
      <AppContext.Provider value={{user: adminRole}}>
        <MemoryRouter
          initialEntries={['/default']}
          initialIndex={0}
        >
          <NamespacesTable
            namespaces={mockNamepacesData}
          />
        </MemoryRouter>
      </AppContext.Provider>
    );


    const searchElement = screen.getByRole('textbox');

    await fireEvent.keyUp(searchElement, { target: { value: 'my-next' } });

    const namepsaceElement = screen.getByText(/my-next-namespace/i);


    expect(namepsaceElement).toBeInTheDocument();
    expect(searchElement.value).toBe('my-next');
    expect(screen.queryAllByRole('button').length).toBe(3);

  });


  it('Create Modal appears on click', async () => {
    render(
      <AppContext.Provider value={{user: adminRole}}>
        <MemoryRouter
          initialEntries={['/default']}
          initialIndex={0}
        >
          <NamespacesTable
            namespaces={mockNamepacesData}
          />
        </MemoryRouter>
      </AppContext.Provider>
    );


    const buttons = screen.queryAllByRole('button');

    await fireEvent.click(buttons[0]);

    const createDatasetHeader = screen.getByText(/Create a new namespace here/i);


    expect(createDatasetHeader).toBeInTheDocument();

  });


  it('Admin user has access to delete buttons and create namespace button', async () => {
    render(
      <AppContext.Provider value={{user: adminRole}}>
        <MemoryRouter
          initialEntries={['/default']}
          initialIndex={0}
        >
          <NamespacesTable
            namespaces={mockNamepacesData}
          />
        </MemoryRouter>
      </AppContext.Provider>
    );


    const buttons = screen.queryAllByRole('button');

    expect(buttons.length).toBe(9);

  });


  it('Privileged user does not have access to delete buttons', async () => {
    render(
      <AppContext.Provider value={{user: privilegedRole}}>
        <MemoryRouter
          initialEntries={['/default']}
          initialIndex={0}
        >
          <NamespacesTable
            namespaces={mockNamepacesData}
          />
        </MemoryRouter>
      </AppContext.Provider>
    );


    const buttons = screen.queryAllByRole('button');

    expect(buttons.length).toBe(4);

  });



  it('Regular user does not have access to delete buttons', async () => {
    render(
      <AppContext.Provider value={{user: userRole}}>
        <MemoryRouter
          initialEntries={['/default']}
          initialIndex={0}
        >
          <NamespacesTable
            namespaces={mockNamepacesData}
          />
        </MemoryRouter>
      </AppContext.Provider>
    );


    const buttons = screen.queryAllByRole('button');

    expect(buttons.length).toBe(4);

  });
});
