// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
// context
import AppContext from 'Src/AppContext';
// components
import DatasetSection from '../DatasetSection';
// data
import dataset from './DatasetSectionData';

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

describe('DatasetSection', () => {
  it('Renders', () => {
    render(
      <AppContext.Provider value={{user: adminRole}}>
        <DatasetSection
          dataset={dataset}
          datasetName="my-dataset"
          namespace="my-namespace"

        />
      </AppContext.Provider>

    );
    const linkElement = screen.getByText(/my-group/i);
    expect(linkElement).toBeInTheDocument();
  });


  it('Has group name admin', () => {
    render(
      <AppContext.Provider value={{user: adminRole}}>

        <DatasetSection
          dataset={dataset}
          datasetName="my-dataset"
          namespace="my-namespace"

        />
      </AppContext.Provider>

    );
    const linkElement = screen.getByText(/admin/i);
    expect(linkElement).toBeInTheDocument();
  });

  it('Admin user has access to actionable buttons', async () => {
    render(
      <AppContext.Provider value={{user: adminRole}}>

        <DatasetSection
          dataset={dataset}
          datasetName="my-dataset"
          namespace="my-namespace"

        />
      </AppContext.Provider>
    );


    const buttons = screen.queryAllByRole('button');

    expect(buttons.length).toBe(7);

  });

  it('Privileged user has access to actionable buttons', async () => {
    render(
      <AppContext.Provider value={{user: privilegedRole}}>

        <DatasetSection
          dataset={dataset}
          datasetName="my-dataset"
          namespace="my-namespace"

        />
      </AppContext.Provider>
    );


    const buttons = screen.queryAllByRole('button');

    expect(buttons.length).toBe(7);

  });

  it('Regular user does not have access to actionable buttons', async () => {
    render(
      <AppContext.Provider value={{user: userRole}}>

        <DatasetSection
          dataset={dataset}
          datasetName="my-dataset"
          namespace="my-namespace"

        />
      </AppContext.Provider>
    );


    const buttons = screen.queryAllByRole('button');

    expect(buttons.length).toBe(0);

  });

});
