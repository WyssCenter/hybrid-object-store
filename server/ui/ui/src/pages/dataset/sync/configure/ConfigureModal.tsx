// vendor
import React, { FC, useState } from 'react';
// components
import Modal from 'Components/modal/Modal';
import { RadioButtons } from 'Components/form/radio/index';
import { PrimaryButton, TernaryButton } from 'Components/button/index';
// css
import './ConfigureModal.scss'



interface Props {
  hideModal: () => void;
  isVisible: boolean;
  namespaceSync: any[];
  syncEnabled: boolean;
  syncPolicy: string;
  handleUpdate: any;
  configError: string;
  syncType: string;
  isButtonEnabled: boolean;
}



const ConfigureModal: FC<Props> = (
  {
    hideModal,
    namespaceSync,
    isVisible,
    syncEnabled,
    syncPolicy,
    handleUpdate,
    configError,
    syncType,
    isButtonEnabled,
  }: Props
) => {
  // states
  const [radioValue, setRadioValue] = useState(syncType || 'simplex');
  const [policy, setPolicy] = useState(syncPolicy || '{"Version": "1", "Effect": "OR", "Statements":[]}');

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
      header="Configure Dataset Sync"
      size="large-full"
    >
      <section className="ConfigureModal flex">
        <div className="ConfigureModal__section">
          <div className="ConfigureModal--centered">Select Sync Type</div>
          <div>
            1-way will sync write & delete operations in this dataset to the same dataset in the target namespace 2-way will sync write & delete operations in both directions.
          </div>
          <div>
            Note, 2-way sync is only available if the namespace is configured for 2-way syncing.
          </div>
          <RadioButtons
            name="SyncDatasetModal"
            radioType="Sync Type"
            value={radioValue}
            updateValue={updateRadioValue}
            values={radioButtonValues}
            disabledValues={disabledValues}
          />
        </div>
        <div className="ConfigureModal__section">
          <div className="ConfigureModal--centered">
            Sync Policy
          </div>
          <div>
            Sync policies modify what files and operations will sync. Review the docs
            {' '}
            <a
              href="https://hoss-client.readthedocs.io/en/stable/datasets.html#dataset-syncing"
              target="__blank"
            >
              here
            </a>
            {' '}
            for more information.
          </div>
          <textarea rows={10} cols={50} onChange={(evt) => setPolicy(evt.target.value)} defaultValue={JSON.stringify(JSON.parse(syncPolicy || '{"Version": "1", "Effect": "OR", "Statements":[]}'), undefined, 2)}>
          </textarea>
        </div>
      </section>
        <div className="ConfigureModal__buttons flex justify--flex-end">
            {
            configError && (
            <div className="ConfigureModal__error">
                {configError}
            </div>
            )
          }
          <TernaryButton
            click={() => {
              hideModal();
            }}
            text="Cancel"
          />
          <PrimaryButton
            click={() => {
              handleUpdate(radioValue, policy);
            }}
            disabled={!isButtonEnabled}
            text={isButtonEnabled ? 'Submit' : 'Submitting...'}
          />
        </div>
          {
          !configError && (
          <div className="ConfigureModal__note">
            Note, it may take up to 1 minute for changes to fully propagate throughout the system.
          </div>
          )
        }
    </Modal>
  )
}


export default ConfigureModal;
