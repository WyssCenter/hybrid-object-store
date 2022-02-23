// vendor
import React,
  {
    FC,
    useCallback,
    useEffect,
    useContext,
    useState,
  } from 'react';
import {
  useParams,
  useLocation,
} from "react-router-dom";
import { useMachine } from '@xstate/react';
import ReactTooltip from 'react-tooltip';
// context
import AppContext from 'Src/AppContext';
// enivornment
import { get, del } from 'Environment/createEnvironment';
// components
import { HierarchyHeader } from 'Components/header/index';
import {SectionCard} from 'Components/card/index';
import {WarningButton} from 'Components/button/index';
import Error from '../machine/error/Error';
import Loading from '../machine/loading/Loading';
import NamespaceList from './list/NamespaceList';
import SyncSection from './sync/SyncSection';
import DeleteModal from './modal/DeleteModal';
// context
import NamespaceContext from './NamespaceContext'
// machine
import namespaceMachine from '../machine/PageMachine';
// css
import './Namespace.scss';

type Data = {
  namespace?: string,
  created: string,
  description: string,
  owner: string,
  root_directory: string,
  name: string,
  sync_enabled: boolean,
  sync_type: string,
  delete_status: string,
}

interface ParamTypes {
  namespace: string;
}

interface RenderMap {
  [key: string]: JSX.Element | undefined;
  idle?: JSX.Element;
  loading?: JSX.Element;
  error?: JSX.Element;
  success?: JSX.Element;
}

interface Props {
  setSyncValue: any;
}

/**
* Method fetches dataset data and handles state changes
* @param {}
* @calls {environment#get}
* @calls {macine#send}
* @return {void}
*/
const Namespace: FC<Props> = ({ setSyncValue }:Props) => {
  const [ isModalVisible, setIsModalVisible ] = useState(false);

  // params
  const { namespace } = useParams<ParamTypes>();
  // machine
  const [state, send] = useMachine(namespaceMachine);
  const location = useLocation();
  const paths = location.pathname.split('/')
  const viewSettings = (paths[2] === 'settings');
  const { user } = useContext(AppContext);
  /**
  * Method sends fetches namespace data and its datasets meta data
  * @param {}
  * @return {void}
  * @fires {#get}
  * @fires {#send}
  */
  const fetchNamespaceData = useCallback(() => {
    Promise.all([
      get(`namespace/${namespace}`),
      get(`namespace/${namespace}/dataset/`),
      get(`namespace/${namespace}/sync`)
    ])
      .then(responses => Promise.all(responses.map(response => response.json())))
      .then(([namespaceData, datasetData, syncData]) => {
        const flattenedDatasets = datasetData.map((data: Data) => {
          return {
            name: data.name,
            description: data.description,
            created: data.created,
            directory: data.root_directory,
            sync: data.sync_enabled ? data.sync_type : 'disabled',
            showDeleteBadge: data.delete_status && (data.delete_status !== 'NOT_SCHEDULED')
          }
        });
        setSyncValue([{
          section: 'namespace',
          data: syncData
        }]);
        send("SUCCESS", {
            namespace: namespaceData,
            datasets: flattenedDatasets,
            sync: syncData,
        });
      })
      .catch((error: Error) => {
        const newErrorMessage = error.toString ? error.toString() : error;
        send("ERROR", { error: newErrorMessage });
      });
  }, [send, namespace])

  useEffect(()=> {
    if (state.value === 'idle') {
      send('SUBMIT');
      fetchNamespaceData();
    }


    if (state.value === 'refetching') {
      fetchNamespaceData();
    }

    return () => {
      // send('RESET');
    }
  }, [send, fetchNamespaceData, state.value]);

  useEffect(() => {
    return () => {
      send('RESET');
    }
  }, [send]);

  useEffect(() => {
    if (user && user.profile) {
      const { email } = user.profile;
      const redirectObject = {
        [email]: window.location.pathname,
      };
      const oldRedirectObject = localStorage.getItem('redirect');
      if (oldRedirectObject) {
        const newRedirectObject = Object.assign({}, JSON.parse(oldRedirectObject), redirectObject);
        localStorage.setItem('redirect', JSON.stringify(newRedirectObject))
      } else {
        localStorage.setItem('redirect', JSON.stringify(redirectObject))
      }
    }
  }, [user])

  const stateValue: string = typeof state.value === 'string' ? state.value : 'idle';

  const deleteNamespace = () => {
    del(`namespace/${namespace}`).then((response: Response) => {
      window.location.pathname = '/ui/';
    }).catch((error: Error) => {
      console.log(error);
    });
  }

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
      <NamespaceList
        datasets={state.event.datasets}
        namespaceData={state.event.namespace}
      />
    ),
    error: (
      <Error errorMessage={state.event.error} />
    )
  };

  const renderSection = () => {
    if (viewSettings) {
      return (
        <>
        <SyncSection
          sync={state.event.sync}
          namespace={namespace}
        />
      <section>
        <h4>Delete Namespace</h4>

        <DeleteModal
          namespaceName={namespace}
          handleDeleteClick={deleteNamespace}
          hideModal={() => setIsModalVisible(false)}
          isVisible={isModalVisible}
        />
        <SectionCard>
          <div className="Namespace__Delete flex align-items--center justify--space-between">
            <p>
             All datasets in a namespace must be deleted and namespace syncing must be disabled before deletion. Deleting a namespace cannot be reverted.
            </p>
            <div
              className="Namespace__Delete--buttons relative"
              data-tip="Only administrators can delete a namespace."
              data-tip-disable={user.profile.role === 'admin'}
            >
              <WarningButton
                click={() => setIsModalVisible(true)}
                disabled={user.profile.role !== 'admin'}
                text="Delete Namespace"
              />
              <ReactTooltip
                place="bottom"
                effect="solid"
              />
            </div>
          </div>
        </ SectionCard>
      </section>
        </>
      )
    }
    return renderMap[stateValue]
  }

  return (
      <div className="Namespace margin--auto column-1-span-12">
        <div className="grid">
          <NamespaceContext.Provider value={{send }}>
              {renderSection()}
          </NamespaceContext.Provider>


        </div>

    </div>
  );
}


export default Namespace;
