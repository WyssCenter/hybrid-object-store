// vendor
import { render, screen } from '@testing-library/react';
import fetchMock, { enableFetchMocks } from 'jest-fetch-mock';
import authMachine from '../auth/machine/AuthStateMachine';
// components
import App from '../App';

describe('App', () => {
  beforeAll(() => {
    const user = JSON.stringify({ id_token: 'id_token' });
    localStorage.setItem(`oidc.user:http://localhost/auth/v1/.well-known/openid-configuration:HossServer`, user);
  })
  beforeEach(() => { // if you have an existing `beforeEach` just add the following line to it
    fetchMock.doMock()

    // TODO figure out tests are failing. Commented out for now.
    fetchMock.mockIf(/^http?:\/\/localhost*$/, req => {
        if (req.url.endsWith("namespace/")) {
          return  { json: () => [{
              name: 'namespace',
              description: 'this is a dataset',
              bucketName: 'data'
            }]
          }
        } else if (req.url.endsWith("/path2")) {
          return {
            body: {

            },
            headers: {
              "X-Some-Response-Header": "Some header value"
            }
          }
        } else {
          return {
            status: 404,
            body: "Not Found"
          }
        }
    })
  })

  it('Renders App', () => {
    render(
      <App
        machine={authMachine}
      />
    );
    const appElement = screen.getByText(/Loading/i);
    expect(appElement).toBeInTheDocument();
  });

})
