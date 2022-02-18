// vendor
import { createMachine, interpret, assign } from 'xstate';
// utils
import {
  authorizeAuthenticateAttempt,
  completeAuthentication,
  checkUserAuthenticated,
  fetchConfig,
  logout,
} from './MachineUtils';
// types
import  {
  Error,
  Config,
  User,
  AuthMachineContext,
  AuthState,
} from './AuthTypes';


type AuthEvent =
  | { type: 'accept', data: Config}
  | { type: 'FETCH_CONFIG' }
  | { type: 'AUTH' }
  | { type: 'LOGGEDIN' }
  | { type: 'LOGGED_OUT' }
  | { type: 'COMPLETE_AUTHENTICATION', data: User }
  | { type: 'USER_LOGGED_IN', actions: any }
  | { type: 'ERROR', error: Error}
  | { type: 'TRY_AGAIN', data: Config }
  | { type: 'loading', data: Config }
  | { type: 'LOGOUT' }
  | { type: 'success' }
  | { type: 'systemError'}
  | { type: 'error', data: Error }
  | { type: 'idleAuth', data: User }
  | { type: 'idle' }
  | { type: 'authorizeAuthenticate', data: User }
  | { type: 'login' };

interface Page {
  waitFor: (value:string) => void;
}

const configuration = {
  id: 'auth',
  initial: 'idle',
  states: {
    idle: {
      on: { FETCH_CONFIG: { target: 'loading'} },
    },
    loading: {
      invoke: {
        id: 'config',
        src: (context: AuthMachineContext):Promise<Config> => fetchConfig('.well-known/openid-configuration'),
        onDone: {
          target: 'checkLogin',
          actions: assign<AuthMachineContext, { type: 'FETCH_CONFIG', data: Config }>({
              config: (context, event) => event.data
          }),
        },
        onError: {
          target: 'systemError',
          actions: assign({ error: (context, event) => event })
        },
      },
    },
    checkLogin: {
      invoke: {
        id: 'checkLogin',
        src: (context: AuthMachineContext) => (callback:any, onRecieve:any):any => checkUserAuthenticated(callback, onRecieve, context),
      },
      on: {
        LOGGED_OUT: { target: 'loggedOut' },
        REMOVE_USER: { target: 'loggedOut' },
        COMPLETE_AUTHENTICATION: { target: 'completeAuthentication' },
        USER_LOGGED_IN: { target: 'loggedIn' },
      }
    },
    loggedIn: {
      on: { LOGOUT: {
        target: 'logout',
      }},
    },
    logout: {
      invoke: {
        id: 'config',
        src: (context: AuthMachineContext):Promise<any> => logout(context),
        onDone: {
          target: 'idle',
          actions: assign<AuthMachineContext, { type: 'LOGOUT', data: Config }>({
              config: (context, event) => {
                context.user = null;
                return event.data;
              }
          })
        },
        onError: {
          target: 'error',
          actions: assign({ error: (context, event) => event })
        }
      },
    },
    loggedOut: {
      on: { AUTH: { target: 'authorizeAuthenticate'} }
    },
    authorizeAuthenticate: {
      invoke: {
        id: 'authorizeAuthenticate',
        src: (context: AuthMachineContext):Promise<any> => authorizeAuthenticateAttempt(context),
        onDone: {
          target: 'redirect',
          actions: assign<AuthMachineContext, { type: 'AUTH', data: User }>({
            user: (context, event) => event.data
          })
        },
        onError: {
          target: 'error',
          actions: assign({ error: (context, event) => event })
        }
      },
    },
    completeAuthentication: {
      invoke: {
        id: 'completeAuthentication',
        src: (context: AuthMachineContext):Promise<any> => completeAuthentication(context),
        onDone: {
          target: 'loggedIn',
          internal: false,
          actions: assign<AuthMachineContext,  { type: 'COMPLETE_AUTHENTICATION', data: User }>({
            user: (context, event) => event.data
          })
        },
        onError: {
          target: 'error',
          actions: assign({ error: (context, event) => event })
        }
      },
    },
    error: {
      on: {
        TRY_AGAIN: {
          target: 'loggedOut'
        },
      },
    },
    systemError: {
    },
    redirect: {
    }
  }
};

const authMachine = createMachine<AuthMachineContext, AuthEvent, AuthState>(configuration);


// TODO add dev only wrapper for this code, commented out for now
// // Edit your service(s) here
// const service = interpret(authMachine, { devTools: true }).onTransition(state => {
//   console.log(state.value);
// });
//
// service.start();

// service.send("FETCH_CONFIG");
// service.send("AUTH");


export {
  configuration,
}


export default authMachine;
