// vendor
import { createContext, Context } from 'react';

const NamespaceListingContext: Context<{ send: any; }> = createContext({ send: null });


export default NamespaceListingContext;
