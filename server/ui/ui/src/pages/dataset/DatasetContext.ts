// vendor
import { createContext, Context } from 'react';

const DatasetContext: Context<{ send: any; }> = createContext({ send: null });


export default DatasetContext;
