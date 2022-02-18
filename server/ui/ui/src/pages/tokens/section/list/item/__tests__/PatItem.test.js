// vendor
import React from 'react';
import { render, screen, waitFor, fireEvent } from '@testing-library/react';
// components
import PatItem from '../PatItem';

const pat = {
  id: 2,
  description: 'Gigantum Client'
}


// mocks the environment interface
jest.mock('Environment/createEnvironment', () => {
  return {
    del: () => new Promise((resolve) => {
      resolve(
        {
          text: () => new Promise((resolve) => {
            resolve();
          })
        }
      )
    })
  }
})

describe('PatItem', () => {
  beforeAll(() => {
    jest.resetModules();
    const user = JSON.stringify({ id_token: 'id_token' });
    localStorage.setItem(`oidc.user:http://localhost/auth/v1/.well-known/openid-configuration:HossServer`, user);
    let modalElement = window.document.createElement('div');
    modalElement.setAttribute('id', 'modal');

    document.body.appendChild(modalElement)
  })

  it('Renders pat item', async () => {
    render(
      <PatItem
        pat={pat}
        send={jest.fn()}
     />
    );
    const firstTableElement = screen.getByText(/Gigantum Client/);

    expect(firstTableElement).toBeInTheDocument();
  });

  it('Clicking cancel doesn\'t delete it', async () => {
    render(
      <PatItem
        pat={pat}
        send={jest.fn()}
      />
    );
    await waitFor(() => jest.mock);
    const buttonElements = screen.queryAllByRole('button');
    fireEvent.click(buttonElements[0]);

    const newButtonElements = screen.queryAllByRole('button');
    fireEvent.click(newButtonElements[1]);
    const firstTableElement = screen.getByText(/Gigantum Client/);
    expect(firstTableElement).toBeInTheDocument();
  });


  it('Deleting fires reset', async () => {
    const reset = jest.fn();
    render(
      <PatItem
        pat={pat}
        send={reset}
      />
    );
    const buttonElements = screen.queryAllByRole('button');
    await fireEvent.click(buttonElements[0]);
    const newButtonElements = screen.queryAllByRole('button');
    fireEvent.click(newButtonElements[2]);
    await waitFor(() => jest.mock);


    expect(reset).toHaveBeenCalledTimes(1);
  });

})
