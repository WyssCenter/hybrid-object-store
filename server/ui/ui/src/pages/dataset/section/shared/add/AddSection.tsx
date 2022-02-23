// vendor
import React,
{
  FC,
  KeyboardEvent,
  useCallback,
  useContext,
  useRef,
  useState,
} from 'react';
import {
  useParams,
} from "react-router-dom";
// environment
import { get, put } from 'Environment/createEnvironment';
// components
import Button from 'Components/button/index';
import UserInput from './input/UserInput';
import PermissionDropdown from 'Shared/permission/PermissionDropdown';
// context
import DatasetContext from '../../../DatasetContext';
// css
import './AddSection.scss';

interface Props {
  sectionType: string;
}

interface Permission {
  name:string;
}

interface ParamTypes {
  datasetName: string;
  namespace: string;
}

const AddSection: FC<Props> = ({ sectionType }:Props) => {
  // refs
  const inputRef = useRef(null);
  // context
  const { send } = useContext(DatasetContext);
  // params
  const { datasetName, namespace } = useParams<ParamTypes>();
  // state
  const [permissions, setPermissions ] = useState({name: 'r'});
  const [name, setName] = useState('');
  const [errorMessage, setErrorMessage] = useState('');
  const [nameList, setNameList] = useState(null);
  const [buttonDisabled, setButtonDisabled] = useState(false);

  // functions

  /**
  * Method updates permissions in state
  * @param {string} newPermissions
  * @return {void}
  * @calls {this#setPermissions}
  */
  const updatePermissions = useCallback((newPermissions) => {
    setPermissions(newPermissions)
  }, [setPermissions]);


  const checkName = (evt: Event) => {
    const target = evt.target as HTMLInputElement;
    if (!target.value) {
      setButtonDisabled(false);
      setNameList(null);
    }
  }

  const submitFetch = (name: string) => {
    get(`${sectionType}/${name}`, true)
    .then((response: any) => {
      if (response.status === 200) {
        inputRef.current.value = '';
        let permName = 'r';
        if (permissions.name === 'Read & Write') {
          permName = 'rw';
        }
        if (permissions.name === 'Read Only') {
          permName = 'r';
        }
        put(`namespace/${namespace}/dataset/${datasetName}/${sectionType}/${name}/access/${permName}`)
        .then((response) => {
          if (response.ok) {
            send('REFETCH');
            setPermissions({name: 'r'});
            setName('');
            setErrorMessage('');
            return;
          }
          return response.json()
        })
        .then((data: any) => {
          if(data && data.error) {
            setErrorMessage(data.error);
            setTimeout(() => {
              setErrorMessage('');
            }, 5000);
          }
        })
        .catch((error) => {
          const newErrorMessage = error.toString ? error.toString() : error;
          setErrorMessage(newErrorMessage);
        })
      } else {
          setErrorMessage(`${sectionType} with that name does not exist.`);
          setTimeout(() => {
            setErrorMessage('');
          }, 5000);
      }
    })
  };
  /**
  * Method resets input and requets via put add a user or group
  * @param {}
  * @return {void}
  * @calls {environment#put}
  * @calls {Dataset#fetchDatasetData}
  * @calls {this#setPermissions}
  * @calls {this#setName}
  */
  const submit = useCallback(() => {
    if (sectionType === 'user' && name.indexOf('@') > -1) {
      get(`usernames?email=${name}`, true)
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
    submitFetch(name);
  }, [datasetName, name, namespace, permissions, sectionType, send]);



    /**
    * Method handles keyUp event and updates name value
    * @param {KeyboardEvent} event
    * @return {void}
    * @calls {state#setName}
    */
    const updateName = (event: KeyboardEvent) => {
      const element = event.currentTarget as HTMLInputElement
      const value = element.value

      if (event.key === "Enter") {
        submit();
      } else {

        setName(value);
      }
    }

  return (
    <>
      <div className="AddSection">
        <div className="relative">
            <UserInput
              inputRef={inputRef}
              permissionType={sectionType}
              updateName={updateName}
              checkName={checkName}
            />
              {
                (nameList !== null) && (
                <div className="AddSection__auto">
                  {
                    nameList.map((key: any) => {
                      return (
                        <button
                          className="AddSection__auto--button"
                          key={key}
                          onClick={() => {
                            submitFetch(key);
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
            <div className="AddSection__dropdown flex">
              <div>
                <PermissionDropdown
                  name=""
                  permission={permissions}
                  updatePermissions={updatePermissions}
                />
              </div>
              <div className="AddSection__column--text-right">
                <Button
                  click={submit}
                  disabled={name.length === 0 || buttonDisabled}
                  text={`Add ${sectionType}`}
                />
              </div>
            </div>


      </div>
      { errorMessage && (
        <div className="AddSection">
        <div>
          <p className="error">{errorMessage}</p>
          </div>
        </div>
      )}
    </>
  )
}

export default AddSection;
