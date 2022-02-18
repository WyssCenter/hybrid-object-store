// vendor
import React from 'react';
import { render, screen, waitFor, fireEvent } from '@testing-library/react';
// components
import ListPat from '../ListPat';



// mocks the environment interface
jest.mock('Environment/createEnvironment', () => {
  return {

    del: () => new Promise((resolve) => {
      resolve(
        {
          text: () => new Promise((resolve) => {
            resolve({});
          })
        }
      )
    })
  }
})

const list = [
    {"id": 1, "description": 'Gigantum Client'},
    {"id": 2, "description": 'Gigantum Hub'},
    {"id": 3, "description": 'Gigantum Desktop'},
  ]

describe('ListPat', () => {
  beforeAll(() => {
    const user = JSON.stringify({ id_token: 'id_token' });
    localStorage.setItem(`oidc.user:http://localhost/auth/v1/.well-known/openid-configuration:HossServer`, user);
    let modalElement = window.document.createElement('div');
    modalElement.setAttribute('id', 'modal');

    document.body.appendChild(modalElement);
  })

  it('Renders all tokens fetched from api', async () => {
    render(
      <ListPat
        list={list}
        send={jest.fn()}
      />
    );
    await waitFor(() => jest.mock);
    const firstTableElement = screen.getByText(/Gigantum Client/);
    const secondTableElement = screen.getByText(/Gigantum Hub/);
    const thirdTableElement = screen.getByText(/Gigantum Desktop/);

    expect(firstTableElement).toBeInTheDocument();
    expect(secondTableElement).toBeInTheDocument();
    expect(thirdTableElement).toBeInTheDocument();
  });

  it('Clicking cancel fires handle close', async () => {
    render(
      <ListPat
        list={list}
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
    expect(newButtonElements[1]).not.toBeInTheDocument();
  });

})
