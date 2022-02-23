// vendor
import React,
{
  FC,
  useState,
  useContext,
} from 'react';
import {
  useParams,
} from "react-router-dom";
import ReactTooltip from 'react-tooltip';
// environment
import { put } from 'Environment/createEnvironment';
// context
import AppContext from 'Src/AppContext';
// components
import PermissionDropdown from 'Shared/permission/PermissionDropdown';
import Delete from './delete/Delete';
// css
import './PermissionsItem.scss';

interface Group {
  group_name: string;
}

interface Permission {
  name: string;
}

interface Item {
  group?: Group;
  permission: Permission;
}

interface Props {
  item: Item;
  sectionType: string;
  owner: string;
}

interface ParamTypes {
  datasetName: string;
  namespace: string;
}


const PermissionsItem: FC<Props> = ({
  item,
  sectionType,
  owner,
}: Props) => {
  // context
  const { user } = useContext(AppContext);
  // params
  const { datasetName, namespace } = useParams<ParamTypes>();
  const [ errorMessage, setErrorMessage ] = useState(null);
  // vars
  // strip hoss-deault-group from username
  const name = item.group.group_name.indexOf('-hoss-default-group') > -1
    ? item.group.group_name.replace('-hoss-default-group', '')
    : item.group.group_name;

  /**
  * Method calls put method to set permisstions for a user or group
  * @param { String } permission
  * @param { String } name
  * @return {Void}
  * @calls {environment#put}
  * @calls {state#setErrorMessage}
  */
  const updatePermissions = (permission: Permission, name: string, errorCallback: () => void) => {
    let permName = 'r';
    if (permission.name === 'Read & Write') {
      permName = 'rw';
    }
    if (permission.name === 'Read Only') {
      permName = 'r';
    }
    put(`namespace/${namespace}/dataset/${datasetName}/${sectionType}/${name}/access/${permName}`)
    .then((response) => {
      if(response.ok) {
        return;
      }
      return response.json();
    })
    .then((data: any) => {
      if (data && data.error) {
        errorCallback();
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
  };

  const authorized = user.profile.role !== 'user';

  return (
    <>
      <tr>
        <td className="PermissionsItem__name">{name}</td>
        <td>
          <PermissionDropdown
            name={name}
            permission={item.permission}
            updatePermissions={updatePermissions}
          />
        </td>
        {
          authorized && (
            <>
              <td></td>
              <td>
                <div
                  data-tip="You cannot remove the owner from the dataset."
                  data-tip-disable={!(owner === name && sectionType === 'user')}
                >
                  <Delete
                    datasetName={datasetName}
                    namespace={namespace}
                    name={name}
                    sectionType={sectionType}
                    setErrorMessage={setErrorMessage}
                    isOwner={owner === name && sectionType === 'user'}
                  />
                  <ReactTooltip
                    place="bottom"
                    effect="solid"
                  />
                </div>
              </td>
            </>
          )
        }
      </tr>


      { errorMessage && (
        <p className="PermissionsItem__error error">{errorMessage}</p>
      )}
    </>

  );
}


export default PermissionsItem;
