// vendor
import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
// components
import CreatePatModal from '../CreatePatModal';


// mocks the environment interface
jest.mock('Environment/createEnvironment', () => {
  return {
    post: () => new Promise((resolve) => {
      resolve(
        {
          json: () => new Promise((resolve) => {
            resolve({
              id: 2,
              token: 'tH1s1S4T0K3n',
            });
          })
        }
      )
    })
  }
})

describe('CreatePatModal', () => {
  beforeAll(() => {
    let modalElement = window.document.createElement('div');
    modalElement.setAttribute('id', 'modal');

    document.body.appendChild(modalElement)
  })

  it('Renders create pat', () => {
    render(
      <CreatePatModal />
    );
    const createModalElement = screen.getByText(/What Is This Token For/i);
    expect(createModalElement).toBeInTheDocument();
  });

  it('Input updates description in input[text]', () => {
    render(
      <CreatePatModal />
    );
    const inputElement = screen.getByRole(/textbox/)
    fireEvent.keyUp(inputElement, { target: { value: 'gigantum token'}});
    expect(inputElement.value).toBe('gigantum token');
  });


  it('NewPat appears after submitting a request for a new token', async () => {
    render(
      <CreatePatModal />
    );
    const inputElement = screen.getByRole(/textbox/);
    await fireEvent.keyUp(inputElement, { target: { value: 'gigantum token'}});

    const buttonElement = screen.getByRole(/button/);
    await fireEvent.click(buttonElement);
    await waitFor(() => jest.mock)

    const newInputElement = screen.getByRole(/textbox/);

    expect(newInputElement.value).toBe('tH1s1S4T0K3n');
  });

})
