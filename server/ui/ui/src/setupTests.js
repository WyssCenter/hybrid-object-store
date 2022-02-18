// jest-dom adds custom jest matchers for asserting on DOM nodes.
// allows you to do things like:
// expect(element).toHaveTextContent(/react/i)
// learn more: https://github.com/testing-library/jest-dom
import '@testing-library/jest-dom';
import { JSDOM } from "jsdom";
import 'jest-fetch-mock';

const indexHTML =
  `<html lang="en">
    <body>
      <div id="modal__cover" class="modal__cover hidden"></div>
      <div id="modal" class="ReactDom"></div>
      <div id="header" class="ReactDom"></div>
      <div id="loader" class="Loader fixed--important hidden"></div>
    </body>
  </html>`;
const dom = new JSDOM(indexHTML);



jest.mock('react-router-dom', () => ({
    ...jest.requireActual('react-router-dom'),
    useLocation: () => ({
      pathname: '/namespace',
      search: '',
      hash: '',
      state: null,
      key: '5nvxpbdafa',
    }),

    useParams: () => ({
      namespace: 'default',
      datasetName: 'datasetName',
      groupname: 'groupname'
    }),
}));


const user = JSON.stringify({ id_token: 'id_token' });


global.document = dom.window.document;
global.window = dom.window;
global.localStorage.setItem('oidc.user:http:/localhost/auth/v1/.well-known/openid-configuration:HossServer', user);
global.navigator = { userAgent: "jest" };
