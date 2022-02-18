// vendor
import React,
{
  FC,
  useCallback,
  useContext,
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
// context
import GroupsContext from '../../GroupsContext';
// css
import './CreateGroupModal.scss';


interface Props {
  isVisible: boolean;
  handleClose: () => void;
}


const CreateGroup: FC<Props> = ({
  isVisible,
  handleClose,
}: Props) => {
  // context
  const { send } = useContext(GroupsContext);

  // refs
  const groupInputRef = useRef(null);
  const descriptionInputRef = useRef(null);
  // states
  const [errorMessage, setErrorMessage] = useState(null);
  const [groupName, setGroupName] = useState(null);
  const [description, setDescription] = useState(null);
  // vars
  const createDisabled = (groupName === null) || (groupName === '')
    || (description === null) || (description === '');

  // functions
  /**
  * Method generates a personal access token
  * @param {}
  * @return {void}
  * @calls {state#setPat}
  * @calls {state#setErrorMessage}
  */
  const handleAddGroup = useCallback(() => {
    const body = { description, name: groupName };
    post('group/', body, true).then((response) => {
      return response.json();
    }).then((data) => {
      if (data.error) {
        setErrorMessage(data.error);
        return;
      }
      setErrorMessage(null);
      groupInputRef.current.value = '';
      descriptionInputRef.current.value = '';
      send('REFETCH');
    }).catch((error) => {
      setErrorMessage(error.toString());
    });
  }, [description, groupName, send]);

  /**
  * Method updates desctipion
  * @param {KeyboardEvent} event
  * @return {void}
  * @calls {state#setDescription}
  */
  const updateDescrption = useCallback((event: KeyboardEvent) => {
    const element = event.currentTarget as HTMLInputElement;

    const newDescription = element.value;


    if((event.key === 'Enter') && (newDescription.length > 0) && (groupName.length > 0)){
      handleAddGroup();
    }

    setDescription(newDescription);

  }, [setDescription, groupName, handleAddGroup]);

  /**
  * Method updates desctipion
  * @param {KeyboardEvent} event
  * @return {void}
  * @calls {state#setDescription}
  */
  const updateGroupName = useCallback((event: KeyboardEvent) => {
    const element = event.currentTarget as HTMLInputElement;

    const newGroupName = element.value;

    if((event.key === 'Enter') && (newGroupName.length > 0) && (description.length > 0)){
      handleAddGroup();
    }

    setGroupName(newGroupName);

  }, [setGroupName, handleAddGroup, description]);

  if(!isVisible) {
    return null;
  }


  return (
    <Modal
      handleClose={handleClose}
      header="Create Group"
      icon=""
      overflow="visible"
      subheader=""
    >
      <p>Groups are collections of users, create a group and add users to it.</p>
      <section className="CreateGroupModal">
          <InputText
            css="column-1-span-4"
            inputRef={groupInputRef}
            label="Group Name"
            placeholder="my-group"
            updateValue={updateGroupName}
          />

          <InputText
            css="column-1-span-4"
            inputRef={descriptionInputRef}
            label="Description of group"
            placeholder="This is my group"
            updateValue={updateDescrption}
          />
        <div className="CreateGroupModal__buttons text-right">
          <PrimaryButton
            click={handleAddGroup}
            disabled={createDisabled}
            text="Create Group"
          />
        </div>
      </section>

      { errorMessage &&
        <section>
          <p className="error">{errorMessage}</p>
        </section>
      }
    </Modal>
  )
}

export default CreateGroup;
