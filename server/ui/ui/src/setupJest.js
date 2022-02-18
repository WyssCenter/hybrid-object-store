// jest-dom adds custom jest matchers for asserting on DOM nodes.
// allows you to do things like:
// expect(element).toHaveTextContent(/react/i)
// learn more: https://github.com/testing-library/jest-dom
import fetchMock, { enableFetchMocks } from 'jest-fetch-mock';
enableFetchMocks();





const user = JSON.stringify({ id_token: 'id_token' });
global.localStorage.setItem('oidc.user:http:/localhost/auth/v1/.well-known/openid-configuration:HossServer', user);
