// vendor

import { createContext } from 'react';

interface Profile {
  name: string;
  role: string;
  groups: string;
  email: string;
  given_name: string;
  family_name: string;
  nickname: string;
}

interface UserAuthManager {
  getUser: () => Promise<User>;
}

interface User {
  profile: Profile;
}

interface AppContextType {
  wellKnown?: any;
  user?: User;
  userAuthManager?: UserAuthManager;
}

const appContext:AppContextType = {
    wellKnown: null,
    user: null
};


const AppContext = createContext(appContext)

export default AppContext;
