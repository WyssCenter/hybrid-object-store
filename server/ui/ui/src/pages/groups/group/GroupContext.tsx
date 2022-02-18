// vendor
import { createContext, Context } from 'react';



const GroupContext: Context<{
  send: any;
  groupname: string | null;
}> = createContext({
  send: null,
  groupname: null,
});


export default GroupContext;
