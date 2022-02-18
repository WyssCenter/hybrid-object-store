// vendor
import React, { FC, useCallback, useRef, useState, useContext } from 'react';
import { useParams } from 'react-router-dom';
import classNames from 'classnames';
import ReactTooltip from 'react-tooltip';
// app context
import AppContext from 'Src/AppContext';
// environemnt
import { put, del } from 'Environment/createEnvironment';
// components
import {SectionCard} from 'Components/card/index';
import {PrimaryButton} from 'Components/button/index';
import {TooltipConfirm} from 'Components/tooltip/index';
import SyncModal from './modal/SyncModal';
import ConfigureModal from './configure/ConfigureModal';
// assets
import { ReactComponent as EnabledArrow } from 'Images/icons/line-enabled.svg';
import { ReactComponent as DisabledArrow } from 'Images/icons/line-disabled.svg';
import './SyncDataset.scss';

// types
import {
  Dataset,
} from '../section/DatasetSectionTypes';

interface Props {
  dataset: Dataset;
  send: any;
  namespaceSync: any;
}

interface ParamTypes {
  datasetName: string;
  namespace: string;
}



const SyncDataset: FC<Props> = ({ dataset, send, namespaceSync }: Props) => {
  // context
  const { user } = useContext(AppContext);
  // state
  const [hasError, setHasError] = useState(false);
  const [ isModalVisible, setIsModalVisible ] = useState(false);
  const [ configModalVisible, setConfigModalVisible ] = useState(false);
  const [ configError, setConfigError ] = useState(null);
  const [ isTooltipVisible, setIsTooltipVisible] = useState(false);
  const [ isButtonEnabled, setIsButtonEnabled] = useState(true);
  // refs
  const enableTooltipRef = useRef();
  const disableTooltipRef = useRef();
  // params
  const { datasetName, namespace } = useParams<ParamTypes>();
  const syncEnabled = dataset.sync_enabled;

  // click functions
  const handleEnableClick = useCallback((syncType) => {
    setIsButtonEnabled(false);
    setIsModalVisible(false);

    put(
      `namespace/${namespace}/dataset/${datasetName}/sync`,
      {
        sync_type: syncType,
      }
    ).then(response => {
      setIsButtonEnabled(true);
      if (response.ok) {
        send('REFETCH');
      } else {
        setHasError(true)
      }
    }).catch(error => {
      console.log(error)
    });
  }, [namespace, datasetName, send])

    // click functions
  const handleUpdate = useCallback((syncType, policy) => {
    setIsButtonEnabled(false);
    setIsModalVisible(false);
    try {
      JSON.parse(policy);
    } catch (e) {
      setConfigError('Policy entered is not valid JSON.');
      setIsButtonEnabled(true);
      setTimeout(() => {
        setConfigError(null);
      }, 5000)
      return;
    }
    put(
      `namespace/${namespace}/dataset/${datasetName}/sync`,
      {
        sync_type: syncType,
        sync_policy: JSON.stringify(JSON.parse(policy)),
      }
    ).then(response => {
      setIsButtonEnabled(true);
      if (response.ok) {
        send('REFETCH');
        return;
      } else {
        setHasError(true)
      }
      return response.json();
    })
    .then((data) => {
      if (data && data.error) {
        setConfigError(data.error);
        setTimeout(() => {
          setConfigError(null);
        }, 5000)
      }
    })
    .catch(error => {
      setIsButtonEnabled(true);
      console.log(error)
    });
  }, [namespace, datasetName, send])


  const handleDisableClick = useCallback(() => {
    setIsModalVisible(false);
    setIsButtonEnabled(false);

    del(`namespace/${namespace}/dataset/${datasetName}/sync`).then(response => {
      setIsButtonEnabled(true);
      if (response.ok) {
        send('REFETCH');
      } else {
        setHasError(true)
      }
    }).catch(error => {
      console.log(error)
    });
  }, [namespace, datasetName, send])

  let syncText = 'Sync Disabled';
  let buttonText = 'Enable Sync'
  if (dataset.sync_enabled && dataset.sync_type === 'simplex') {
    syncText = '1-way Sync Enabled';
    buttonText = 'Disable Sync'
  } else if (dataset.sync_enabled && dataset.sync_type === 'duplex') {
    syncText = '2-way Sync Enabled';
    buttonText = 'Disable Sync'
  }

  let disabledMessage = '';
  if (namespaceSync && namespaceSync.length === 0 && !dataset.sync_enabled) {
    disabledMessage = 'Namespace Sync must be enabled first';
  }
  if (user.profile.role === 'user') {
    disabledMessage = 'Only admins and privileged users can change the configuration';
  }

  const targetHostname = () => {
    if(!namespaceSync || !namespaceSync[0]) {
      return '';
    } else {
      return new URL(namespaceSync[0].target_core_service).hostname;
    }
  }

  const squareCSS = classNames({
    SyncDataset__square: true,
    'SyncDataset__square--disabled': !dataset.sync_enabled,
  })

  const arrowCSS = classNames({
    SyncDataset__arrows: true,
    'SyncDataset__arrows--disabled': !dataset.sync_enabled,
    'SyncDataset__arrows--simplex': dataset.sync_enabled && dataset.sync_type === 'simplex',
    'SyncDataset__arrows--duplex': dataset.sync_enabled && dataset.sync_type === 'duplex',
  });

  return (
    <section>
     <h3>Dataset Sync Configuration</h3>
     <SyncModal
        handleEnableClick={handleEnableClick}
        handleDisableClick={handleDisableClick}
        hideModal={() => setIsModalVisible(false)}
        namespaceSync={namespaceSync}
        isVisible={isModalVisible}
        syncEnabled={syncEnabled}
     />
     <ConfigureModal
        handleUpdate={handleUpdate}
        hideModal={() => setConfigModalVisible(false)}
        namespaceSync={namespaceSync}
        isVisible={configModalVisible}
        syncEnabled={syncEnabled}
        syncPolicy={dataset.sync_policy}
        syncType={dataset.sync_type}
        configError={configError}
        isButtonEnabled={isButtonEnabled}
     />

      <SectionCard>
      <div className="SyncDataset flex">
        <div className="SyncDataset__blurb column-1-span-6">
          <p>
            A dataset can be configured for 1-way or 2-way syncing. In 1-way syncing, data will be replicated from the dataset in this namespace to the dataset in the linked namespace. In 2-way syncing, data will be replicated in both directions.
          </p>
          <p>
            The dataset will automatically be created in the linked namespace when dataset syncing is enabled.
          </p>
          <p>
            Write and delete operations will be syncronized automatically.
          </p>
          <p>
            Note, if you wish to delete the dataset and leave it intact in the synced location, you must first disable syncing before performing the delete!
          </p>
        </div>
        <div className="flex flex--column align-items--center justify--center  column-1-span-6">
          <div className="SyncDataset__diagram">
            <div>
              <div className={squareCSS}>
                {(namespaceSync && namespaceSync.length === 0) ? '' : namespace}
                <div className="SyncDataset__square--host">
                {(namespaceSync && namespaceSync.length === 0) ? '' : `host: ${window.location.hostname}`}
                </div>
              </div>
            </div>
            <div className={arrowCSS}>
              <div className="SyncDataset__arrows--left"
              >
                {
                  dataset.sync_enabled &&
                  <EnabledArrow />
                }
                {
                  !dataset.sync_enabled &&
                  <DisabledArrow />
                }
              </div>
              <div className="SyncDataset__arrows--right">
                {
                  (dataset.sync_enabled && dataset.sync_type === 'duplex') &&
                  <EnabledArrow />
                }
                {
                  !(dataset.sync_enabled && dataset.sync_type === 'duplex') &&
                  <DisabledArrow />
                }
              </div>
            </div>
            <div className={squareCSS}>
              {(namespaceSync && namespaceSync[0] && namespaceSync[0].target_namespace) || ''}
              <div className="SyncDataset__square--host">
                {(namespaceSync && namespaceSync.length === 0) ? '' : `host: ${targetHostname()}`}
              </div>
            </div>
          </div>
          <div className="SyncDataset__label">
            {syncText}
          </div>
        </div>
      </div>
      <div className="flex flex-1 justify--center">
        <div
          data-tip={disabledMessage}
          role="presentation"
        >
          <PrimaryButton
            click={() => {
              if (!dataset.sync_enabled) {
                setConfigModalVisible(true);
              } else {
                setIsModalVisible(true)
              }
            }}
            text={buttonText}
            disabled={disabledMessage.length > 0}
          />
          <ReactTooltip
            place="bottom"
            effect="solid"
          />
        </div>
        <div
          data-tip={dataset.sync_enabled ? disabledMessage : 'Enable Dataset Sync to configure'}
          role="presentation"
        >
          <PrimaryButton
            click={() => setConfigModalVisible(true)}
            text="Configure Sync"
            disabled={!dataset.sync_enabled || disabledMessage.length > 0}
          />
          <ReactTooltip
            place="bottom"
            effect="solid"
          />
        </div>
      </div>
      </ SectionCard>
    </section>
  );
}


export default SyncDataset;
