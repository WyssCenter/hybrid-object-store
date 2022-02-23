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
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faInfo } from '@fortawesome/free-solid-svg-icons'
// environment
import { get, put } from 'Environment/createEnvironment';
// components
import { PrimaryButton } from 'Components/button/index';
import { InputText } from 'Components/form/text/index';
// context
import AppContext from 'Src/AppContext';
import GroupContext from '../../GroupContext';
// css
import './AddUser.scss';

const AddUser: FC = () => {
  // context
  const { send, groupname } = useContext(GroupContext);
  const { user } = useContext(AppContext);

  // refs
  const userInputRef = useRef(null);
  // states
  const [errorMessage, setErrorMessage] = useState(null);
  const [username, setUsername] = useState(null);
  const [nameList, setNameList] = useState(null);
  const [buttonDisabled, setButtonDisabled] = useState(false);
  // vars
  const usersGroups = user.profile.groups.split(',');
  const canMemberEdit = (user.profile.role === 'admin') ||  (user.profile.role === 'privileged')
    ? true
    : false;
  const editDisabled = (username === null) || (username === '') || !canMemberEdit;

  /**
  * Method generates a personal access token
  * @param {}
  * @return {void}
  * @calls {state#setPat}
  * @calls {state#setErrorMessage}
  */
  const handleAddUser = useCallback((name: string) => {
    put(
        `group/${groupname}/user/${name}`,
        {},
        true
    ).then((response) => {
      return response.json();
    }).then((data) => {
      if (data && data.error) {
        setErrorMessage(data.error);
        return;
      }
      setErrorMessage(null);
      userInputRef.current.value = '';
      send('REFETCH');
    }).catch((error) => {
      setErrorMessage(error.toString());
    });
  }, [username, groupname, send]);

  const handleWrapper = useCallback((name: string) => {
    if (username.indexOf('@') > -1) {
      get(`usernames?email=${username}`, true)
      .then((res) => res.json())
      .then((data) => {
        if (!data.usernames) {
          setErrorMessage('No users with that e-mail exist.');
          setTimeout(() => {
            setErrorMessage('');
          }, 5000);
        } else {
          setButtonDisabled(true);
          setNameList(data.usernames);
        }
      })
      return;
    }
    handleAddUser(name);
  }, [handleAddUser]);

  /**
  * Method updates desctipion
  * @param {KeyboardEvent} event
  * @return {void}
  * @calls {state#setDescription}
  */
  const updateUsername = useCallback((event: KeyboardEvent) => {
    const element = event.currentTarget as HTMLInputElement;
    const newUsername = element.value;
    if (!newUsername) {
      setButtonDisabled(false);
      setNameList(null);
    }

    if((event.key === 'Enter') && (newUsername.length > 0)){
      handleWrapper(username);
    }

    setUsername(newUsername);

  }, [username, setUsername, handleWrapper]);


  return (
    <>
    <tr className="AddUser">

      <td colSpan={3}>
        <div className="AddUser__form">
          <InputText
            css="column-1-span-4"
            disabled={!canMemberEdit}
            inputRef={userInputRef}
            label="Add a user by username or email address"
            placeholder="user-1"
            updateValue={updateUsername}
          />
          {
            (nameList !== null) && (
            <div className="AddUser__auto">
              {
                nameList.map((key: any) => {
                  return (
                    <button
                      className="AddUser__auto--button"
                      key={key}
                      onClick={() => {
                        handleAddUser(key);
                      }}
                    >
                      {key}
                    </button>
                  )
                })
              }
            </div>
            )
          }
        </div>
      </td>
      <td>
        <div className="flex justify--center">
          <PrimaryButton
            click={() => handleWrapper(username)}
            disabled={editDisabled || buttonDisabled}
            text="Add User"
          />
        </div>
      </td>

    </tr>
    { errorMessage &&
      <tr>
        <td colSpan={4}>
        <p className="text-center error">Error adding user; {errorMessage}. Make sure you have entered the username correctly and try again.</p>

      </td>
    </tr>
    }

    { !canMemberEdit &&
      <tr className="AddUser__container">
        <td colSpan={6}>
          <div className="AddUser__warning">
            <FontAwesomeIcon
              color="rgb(236, 128, 11)"
              icon={faInfo}
              size="lg"
            />
          </div>
          <h5 className="AddUser__h5--warning">
            Privileged users must be a member of a group to have write access.
          </h5>
        </td>
      </tr>
    }
    </>
  )
}

export default AddUser;
