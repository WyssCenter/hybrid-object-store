// vendor
import { UserManager, WebStorageStateStore } from 'oidc-client';
// api
import { getWellKnown } from 'Environment/createEnvironment';
// types
import  {
  Config,
  User,
  AuthMachineContext,
} from './AuthTypes';

const id = 'tokenExchange';

const baseHeaders = {
  "Access-Control-Allow-Origin": "*",
  "Content-Type": 'application/json',
  "Origin": window.location.origin
}

const baseUrl = `${window.location.protocol}//${window.location.hostname}/core/v1`;
let authUrl = `${window.location.protocol}//${window.location.hostname}/auth/v1`;

fetch(`${baseUrl}/discover`, { headers: baseHeaders, method: 'GET' })
.then(response => response.json())
.then(data => authUrl = data.auth_service);

/**
* Method is passed a route to get the wellKnown service config.
* @param {string} route
* @return {Promise<Config>}
* @fires {#getWellKnown}
*/
const fetchConfig: (route: string) => Promise<Config> = route => {
  return new Promise((resolve, reject) => {
    getWellKnown(route)
      .then((response) => {
        return (response.json())
      }).then(data => {
        const userAuthManager = new UserManager({
          authority: `${authUrl}/.well-known/openid-configuration`,
          client_id: 'HossServer',
          redirect_uri: `${window.location.origin}/ui/`,
          post_logout_redirect_uri: `${window.location.origin}/ui/`,
          response_mode: 'fragment',
          scope: data.scopes_supported.join(' '),
          response_type: 'id_token',
          filterProtocolClaims: false,
          metadata: {
            ...data,
          },
          userStore: new WebStorageStateStore({
            store: window.localStorage
          }),
          stateStore: new WebStorageStateStore({
           store: window.localStorage
         })
        });

        resolve({ wellKnown: data, userAuthManager })
      })
      .catch((error) => {
        reject(error);
      });


  });
};

/**
* Method is passed contexxt and signs in if id_token is not visible
* @param {AuthMachineContext} context
* @return {Promise<any>}
* @fires {userAuthManager#signinRedirect}
*/
const authorizeAuthenticateAttempt: (context: AuthMachineContext) => Promise<any> = (context) => {
  const { userAuthManager } = context.config;

  if (window.location.hash.indexOf('id_token') > -1) {
    return new Promise((resolve, reject) => {
      resolve({ state: 'tokenExchange'});
    })
  }

  return new Promise((resolve, reject) => {
      userAuthManager.signinRedirect({ state: { bar: id }}).then(() => {
        resolve({ state: { bar : id }});
      }).catch(error => {
        reject(error);
      });

  });
}

/**
* Method is passed contexxt and completes authentication
* @param {AuthMachineContext} context
* @return {Promise<any>}
* @fires {userAuthManager#signinRedirectCallback}
*/
const completeAuthentication: (context: AuthMachineContext) => Promise<any> = (context) => {
  const { userAuthManager } = context.config

  return new Promise((resolve, reject) => {

      userAuthManager.signinRedirectCallback().then((user) => {
        window.location.hash = '';
        if (localStorage.getItem('redirect')) {
          const parsedObject = JSON.parse(localStorage.getItem('redirect'));
          if (parsedObject && user && user.profile) {
            if (parsedObject[user.profile.email]) {
              window.location.pathname = parsedObject[user.profile.email];
              const newParsedObject = Object.assign({}, parsedObject, {[user.profile.email]: null})
              localStorage.setItem('redirect', JSON.stringify(newParsedObject));
            }
          }
        }
        context.user = user;
        resolve(user);
      }).catch(error => {
        reject(error);
      });

  });
}

/**
* Method checks if the user is authenticated
* @param {Function} callback
* @param {Function} onRecieve
* @param {AuthMachineContext} context
* @return {string}
* @fires {userAuthManager#getUser}

Note: there are issues with typescript and xstate. For now these types will be set as any until xstate fully supports typescript.
*/
const checkUserAuthenticated = (callback: any, onRecieve: any, context: AuthMachineContext):any => {
  context.config.userAuthManager.getUser().then((user: User) => {
    const timeNow = (new Date().getTime()/1000)
    const expired = user && user.profile  && (user.profile.exp - timeNow) < 0


    if (window.location.hash.indexOf('id_token') > -1) {
      callback('COMPLETE_AUTHENTICATION');
      context.user = user;
    } else if (user && user.profile && expired) {
      callback('REMOVE_USER')
    } else if (user) {
      context.user = user;
      callback('USER_LOGGED_IN');
    } else if (user ===  null) {
      callback('LOGGED_OUT');
    }
  });

  return 'done';
}

/**
* Method removes user from localStorage
* @param {AuthMachineContext} context
* @return {Promise<any>}
* @fires {userAuthManager#removeUser}
*/
const logout: (context: AuthMachineContext) => Promise<any> = (context) => {
  const { userAuthManager } = context.config
  return new Promise((resolve, reject) => {

      userAuthManager.removeUser().then((data) => {
        resolve({ user: null });
      }).catch(error => {
        reject(error);
      });

  });
}

export {
  authorizeAuthenticateAttempt,
  completeAuthentication,
  checkUserAuthenticated,
  fetchConfig,
  logout
}


const machineUtils = {
  authorizeAuthenticateAttempt,
  completeAuthentication,
  checkUserAuthenticated,
  fetchConfig,
  logout
}

export default machineUtils;
