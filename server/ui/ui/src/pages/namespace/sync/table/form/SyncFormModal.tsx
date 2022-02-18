// vendor
import React, {
  FC,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
} from 'react';
// params
import {
  useParams,
} from "react-router-dom";
// context
import NamespaceContext from '../../../NamespaceContext';
// environment
import { put } from 'Environment/createEnvironment';
// components
import Modal from 'Components/modal/Modal';
import Checkbox from 'Components/form/checkbox/index';
import { PrimaryButton } from 'Components/button/index';
import { InputText } from 'Components/form/text/index';
import { RadioButtons } from 'Components/form/radio/index';
// css
import './SyncFormModal.scss';

interface Props {
  hideModal: () => void;
  isVisible: boolean;
  sync: any;
}

interface ParamTypes {
  namespace: string;
}

const SyncForm: FC<Props> = ({
  hideModal,
  isVisible,
  sync,
}: Props) => {

  // context
  const { send } = useContext(NamespaceContext);
  // params
  const { namespace } = useParams<ParamTypes>();
  // ref
  const inputTextRef = useRef(null);
  const inputTextNamespaceRef = useRef(null);
  // state
  const [isHttpsChecked, setHttpsChecked] = useState(true);
  const [hostname, setHostname] = useState('');
  const [namespaceFormItem, setNamespace] = useState('');
  const [error, setError] = useState(null);
  const [radioValue, setRadioValue] = useState('simplex');


  // vars
  const urlPrefix = isHttpsChecked ? 'https://' : 'http://';
  const radioButtonValues = [
    {value: 'simplex', text: '1 Way Syncing'},
    {value: 'duplex', text: '2 Way Syncing'}
  ];


  /**
  * Method updates state for radio button
  */
  const updateRadioValue = (radioValue: string) => {
    setRadioValue(radioValue)
  }

  /**
  * Method updates state for checkbox
  * @param {value} boolean
  */
  const updateCheckbox = useCallback((value:boolean) => {
    setHttpsChecked(value);
  }, [setHttpsChecked]);

  /**
  * Method updates state for checkbox
  * @param {string} value
  */
  const updateHostname: any = useCallback((event:KeyboardEvent) => {
    const element = event.currentTarget as HTMLInputElement;
    setHostname(element.value);
  }, [setHostname]);

  /**
  * Method updates state for checkbox
  * @param {string} value
  */
  const updateNamespace: any = useCallback((event:KeyboardEvent) => {
    const element = event.currentTarget as HTMLInputElement;
    setNamespace(element.value);
  }, [setNamespace]);

  /**
  * Method updates state for checkbox
  * @param {}
  */
  const submitTarget = () => {
    const url = `${urlPrefix}${hostname}/core/v1`;
    put(
      `namespace/${namespace}/sync`,
      {
        target_core_service: url,
        target_namespace: namespaceFormItem,
        sync_type: radioValue,
      }
    ).then((response) => {
      if (response.ok) {
        send('REFETCH');
        hideModal();
        inputTextRef.current.value = '';
        inputTextNamespaceRef.current.value = '';
      } else {
        setError('There was a problem adding this sync service.')
      }
    }).catch((error) => {
      console.log(error);
    })
  };

  useEffect(() => {
    setError(null);
  }, [isVisible])

  if(!isVisible) {
    return null;
  }

  return (
    <Modal
      header="Add Sync Target"
      handleClose={hideModal}
      icon={'faSyncAlt'}
      overflow="visible"
      size="medium"
    >
      <section className="SyncForm">
        <p>Add a Hoss server to enable dataset syncing. Make sure that you select the correct protocol for your server.</p>
        <div className="SyncForm__widget">
          <div className="flex">
            <span className="SyncForm__widget--label">Hostname</span>
            <Checkbox
              id="Checkbox__https"
              isChecked={isHttpsChecked}
              label="https"
              updateCheckbox={updateCheckbox}
            />
            <i className="InputText__i InputText__i--orange">* Required</i>
          </div>
          <div className="SyncForm__container">
            <input
              className="SyncForm__input"
              ref={inputTextRef}
              placeholder="localhost"
              onKeyUp={updateHostname}
              type="text"
            />
          </div>
        </div>

        <InputText
          inputRef={inputTextNamespaceRef}
          isRequired={true}
          label="Target Namespace"
          placeholder={namespace}
          updateValue={updateNamespace}
        />

        <RadioButtons
          name="SyncNamespaceModal"
          radioType="Sync Type"
          value={radioValue}
          updateValue={updateRadioValue}
          values={radioButtonValues}
        />


        <div className="SyncForm__buttons">
          <PrimaryButton
            click={submitTarget}
            disabled={hostname.length < 1}
            text="Add Sync Target"
          />
        </div>

        { error && (
          <p className="SyncForm__p--error">{error}</p>
        )}

        </section>
      </Modal>
  );
}

export default SyncForm;
