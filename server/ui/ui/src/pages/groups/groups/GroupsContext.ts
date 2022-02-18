// vendor
import { createContext, Context } from 'react';

const GroupsContext: Context<{ send: any; }> = createContext({ send: null });


export default GroupsContext;
