// vendor
import React, { FC, useState } from 'react';
// components
import Modal from 'Components/modal/Modal';
import { PrimaryButton } from 'Components/button/index';
import { RadioButtons } from 'Components/form/radio/index';
// css
import './SyncModal.scss'



interface Props {
  handleEnableClick: (syncType: string) => void;
  handleDisableClick: () => void;
  hideModal: () => void;
  isVisible: boolean;
  namespaceSync: any[];
  syncEnabled: boolean;
}



const SyncModal: FC<Props> = (
  {
    handleEnableClick,
    handleDisableClick,
    hideModal,
    namespaceSync,
    isVisible,
    syncEnabled,
  }: Props
) => {
  // states
  const [radioValue, setRadioValue] = useState('simplex');

  // vars
  const radioButtonValues = [
    {value: 'simplex', text: '1 Way Syncing'},
    {value: 'duplex', text: '2 Way Syncing'}
  ];

  /**
  * Method updates state for radio button
  *
  */
  const updateRadioValue = (radioValue: string) => {
    setRadioValue(radioValue)
  }

  const handleClick = () => {
    hideModal()
    if (!syncEnabled) {
      handleEnableClick(radioValue);
    } else {
      handleDisableClick();
    }
  }

  const namespaceSyncType = namespaceSync[0] && namespaceSync[0].sync_type;
  const disabledValues = [];
  if (namespaceSyncType === 'simplex') {
    disabledValues.push('duplex');
  }

  if (!isVisible) {
    return null;
  }
  return (
    <Modal
      handleClose={hideModal}
      header={!syncEnabled ? 'Enable Sync' : 'Disable Sync'}
    >
      <section className="SyncModal">
        <p>{!syncEnabled ? 'Enable sync and choose between 1-way and 2-way syncing.' : 'Disabling syncing will stop the server from automatically synchronizing data in this dataset with configured sync targets.'} </p>
        {
          ((disabledValues.length > 0) && !syncEnabled) && (
            <p>
              2-Way Syncing is disabled for this dataset as the namespace only has 1-way syncing enabled.
            </p>
          )
        }
        { !syncEnabled &&
          <RadioButtons
            name="SyncDatasetModal"
            radioType="Sync Type"
            value={radioValue}
            updateValue={updateRadioValue}
            values={radioButtonValues}
            disabledValues={disabledValues}
          />
        }
        <PrimaryButton
          click={handleClick}
          text={!syncEnabled ? 'Enable Sync' : 'Disable Sync'}
        />
      </section>
    </Modal>
  )
}


export default SyncModal;
