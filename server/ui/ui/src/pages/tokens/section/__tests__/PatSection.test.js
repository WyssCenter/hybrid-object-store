
// vendor
import React from 'react';
import { render, screen, waitFor, fireEvent } from '@testing-library/react';
// components
import PatSection from '../PatSection';

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


describe('PatSection', () => {
  beforeAll(() => {
    jest.resetModules();
    const user = JSON.stringify({ id_token: 'id_token' });
    localStorage.setItem(`oidc.user:http://localhost/auth/v1/.well-known/openid-configuration:HossServer`, user);
    let modalElement = window.document.createElement('div');
    modalElement.setAttribute('id', 'modal');

    document.body.appendChild(modalElement)
  })

  it('Renders all tokens fetched from api', async () => {
    render(
      <PatSection />
    );
    await waitFor(() => jest.mock);
    const firstTableElement = screen.getByText(/Gigantum Client/);
    const secondTableElement = screen.getByText(/Gigantum Hub/);
    const thirdTableElement = screen.getByText(/Gigantum Desktop/);

    expect(firstTableElement).toBeInTheDocument();
    expect(secondTableElement).toBeInTheDocument();
    expect(thirdTableElement).toBeInTheDocument();
  });


  it('Renders create pat', async () => {
    render(
      <PatSection />
    );
    await waitFor(() => jest.mock);
    const labelDescription = screen.getByText(/What is this token for/i);

    expect(labelDescription).toBeInTheDocument();
  });


})
