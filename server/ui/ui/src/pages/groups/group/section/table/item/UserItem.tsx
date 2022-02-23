// vendor
import React,
{
  FC,
  useCallback,
  useContext,
  useRef,
  useState
} from 'react';
import { faTrash } from '@fortawesome/free-solid-svg-icons';
// environment
import { del } from 'Environment/createEnvironment';
// components
import { IconButton } from 'Components/button/index';
import { TooltipConfirm } from 'Components/tooltip/index';
// context
import AppContext from 'Src/AppContext';
import GroupContext from '../../../GroupContext'
// css
import './UserItem.scss';

interface User {
  full_name: string;
  role: string;
  username: string;
}

interface Props {
  userItem: User;
}

const UserItem: FC<Props> = ({ userItem }: Props) => {
  // context
  const { user } = useContext(AppContext);
  const { groupname, send } = useContext(GroupContext);
  // refs
  const tooltipRef = useRef(null);
  // state
  const [isTooltipVisible, setIsTooltipVisible] = useState(null);
  const [errorMessage, setErrorMessage] = useState(null);
  // vars
  const usersGroups = user.profile.groups.split(',');
  const canMemberEdit = (user.profile.role === 'admin') || (user.profile.role === 'privileged')
    ? true
    : false;

  /**
  * Method calls delete to remove
  * @param {}
  * @return {void}
  * @calls {environment#del}
  * @calls {machine#send}
  * @calls {state#setErrorMessage}
  */
  const removeUser = useCallback(() => {
    del(`group/${groupname}/user/${userItem.username}`, {}, true).then((response) => {
      if(response.ok) {
        send('REFETCH');
        return;
      } else {
        return response.json();
      }
    })
    .then((data: any) => {
      if (data && data.error) {
        setIsTooltipVisible(false)
        setErrorMessage(data.error);
        setTimeout(() => {
          setErrorMessage('');
        }, 5000);
      }
    })
    .catch((error) => {
      setErrorMessage(error.toString());
    })
  }, [send, groupname, userItem])


  /**
  * Method closes confirmation tooltip
  * @param {}
  * @return {void}}
  * @calls {state#setIsTooltipVisible}
  */
  const hideTooltipVisible = () => {
    setIsTooltipVisible(false);
  }

  return (
    <>
      <tr className="UserItem">
        <td>
            {userItem.username}
        </td>
        <td>
          {userItem.full_name}
        </td>
        <td>
         {userItem.role}
        </td>
        { (user.profile.role !== 'user') && (
          <td className="UserItem__td--actions">
            <div
              className="UserItem__actions"
              ref={tooltipRef}
            >
              <IconButton
               click={() => { setIsTooltipVisible(true)}}
               disabled={!canMemberEdit}
               icon={faTrash}
               color="white"
              />
              <TooltipConfirm
                confirm={removeUser}
                cancel={hideTooltipVisible}
                isVisible={isTooltipVisible}
                text="Are you sure? The user will lose access to all datasets that this group is granted permissions to."
                tooltipRef={tooltipRef}
              />
            </div>
          </td>
        )}
      </tr>

      {errorMessage && (
        <tr className="text-center">
          <td colSpan={4}>
            <p className="error">{errorMessage}</p>
          </td>
        </tr>
      )}
    </>
  )
}


export default UserItem;
