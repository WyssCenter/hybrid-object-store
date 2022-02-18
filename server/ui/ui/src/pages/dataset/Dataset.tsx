// vendor
import React, { FC, useCallback, useEffect, useState, useRef, useContext } from 'react';
import classNames from 'classnames';
import { useMachine } from '@xstate/react';
import moment from 'moment';
import {
  useParams,
  useLocation,
} from "react-router-dom";
// enivornment
import { get, put } from 'Environment/createEnvironment';
import AppContext from 'Src/AppContext';
// components
import { PrimaryButton } from 'Components/button/index';
import { SectionCard } from 'Components/card/index';
import FileBrowser from './filebrowser/FileBrowser';
import Error from '../machine/error/Error'
import Loading from '../machine/loading/Loading';
import DatasetSection from './section/DatasetSection';
import DeleteDataset from './delete/DeleteDataset';
import SyncDataset from './sync/SyncDataset';
import DatasetContext from './DatasetContext';
// machine
import datasetMachine from '../machine/PageMachine';
// css
import './Dataset.scss';


interface RenderMap {
  [key: string]: JSX.Element | undefined;
  idle?: JSX.Element;
  loading?: JSX.Element;
  refetching?: JSX.Element;
  error?: JSX.Element;
  success?: JSX.Element;
}

interface Props {
  setSyncValue: any,
  setDatasetDescription: (description: string) => void;
}

type Data = {
  namespace?: string,
  error?: string,
  sync_enabled: boolean;
}

interface ParamTypes {
  datasetName: string;
  namespace: string;
}

const Dataset: FC<Props> = ({ setSyncValue, setDatasetDescription }:Props) => {
  // state
  const [namespaceSync, setNamespaceSync] = useState(null)
  const [detailsVisible, setDetailsVisible] = useState(false);
  const [columns, setcolumns] = useState(12);

  const { user } = useContext(AppContext);


  const detailsRef = useRef(null);

  detailsRef.current = detailsVisible;

  // machine
  const [state, send ] = useMachine(datasetMachine);
  // params
  const { datasetName, namespace } = useParams<ParamTypes>();
  // vars
  const stateValue: string = typeof state.value === 'string' ? state.value : 'idle';
  const location = useLocation();
  const paths = location.pathname.split('/')
  const viewSettings = (paths[3] === 'settings');

  const restoreDataset = () => {
    put(
      `namespace/${namespace}/dataset/${datasetName}/restore`,
      {
      }
    ).then((response) => {
      if (response.ok) {
        send('REFETCH');
      }
    }).catch((error) => {
      console.log(error);
    })
  }

  /**
  * Method fetches dataset data and handles state changes
  * @param {}
  * @calls {environment#get}
  * @calls {macine#send}
  * @return {void}
  */
  const fetchDatasetData = useCallback(() => {
    Promise.all([
      get(`namespace/${namespace}/dataset/${datasetName}`),
      get(`namespace/${namespace}/sync`)
    ])
    .then(responses => Promise.all(responses.map(response => response.json())))
    .then(([data, namespaceData]) => {
      setNamespaceSync(namespaceData);
      setSyncValue([{
        section: 'namespace',
        data: namespaceData,
      },
      {
        section: 'dataset',
        data: {
          syncEnabled: data.sync_enabled,
          syncType: data.sync_type,
        }
      }]
      )
      setDatasetDescription(data.description)
      if ((typeof data === 'string') || data.error) {
        send("ERROR", { error: data.error || data});
      } else {
        send("SUCCESS", { data: data });
      }
    })
    .catch((error: Error) => {
      const newErrorMessage = error.toString ? error.toString() : error;
      send("ERROR", {error: newErrorMessage});
    });
  }, [send, get]);

  useEffect(()=> {
    if (state.value === 'idle') {
      send('SUBMIT');
      fetchDatasetData();
    }

    if (state.value === 'refetching') {
      send('SUBMIT');
      fetchDatasetData();
    }
  }, [send, fetchDatasetData, state.value])

  useEffect(() => {
    return () => {
      send('RESET');
    }
  }, [send]);

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
    success: (<>
        <DatasetSection
          dataset={state.event.data}
          datasetName={datasetName}
          namespace={namespace}
        />
        <SyncDataset
          dataset={state.event.data}
          send={send}
          namespaceSync={namespaceSync}
        />
        <DeleteDataset
          pendingDelete={state.event.data && (state.event.data.delete_status !== 'NOT_SCHEDULED')}
          route={`namespace/${namespace}/dataset/${datasetName}`}
          redirectRoute={`/ui/${namespace}/`}
          send={send}
        />
      </>
    ),
    error: (
      <Error errorMessage={state.event.error} />
    )
  };

    const handleResize = () => {
      if(detailsRef.current) {
        let columns = 5;
        const width = window.innerWidth - 1100;
        if (width > 150) {
          columns += Math.floor(width / 100)
          if (columns >= 13) {
            columns = 12;
          }
        }
        setcolumns(columns);
      }
    }

  useEffect(() => {
    window.addEventListener("resize", handleResize);
    handleResize();
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  useEffect(() => {
    handleResize();
  }, [detailsVisible])

  const datasetContainerCSS = classNames({
    DatasetContainer: true,
    'DatasetContainer--details': detailsVisible,
  })

  const sectionCSS = classNames({
    'column-1-span-12': !detailsVisible,
    [`column-1-span-${columns}`]: detailsVisible,
  })

  const renderSection = () => {
    if (viewSettings) {
      return (
        renderMap[stateValue]
      )
    }
    return (
      <>
      {
        state.event.data && (state.event.data.delete_status !== 'NOT_SCHEDULED') && (
          <div className="Dataset__delete">
            <div className="Dataset__delete--main">{`This dataset has been scheduled  for deletion and all data will be permanently removed on ${moment.utc(state.event.data.delete_on).format('lll')} (UTC)`}</div>
            {
              user.profile.role !== 'admin' && (state.event.data.delete_status !== 'IN_PROGRESS') && (
                <div className="Dataset__delete--sub">Contact an administrator before this time to restore the data.</div>
              )
            }
            {
              user.profile.role === 'admin' && (state.event.data.delete_status !== 'IN_PROGRESS') && (
              <PrimaryButton
                click={restoreDataset}
                text="Click here to restore this dataset."
              />
              )
            }
          </div>
        )
      }
      <SectionCard noPadding span={detailsVisible ? columns : 12}>
        <div className="flex align-items--center flex--column">
          {
            state.event.data && (
              <FileBrowser
                bucket={state.event.data.namespace.bucket_name}
                namespace={state.event.data.namespace.name}
                dataset={datasetName}
                setDetailsVisible={setDetailsVisible}
                detailsVisible={detailsVisible}
                lockBrowser={(state.event.data.delete_status !== 'NOT_SCHEDULED')}
              />
            )
          }
        </div>
      </SectionCard>
      </>
    )
  }

  return (
    <DatasetContext.Provider value={{ send }}>
      <div className={datasetContainerCSS}>
      <section className={sectionCSS}>
      </section>
      <div className="Dataset">
        <div className={`Dataset__section ${sectionCSS}`}>
          {renderSection()}
        </div>
      </div>
      </div>
    </DatasetContext.Provider>
  );
}


export default Dataset;
