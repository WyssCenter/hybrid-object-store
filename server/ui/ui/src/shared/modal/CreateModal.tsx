// vendor
import React, {
  FC,
  useEffect,
  useRef,
  useState,
  KeyboardEvent,
} from 'react';
import { useMachine } from '@xstate/react';
// environment
import { post } from 'Environment/createEnvironment';
// components
import Modal from 'Components/modal/Modal';
import PrimaryButton, { TernaryButton } from 'Components/button/index';
import InputText from 'Components/form/text/index';
import Processing from './processing/Processing';
import Success from './success/Success';
import Error from './error/Error';
import NamespaceFields from './namespace/NamespaceFields';
// machine
import ModalMachine from './machine/ModalMachine';
// css
import './CreateModal.scss';

interface ObjectStore {
  name: string;
  type: string;
}

interface Props {
  handleClose: () => void;
  isVisible: boolean;
  modalType: string;
  objectStoreList?: Array<ObjectStore>;
  postRoute: string;
  sendRefetch?: any;
}

interface RenderMap {
  [key: string]: JSX.Element | undefined;
  idle?: JSX.Element;
  processing?: JSX.Element;
  error?: JSX.Element;
  success?: JSX.Element;
}


const CreateModal: FC<Props> = ({
  handleClose,
  isVisible,
  modalType,
  objectStoreList,
  postRoute,
  sendRefetch,
}:Props )=> {
  const [state, send] = useMachine(ModalMachine);
  // refs
  const inputNameRef = useRef(null);
  const inputDescriptionRef = useRef(null);
  // state
  const [ name, setName ] = useState('');
  const [ description, setDescription ] = useState('');
  const [ bucketName, setBucketName ] = useState(null);
  const [ objectStoreType, setObjectStoreType ] = useState(null);
  const [ errorMessage, setErrorMessage ] = useState('');
  // vars
  const modalSize = (modalType === 'namespace') ? 'longer' : 'medium';



  // set default
  useEffect(() => {
    if(objectStoreList && objectStoreList[0]) {
      setObjectStoreType(objectStoreList[0]);
    }
  }, [objectStoreList]);

  /**
  * Metod handles close for model and checks if error state is present
  * @param {}
  * @return {void}
  * @fires {#send}
  * @fires {#handleClose}
  */
  const manageClose = () => {
    if (state.value === 'error') {
      send("TRY_AGAIN");
    }
    handleClose();
  }

  /**
  * Metod updates the name value;
  * @param {KeyboardEvent} event
  * @return {void}
  * @fires {#setName}
  */
  const handleNameEvent = (event:KeyboardEvent) => {
    const element = event.currentTarget as HTMLInputElement
    const value = element.value

    setName(value);
  }


  /**
  * Metod updates the name value;
  * @param {KeyboardEvent} event
  * @return {void}
  * @fires {#setName}
  */
  const handleBucketNameEvent = (event:KeyboardEvent) => {
    const element = event.currentTarget as HTMLInputElement
    const value = element.value

    setBucketName(value);
  }


  /**
  * Metod updates the name value;
  * @param {KeyboardEvent} event
  * @return {void}
  * @fires {#setName}
  */
  const handleObjectChangeEvent = (value:ObjectStore) => {

    setObjectStoreType(value);
  }

  /**
  * Metod updates the description value;
  * @param {KeyboardEvent} event
  * @return {void}
  * @fires {#setDescription}
  */
  const handleDescriptionEvent = (event:KeyboardEvent) => {
    const element = event.currentTarget as HTMLInputElement
    const value = element.value

    setDescription(value);
  }

  /**
  * Metod submits data to create a new name
  * @param {}
  * @return {void}
  * @fires {#post}
  * @fires {#handleClose}
  * @fires {#updateFetchId}
  */
  const handleCreate = () => {
    send("SUBMIT");

    if (
      (name.length < 3)
      || ((bucketName && bucketName.length < 3) && (modalType !== 'namespace'))
    ) {
      return;
    }

    const postBody = {
      name,
      description,
      bucket_name: modalType === 'namespace' ? bucketName : undefined,
      object_store_name: modalType === 'namespace' ? objectStoreType.name : undefined
    }

    post(postRoute, postBody).then((response) => {
      return response.json();
    }).then((data) => {
      send("SUCCESS");
      setTimeout(() => {
        handleClose();
        send("RESET");
        sendRefetch();
      }, 1000)

    }).catch((error) => {
      send("ERROR");
      setErrorMessage(error)
    });

  }

  if (!isVisible) {
    return null;
  }

  const stateValue: string = typeof state.value === 'string' ? state.value : 'idle';


  const renderMap: RenderMap = {
    idle: (
      <div className={`CreateModal CreateModal--${modalType}`}>
        <div className="CreateModal__form">
          <p>{`Create a new ${modalType}`}</p>
          <InputText
            css=""
            isRequired
            inputRef={inputNameRef}
            label={modalType}
            updateValue={handleNameEvent}
          />

          <InputText
            css=""
            inputRef={inputDescriptionRef}
            label="Description"
            updateValue={handleDescriptionEvent}
          />


          { (modalType === 'namespace') &&
            <NamespaceFields
              handleBucketNameEvent={handleBucketNameEvent}
              handleObjectChangeEvent={handleObjectChangeEvent}
              objectStoreList={objectStoreList}
            />
          }
        </div>




        <div className="CreateModal__buttons">
          <TernaryButton
            click={handleClose}
            text="Cancel"
          />
          <PrimaryButton
            click={handleCreate}
            disabled={name.length < 3}
            text={`Create ${modalType}`}
          />
        </div>
      </div>
    ),
    processing: (
      <Processing
        modalType={modalType}
        name={name}
      />
    ),
    success: (
      <Success
        modalType={modalType}
        name={name}
      />
    ),
    error: (
      <Error
        errorMessage={errorMessage}
        name={name}
        send={send}
      />
    )
  }

  return (
    <Modal
      handleClose={manageClose}
      header={`Create ${modalType}`}
      icon=""
      overflow=""
      size={modalSize}
      subheader=""
    >
      {renderMap[stateValue]}
    </Modal>
  )
}

export default CreateModal;
