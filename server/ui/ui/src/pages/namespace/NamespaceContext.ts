// vendor
import { createContext, Context } from 'react';

const NamespaceContext: Context<{ send: any; }> = createContext({ send: null });


export default NamespaceContext;
