import { UserManager } from 'oidc-client';

export interface Config {
  wellKnown: {
    issuer: string;
    authorization_endpoint: string;
    end_session_endpoint: string;
    jwks_uri: string;
    userinfo_endpoint: string;
  },
  userAuthManager: UserManager;
}

interface Profile {
  exp: number;
}

export interface User {
  expired: boolean;
  profile: Profile;
  state: any;
}

export interface Error {
  message: string;
}

export interface AuthMachineContext {
  config?: Config | undefined;
  user?: User | undefined
  error?: Error | undefined
}

export type AuthState =
  | {
      value: 'idle';
      context: AuthMachineContext & {
        config: undefined;
        error: undefined;
      };
    }
  | {
      value: 'loading';
      context: AuthMachineContext & { config: Config | undefined; user: User | undefined };
    }
  | {
      value: 'checkLogin';
      context: AuthMachineContext & { config: Config | undefined; user: User | undefined; };
    }
  | {
      value: 'loggedIn';
      context: AuthMachineContext & { config: Config; user: User };
    }
  | {
      value: 'logout';
      context: AuthMachineContext & { config: Config; user: User | undefined };
    }
  | {
      value: 'authorizeAuthenticate';
      context: AuthMachineContext & { config: Config; user: User | undefined };
    }
  | {
      value: 'completeAuthentication';
      context: AuthMachineContext;
    }
  | {
      value: 'error';
      context: AuthMachineContext & { config: Config | undefined; user: User | undefined };
    }
  | {
      value: 'systemError';
      context: AuthMachineContext & { config: undefined; user: undefined  };
    }
  | {
      value: 'redirect';
      context: AuthMachineContext & { config: Config; user: User };
    };
