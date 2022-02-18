// vendor
import fetchMock, { enableFetchMocks } from 'jest-fetch-mock';
import { act } from 'react-dom/test-utils';
import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
// context
import DatasetContext from '../../../../../DatasetContext';
// components
import Delete from '../Delete';

jest.mock('Environment/createEnvironment', () => {
  return {
    del: () => new Promise((resolve) => {
      resolve(
        {
          ok: true,
          json: () => new Promise((resolve) => { resolve({ success: true})})
        }
      )
    })
  }
});


describe('Delete', () => {
  it('Renders', () => {
    render(
      <Delete
        datasetName="my-dataset"
        name="user"
        namespace="my-namespace"
        sectionType="user"
      />

    );
    const name = screen.getByRole('button');
    expect(name).toBeInTheDocument();
  });


  it('Tooltip appears after delete button click', () => {
    const send = jest.fn();
    render(
      <DatasetContext.Provider value={{ send }}>
        <Delete
          datasetName="my-dataset"
          name="user"
          namespace="my-namespace"
          sectionType="user"
        />
      </DatasetContext.Provider>
    );
    const button = screen.getByRole('button');
    fireEvent.click(button);

    expect(screen.getByText(/Are you sure/i)).toBeInTheDocument();
  });


  it('Cancel clears tooltip', async () => {
    const send = jest.fn();
    render(
      <DatasetContext.Provider value={{ send }}>
        <Delete
          datasetName="my-dataset"
          name="user"
          namespace="my-namespace"
          sectionType="user"
        />
      </DatasetContext.Provider>
    );
    const button = screen.getByRole('button');
    fireEvent.click(button);

    expect(screen.getByText(/Are you sure/i)).toBeInTheDocument();

    const buttons = screen.queryAllByRole('button');
    expect(buttons.length).toBe(3);

    fireEvent.click(buttons[1]);

    const buttonsAfterCancel = screen.queryAllByRole('button');

    expect(buttonsAfterCancel.length).toBe(1);
  });


  it('Submit Button Fires open tooltip and confirm fires delete, send is called to reset table', async () => {
    const send = jest.fn();
    act(() => {
        render(
          <DatasetContext.Provider value={{ send }}>
            <Delete
              datasetName="my-dataset"
              name="user"
              namespace="my-namespace"
              sectionType="user"
            />
          </DatasetContext.Provider>
        );
      }
    );
    const button = screen.getByRole('button');
    fireEvent.click(button);

    expect(screen.getByText(/Are you sure/i)).toBeInTheDocument();

    const confirmButton = screen.getByText(/confirm/i);
    fireEvent.click(confirmButton);

    await waitFor(() => jest.mock);

    expect(send).toBeCalledTimes(1);
  });


})
