// import
import React from 'react';
import fetchMock from 'jest-fetch-mock';
import { Machine } from 'xstate';
import { createModel } from '@xstate/test';
import { render, fireEvent, cleanup, act, waitFor} from '@testing-library/react';
// components
import App from '../../../App';

import { configuration } from '../AuthStateMachine';




// TODO figure out tests are failing. Commented out for now.
configuration.states.loading.invoke.src = waitFor(async () => {
  await new Promise((resolve) => {
    resolve({
      wellKnown: {
        authority: 'http:/localhost/auth/v1/.well-known/openid-configuration',
      },
      user: {
        authority: 'http:/localhost/auth/v1/.well-known/openid-configuration',
      }
    });
  });
});
//
//
configuration.states.checkLogin.invoke.src = waitFor(async (context) =>
  await new Promise((resolve) => {
    if (window.location.hash.indexOf('id_token') > -1) {
      resolve('COMPLETE_AUTHENTICATION');
    } else if (user) {
      resolve('USER_LOGGED_IN');
    } else if (user ===  null) {
      resolve('LOGGED_OUT');
    }
  }))

  if (window.location.hash.indexOf('id_token') > -1) {
    await callback('COMPLETE_AUTHENTICATION');
  } else if (user) {
    await callback('USER_LOGGED_IN');
  } else if (user ===  null) {
    await callback('LOGGED_OUT');
  }
  return undefined;
});

configuration.states.idle.meta = {
  test: async (page) => {
    page.debug()
    await page.getByText('Loading');
  }
};

configuration.states.loading.meta = {
  test: async (page) => {
    await page.getByText('Loading');
  }
};

configuration.states.checkLogin.meta = {
  test: async (page) => {
    await page.getByText('Error');
  }
};

configuration.states.loggedIn.meta = {
  test: async (page) => {
    await page.getByText('Server');
  }
};


configuration.states.logout.meta = {
  test: async (page) => {
    await page.getByText('Username');
  }
};

configuration.states.loggedOut.meta = {
  test: async (page) => {
    await page.getByText('Login');
  }
};

configuration.states.authorizeAuthenticate.meta = {
  test: async (page) => {
    await page.getByText('Loading');
  }
};


configuration.states.completeAuthentication.meta = {
  test: async (page) => {
    await page.getByText('Loading');
  }
};


configuration.states.error.meta = {
  test: async (page) => {
    await page.getByText('Error');
  }
};

configuration.states.systemError.meta = {
  test: async (page) => {
    await page.getByText('Error');
  }
};

configuration.states.redirect.meta = {
  test: async (page) => {
    await page.getByText('Loading');
  }
};

const AuthMachineMock = Machine(configuration);

const user = JSON.stringify({ id_token: 'id_token' });

global.localStorage.removeItem('oidc.user:http:/localhost/auth/v1/.well-known/openid-configuration:HossServer', user);


describe('AuthMachine model testing', () => {
  beforeEach(() => {
    fetchMock.mockIf(/^https?:\/\/localhost.*$/, async req => {
      if (req.url.endsWith("/openid-configuration")) {
          return {
            "issuer": "http://localhost/auth/v1",
            "authorization_endpoint": "http://localhost/auth/v1/authorize",
            "token_endpoint": "http://localhost/auth/v1/token",
            "jwks_uri": "http://localhost/auth/v1/keys",
            "userinfo_endpoint": "http://localhost/auth/v1/userinfo",
            "grant_types_supported": [
              "authorization_code",
              "implicit"
            ],
            "response_types_supported": [
              "code",
              "id_token"
            ],
            "subject_types_supported": [
              "public"
            ],
            "id_token_signing_alg_values_supported": [
              "RS256"
            ],
            "scopes_supported": [
              "openid",
              "email",
              "profile",
              "hos"
            ],
            "token_endpoint_auth_methods_supported": [
              "client_secret_post"
            ],
            "claims_supported": [
              "iss",
              "sub",
              "aud",
              "iat",
              "exp",

              "email",
              "email_verified",
              "name",
              "given_name",
              "family_name",
              "nickname"
            ]
          }
        } else if (req.url.endsWith("/path2")) {
          return {
            body: "another response body",
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
  });

  const authModel = createModel(AuthMachineMock, {
    events: {
      FETCH_CONFIG: {
        exec: async (page) => {
          await page.getByText('Loading')
        },
        // cases:['loading']
      },
      AUTH: {
        exec: async ({ waitForNextUpdate }) => {
          await waitForNextUpdate();
        },
        cases:['authorizeAuthenticate', 'error']
      },
      COMPLETE_AUTHENTICATION: {
        exec: async ({ getByText }) => {
          global.localStorage.setItem('oidc.user:http:/localhost/auth/v1/.well-known/openid-configuration:HossServer', user);
        },
        cases:['loggedIn', 'error']
      },
      LOGIN: {
        exec: async ({ getByText }) => {

          global.localStorage.setItem('oidc.user:http:/localhost/auth/v1/.well-known/openid-configuration:HossServer', user);
        },
        cases:['loggedOut', 'error']
      },
      LOGGED_OUT: {
        exec: async ({ getByText }) => {

          global.localStorage.removeItem('oidc.user:http:/localhost/auth/v1/.well-known/openid-configuration:HossServer');

          getByText('Error')
        }
      },
      USER_LOGGED_IN: {
        exec: async ({ getByText }) => {
          global.localStorage.setItem('oidc.user:http:/localhost/auth/v1/.well-known/openid-configuration:HossServer', user);
          getByText('Server')
        }
      },
      ERROR: {
        exec: async ({ getByText }) => {
          global.localStorage.setItem('oidc.user:http:/localhost/auth/v1/.well-known/openid-configuration:HossServer', user);
          getByText('Error')
        }
      },
      LOGOUT:  {
        exec: async ({ getByText }) => {
          global.localStorage.setItem('oidc.user:http:/localhost/auth/v1/.well-known/openid-configuration:HossServer', user);
        }
      },
    }
  });


  const testPlans = authModel.getShortestPathPlans();
  //
  testPlans.forEach((plan) => {
    describe(plan.description, () => {
      afterEach(cleanup);
      plan.paths.forEach((path) => {
        it(path.description, async () => {
          // do any setup, then...
          const rendered = render(<App machine={AuthMachineMock} />);
          await path.test(rendered);
        });
      });
    });
  });

  it('should have full coverage', () => {
    return authModel.testCoverage();
  });
});
