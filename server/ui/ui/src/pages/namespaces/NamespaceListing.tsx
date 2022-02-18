// vendor
import React,
{
  FC,
  useCallback,
  useEffect,
} from 'react';
import { useMachine } from '@xstate/react';
// environment
import { get } from 'Environment/createEnvironment';
// Components
import { HierarchyHeader } from 'Components/header/index';
import NamespacesTable from './table/NamespacesTable';
import Error from '../machine/error/Error';
import Loading from '../machine/loading/Loading';
// mahcine
import namespaceListingMachine from '../machine/PageMachine';
// context
import NamespaceListingContext from './NamespaceListingContext';
// css
import './NamespaceListing.scss';

type Data = {
  namespace?: string,
} | string;


interface RenderMap {
  [key: string]: JSX.Element | undefined;
  idle?: JSX.Element;
  loading?: JSX.Element;
  refetching?: JSX.Element;
  error?: JSX.Element;
  success?: JSX.Element;
}


const NamespaceListing: FC = () => {
  // machine
  const [state, send] = useMachine(namespaceListingMachine);

  /**
  * Method shows fetches data for pat list
  * @param {}
  * @return {void}
  * @call {machine#send}
  */
  const getNamespaces = useCallback(() => {
    get('namespace/').then((response: Response) => {
        if (response.statusText === 'Unauthorized') {
          send("ERROR", { error: 'User is not authorized to view these namespaces' });
        } else {
          return response.json();
        }
      })
      .then((data: Data) => {
        if (typeof data === 'string') {
          send("ERROR", { error: data });
        } else {
          send("SUCCESS", { data });
        }
      })
      .catch((error: Error) => {
        const newErrorMessage = error.toString ? error.toString() : error;
        send("ERROR", { error: newErrorMessage });
      })
  }, [send]);

  useEffect(()=> {
    if (state.value === 'idle') {
      send('SUBMIT');
      getNamespaces();
    }

    if (state.value === 'refetching') {
      getNamespaces();
    }
  }, [state.value, getNamespaces, send])

  const stateValue: string = typeof state.value === 'string' ? state.value : 'idle';

  const renderMap:RenderMap = {
    idle: (
      <div />
    ),
    loading: (
      <Loading />
    ),
    refetching: (
      <div />
    ),
    success: (
      <NamespacesTable
        namespaces={state.event.data}
      />
    ),
    error: (
      <Error errorMessage={state.event.error} />
    )
  };

  return (
      <>
        <div className="NamespaceListing">

          <div className="grid">
            <NamespaceListingContext.Provider value={{send}}>
              {renderMap[stateValue]}
            </NamespaceListingContext.Provider>
          </div>
        </div>
      </>
  );
}


export default NamespaceListing;
