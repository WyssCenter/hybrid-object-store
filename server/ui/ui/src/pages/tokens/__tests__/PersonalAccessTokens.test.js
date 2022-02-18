
// vendor
import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
// components
import PersonalAccessTokens from '../PersonalAccessTokens';

const pat = {
  id: 2,
  description: 'Gigantum Client'
}


// mocks the environment interface
jest.mock('Environment/createEnvironment', () => {
  return {
    get: () => new Promise((resolve) => {
      resolve(
        {

          json: () => new Promise((resolve) => {

            resolve(
              [
                {"id": 1, "description": 'Gigantum Client'},
                {"id": 2, "description": 'Gigantum Hub'},
                {"id": 3, "description": 'Gigantum Desktop'},
              ]
            );
          })
        }
      )
    })
  }
});


describe('PersonalAccessTokens', () => {
  beforeAll(() => {
    jest.resetModules();
    const user = JSON.stringify({ id_token: 'id_token' });
    localStorage.setItem(`oidc.user:http://localhost/auth/v1/.well-known/openid-configuration:HossServer`, user);
    let modalElement = window.document.createElement('div');
    modalElement.setAttribute('id', 'modal');

    document.body.appendChild(modalElement)
  })


  it('Renders header', async () => {
    render(
      <PersonalAccessTokens />
    );
    await waitFor(() => jest.mock);
    const labelDescription = screen.getAllByText(/Personal Access Tokens/i);

    expect(labelDescription[0]).toBeInTheDocument();
  });

  it('Renders all tokens fetched from api', async () => {
    render(
      <PersonalAccessTokens />
    );
    await waitFor(() => jest.mock);
    const firstTableElement = screen.getByText(/Gigantum Client/);
    const secondTableElement = screen.getByText(/Gigantum Hub/);
    const thirdTableElement = screen.getByText(/Gigantum Desktop/);

    expect(firstTableElement).toBeInTheDocument();
    expect(secondTableElement).toBeInTheDocument();
    expect(thirdTableElement).toBeInTheDocument();
  });



})
