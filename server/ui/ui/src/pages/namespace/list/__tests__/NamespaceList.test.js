// vendor
import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
// context
import AppContext from 'Src/AppContext';
// components
import NamespaceList from '../NamespaceList';
import mockDatsetDetails from './NamespaceListData';

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

describe('NamespaceList', () => {
  beforeAll(() => {
    let modalElement = window.document.createElement('div');
    modalElement.setAttribute('id', 'modal');

    document.body.appendChild(modalElement)
  })

  it('Renders all rows', () => {

    render(
      <AppContext.Provider value={{user: adminRole}}>
        <MemoryRouter
          initialEntries={['/default']}
          initialIndex={0}
        >
          <NamespaceList
            datasets={mockDatsetDetails}
          />
        </MemoryRouter>
      </AppContext.Provider>
    );

    const datasetName1Element = screen.getByText(/hanks-dataset/i);
    const datasetName2Element = screen.getByText(/my-dataset/i);
    const datasetName3Element = screen.getByText(/my-next-dataset/i);


    expect(datasetName1Element).toBeInTheDocument();
    expect(datasetName2Element).toBeInTheDocument();
    expect(datasetName3Element).toBeInTheDocument();

  });



  it('Search filters correctly', async () => {
    render(
      <AppContext.Provider value={{user: adminRole}}>
        <MemoryRouter
          initialEntries={['/default']}
          initialIndex={0}
        >
          <NamespaceList
            datasets={mockDatsetDetails}
          />
        </MemoryRouter>
      </AppContext.Provider>
    );


    const searchElement = screen.getByRole('textbox');

    await fireEvent.keyUp(searchElement, { target: { value: 'hanks' } });

    const datasetName1Element = screen.getByText(/hanks-dataset/i);


    expect(datasetName1Element).toBeInTheDocument();
    expect(searchElement.value).toBe('hanks');
    expect(screen.queryAllByRole('button').length).toBe(3);

  });


  it('Create Modal appears on click', async () => {
    render(
      <AppContext.Provider value={{user: adminRole}}>
        <MemoryRouter
          initialEntries={['/default']}
          initialIndex={0}
        >
          <NamespaceList
            datasets={mockDatsetDetails}
          />
        </MemoryRouter>
      </AppContext.Provider>
    );


    const buttons = screen.queryAllByRole('button');

    await fireEvent.click(buttons[0]);

    const createDatasetHeader = screen.getByText(/Create a new dataset here/i);


    expect(createDatasetHeader).toBeInTheDocument();

  });

  it('Admin user has access to delete buttons and create dataset button', async () => {
    render(
      <AppContext.Provider value={{user: adminRole}}>
        <MemoryRouter
          initialEntries={['/default']}
          initialIndex={0}
        >
          <NamespaceList
            datasets={mockDatsetDetails}
          />
        </MemoryRouter>
      </AppContext.Provider>
    );


    const buttons = screen.queryAllByRole('button');

    expect(buttons.length).toBe(7);

  });

  it('Privileged user has access to delete buttons and create dataset button', async () => {
    render(
      <AppContext.Provider value={{user: privilegedRole}}>
        <MemoryRouter
          initialEntries={['/default']}
          initialIndex={0}
        >
          <NamespaceList
            datasets={mockDatsetDetails}
          />
        </MemoryRouter>
      </AppContext.Provider>
    );


    const buttons = screen.queryAllByRole('button');

    expect(buttons.length).toBe(7);

  });

  it('Regular user does not have access to delete buttons and create dataset button', async () => {
    render(
      <AppContext.Provider value={{user: userRole}}>
        <MemoryRouter
          initialEntries={['/default']}
          initialIndex={0}
        >
          <NamespaceList
            datasets={mockDatsetDetails}
          />
        </MemoryRouter>
      </AppContext.Provider>
    );


    const buttons = screen.queryAllByRole('button');

    expect(buttons.length).toBe(3);

  });
});
