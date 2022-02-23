// vendor
import React,
{
  FC,
  useCallback,
  useState,
  useRef,
  KeyboardEvent,
} from 'react';
// environment
import { post } from 'Environment/createEnvironment';
// components
import { PrimaryButton } from 'Components/button/index';
import { InputText } from 'Components/form/text/index';
import Modal from 'Components/modal/Modal';
import NewPat from './new/NewPat';
// css
import './CreatePatModal.scss';

interface Props {
  handleClose: () => void;
  isVisible: boolean;
  send: any;
}

const CreatePatModal: FC<Props> = ({
  handleClose,
  isVisible,
  send,
}: Props) => {
  // refs
  const inputRef = useRef(null);
  // states
  const [pat, setPat] = useState(null);
  const [errorMessage, setErrorMessage] = useState(null);
  const [description, setDescription] = useState('this is a description');

  // functions
  /**
  * Method updates desctipion
  * @param {KeyboardEvent} event
  * @return {void}
  * @calls {state#setDescription}
  */
  const updateDescrption = useCallback((event: KeyboardEvent) => {
    const element = event.currentTarget as HTMLInputElement;

    const newDescription = element.value;

    setDescription(newDescription);

  }, [setDescription]);

  /**
  * Method generates a personal access token
  * @param {}
  * @return {void}
  * @calls {state#setPat}
  * @calls {state#setErrorMessage}
  */
  const handleGeneratePat = useCallback(() => {
    const body = {description};
    post('pat', body, true).then((response) => {
      return response.json();
    }).then((data) => {
      if (data && data.error) {
        setErrorMessage(data.error);
        return;
      }
      setErrorMessage(null);
      inputRef.current.value = '';
      setPat(data);
    }).catch((error) => {
      setErrorMessage(error.toString());
    });
  }, [description, setPat, setErrorMessage]);

  /**
  * Method removes generated personal access token from state a refreshes the ui
  * @param {}
  * @return {void}
  * @calls {state#setPat}
  */
  const dismissPat = useCallback(() => {
    setPat(null);
    send('REFETCH');
  }, [setPat, send]);

  if (!isVisible) {
    return null;
  }

  return (
    <Modal
      handleClose={() => {
        handleClose();
        if (pat) {
          dismissPat();
        }
      }}
      header="Create Personal Access Token"
    >
      <section className="CreatePatModal">
        <div className="CreatePatModal__form">

          { (pat === null) && (
            <>
              <p>Personal access tokens are access tokens that do not expire. These are often ideal for programmatic access. PATs are exchanged for time-limited access tokens that can be used to interact with the Hoss API.</p>
              <div className="CreatePatModal__flex">
                <InputText
                  flexParent
                  inputRef={inputRef}
                  label="What is this token for?"
                  placeholder="A brief token description"
                  updateValue={updateDescrption}
                />
                <div className="CreatePatModal__buttons">
                  <PrimaryButton
                    click={handleGeneratePat}
                    text="Generate Token"
                  />
                </div>
              </div>
            </>
          )}
          { pat &&
            <NewPat
              dismissPat={dismissPat}
              pat={pat}
            />
          }
        </div>


        { errorMessage &&
          <p className="error">{errorMessage}</p>
        }

      </section>
    </Modal>
  )
}

export default CreatePatModal;
