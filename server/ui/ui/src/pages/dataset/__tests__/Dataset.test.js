// vendor
import React from 'react';
import { act } from 'react-dom/test-utils';
import { render, screen, waitFor } from '@testing-library/react';
// components
import Dataset from '../Dataset';
// data
import mockDatasetData  from './DatasetData';

// mocks the environment interface
jest.mock('Environment/createEnvironment', () => {
  return {
    get: () => new Promise((resolve) => {
      resolve(
        {
          json: () => new Promise((resolve) => { resolve(mockDatasetData)})
        }
      )
    })
  }
})


describe('Dataset', () => {

  it('Renders loading', () => {
    act(() => {
      render(
        <Dataset />
      );
    });
    const fetchingText = screen.getByText(/Fetching/i);
    expect(fetchingText).toBeInTheDocument();
  })

  it('Renders Name and Description', async () => {
     act(() => {
      render(
        <Dataset />
      );
    });

    await waitFor(() => jest.mock)

    const nameText = screen.getByText(/my/i);
    const descriptionText = screen.getByText(/testing this out/i);

    expect(nameText).toBeInTheDocument();
    expect(descriptionText).toBeInTheDocument();
  });

})
