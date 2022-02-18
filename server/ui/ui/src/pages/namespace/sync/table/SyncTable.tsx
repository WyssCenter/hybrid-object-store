// vendor
import React, {
  FC,
  useContext,
  useState,
  useCallback,
  useRef,
} from 'react';
import { faPlus } from '@fortawesome/free-solid-svg-icons';
import classNames from 'classnames';
import ReactTooltip from 'react-tooltip';
// app context
import AppContext from 'Src/AppContext';
import NamespaceContext from '../../NamespaceContext';
// environment
import { del } from 'Environment/createEnvironment';
// components
import { PrimaryButton } from 'Components/button/index';
import SyncFormModal from './form/SyncFormModal';
import Row from './row/Row';
import DeleteModal from './modal/DeleteModal';
import { TooltipConfirm } from 'Components/tooltip/index';
// assets
import { ReactComponent as EnabledArrow } from 'Images/icons/line-enabled.svg';
import { ReactComponent as DisabledArrow } from 'Images/icons/line-disabled.svg';
import './SyncTable.scss';

interface Target {
  sync_type: string;
  target_core_service: string;
  target_namespace: string;
}

interface Props {
  sync: Array<Target>;
  namespace: string;
}


const SyncTable: FC<Props> = ({sync, namespace}: Props) => {
  // context
  const { user } = useContext(AppContext);
  const { send } = useContext(NamespaceContext)
  // refs
  const tooltipRef = useRef(null);
  // state
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [isDeleteModalVisible, setIsDeleteModalVisible] = useState(false);
  const [isTooltipVisible, setIsTooltipVisible] = useState(false);

    // events
  /**
  * Method updates tooltip visibility
  * @param {}
  */
  const shopwIsTooltipVisible = useCallback(() => {
    setIsTooltipVisible(true);
  }, []);
  /**
  * Method updates tooltip visibility
  * @param {}
  */
  const hideIsTooltipVisible = () => {
    setIsTooltipVisible(false);
  }

  /**
  * Method removes sync item from table
  * @param {event} Object
  */
  const removeSyncItem:any = useCallback((event: MouseEvent) => {
    del(`namespace/${namespace}/sync`, sync[0]).then((response) => {

      if(response.ok) {
        setIsTooltipVisible(false);
        send("REFETCH")
      }
    }).catch((error) => {
      console.log(error);
    });
  }, [namespace, send, sync]);

  // vars
  const syncActive = sync && (sync.length > 0);

  const permissions = user.profile.role !== 'user';



  let syncText = 'Sync Disabled';
  let buttonText = 'Enable Sync'
  if (sync && sync[0] && sync[0].sync_type === 'simplex') {
    syncText = '1-way Sync Enabled';
    buttonText = 'Disable Sync'
  } else if (sync && sync[0] && sync[0].sync_type === 'duplex') {
    syncText = '2-way Sync Enabled';
    buttonText = 'Disable Sync'
  }

  const targetHostname = () => {
    if(!sync || !sync[0]) {
      return '';
    } else {
      return new URL(sync[0].target_core_service).hostname;
    }
  }

  const syncDisabled = !sync || (sync && sync.length === 0);

  const buttonAction = () => {
    if (sync && sync[0]) {
      setIsDeleteModalVisible(true);
    } else {
      setIsModalVisible(true)
    }
  }

  const squareCSS = classNames({
    SyncTable__square: true,
    'SyncTable__square--disabled': syncDisabled,
  })

  const arrowCSS = classNames({
    SyncTable__arrows: true,
    'SyncTable__arrows--disabled': syncDisabled,
    'SyncTable__arrows--simplex': sync && sync[0] && sync[0].sync_type === 'simplex',
    'SyncTable__arrows--duplex': sync && sync[0] && sync[0].sync_type === 'duplex',
  });

  return (
    <div className="flex flex--column">
      <div className="SyncTable flex">
        <SyncFormModal
          hideModal={() => setIsModalVisible(false)}
          isVisible={isModalVisible}
          sync={sync}
        />
        <DeleteModal
          namespaceName={namespace}
          handleDeleteClick={removeSyncItem}
          hideModal={() => setIsDeleteModalVisible(false)}
          isVisible={isDeleteModalVisible}
        />
        <div className="SyncTable__blurb column-1-span-6">
          <p>
            A namespace must be configured for syncing before datasets in the namespace can be synced. This action links two namespaces together and tells the server to monitor the associated buckets for events.
          </p>
          <p>
            A namespace can be configured for 1-way or 2-way syncing. In 1-way syncing, data will be replicated from this namespace to the linked namespace. In 2-way syncing, data will be replicated in both directions.
          </p>
          <p>
            A namespace that is configured for 1-way syncing can only support 1-way dataset syncing. A namespace that is configured for 2-way syncing will support both 1-way and 2-way dataset syncing, depending on the individual dataset’s sync configuration.
          </p>
          <p>
            When syncing between servers, typically the name of the linked namespace is the same to minimize code changes for users, but the names can be different if your use case requires it.
          </p>
        </div>
        <div className="flex flex--column align-items--center justify--center  column-1-span-6">
          <div className="SyncTable__diagram">
            <div>
              <div className={squareCSS}>
                {syncDisabled ? '' : namespace}
                <div className="SyncTable__square--host">
                  {syncDisabled ? '' : `host: ${window.location.hostname}`}
                </div>
              </div>
            </div>
            <div className={arrowCSS}>
              <div className="SyncTable__arrows--left"
              >
                {
                  !syncDisabled &&
                  <EnabledArrow />
                }
                {
                  syncDisabled &&
                  <DisabledArrow />
                }
              </div>
              <div className="SyncTable__arrows--right">
                {
                  (sync && sync[0] && sync[0].sync_type === 'duplex') &&
                  <EnabledArrow />
                }
                {
                  !(sync && sync[0] && sync[0].sync_type === 'duplex') &&
                  <DisabledArrow />
                }
              </div>
            </div>
            <div className={squareCSS}>
              {(sync && sync[0] && sync[0].target_namespace) || ''}
              <div className="SyncTable__square--host">
                {syncDisabled ? '' : `host: ${targetHostname()}`}
              </div>
            </div>
          </div>
          <div className="SyncTable__label">
            {syncText}
          </div>
        </div>
      </div>
      <div className="flex flex-1 justify--center">
        <div ref={tooltipRef}>
          {
            (user.profile.role === 'admin') && (
              <PrimaryButton
                click={buttonAction}
                text={buttonText}
              />
            )
          }
          {
            (user.profile.role !== 'admin') && (
            <>
              <div
                data-tip="Only admin users can change the configuration"
                role="presentation"
              >
                <PrimaryButton
                  click={buttonAction}
                  text={buttonText}
                  disabled
                />
              </div>
              <ReactTooltip
                place="bottom"
                effect="solid"
              />
            </>
            )
          }
        </div>
      </div>
    </div>
  );
}

export default SyncTable;
