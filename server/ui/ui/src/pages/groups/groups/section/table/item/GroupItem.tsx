// vendor
import React,
{
  FC,
  useCallback,
  useContext,
  useRef,
  useState,
} from 'react';
import { faTrash } from '@fortawesome/free-solid-svg-icons';
import { Link } from 'react-router-dom'
// Environment
import { del } from 'Environment/createEnvironment';
// components
import { IconButton } from 'Components/button/index';
import { TooltipConfirm } from 'Components/tooltip/index';
// context
import GroupsContext from '../../../GroupsContext';
// css
import './GroupItem.scss';

interface Group {
  description: string;
  group_name: string;
}

interface Membership {
  group: Group;
}

interface Props {
  membership: Membership;
  authorized: boolean;
}

const GroupItem: FC<Props> = ({ membership, authorized }: Props) => {

  // refs
  const tooltipRef = useRef(null);
  // state
  const [isTooltipVisible, setIsTooltipVisible] = useState(null);
  const [errorMessage, setErrorMessage] = useState(null);
  const { send } = useContext(GroupsContext);

  /**
  * Method calls delete to remove group
  * @param {}
  * @return {void}
  * @calls {environment#del}
  * @calls {machine#send}
  * @calls {state#setErrorMessage}
  */
  const removeGroup = useCallback(() => {
    del(
      `group/${membership.group.group_name}`,
      del,
      true).then((response) => {
      if(response.ok) {
        send('REFETCH');
      } else {
        setErrorMessage("Something went wrong, try again. If this continues contact your support.");
      }
    }).catch((error) => {
      setErrorMessage(error.toString());
    })
  }, [send, membership])


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
      <tr className="GroupItem">
        <td>
          <Link
            to={`/groups/${membership.group.group_name}`}
          >
            {membership.group.group_name}
          </Link>
        </td>
        <td>{membership.group.description}</td>
        {
          authorized && (
            <td className="GroupItem__td--actions">
              <div
                className="GroupItem__actions"
                ref={tooltipRef}
              >
                <IconButton
                click={() => { setIsTooltipVisible(true)}}
                icon={faTrash}
                color="white"
                />
                <TooltipConfirm
                  confirm={removeGroup}
                  cancel={hideTooltipVisible}
                  isVisible={isTooltipVisible}
                  text="Are you sure? If you haven't copied and saved this token you will not be able use it."
                  tooltipRef={tooltipRef}
                />
              </div>
            </td>
          )
        }
      </tr>

      { errorMessage && (
        <tr className="text-center">
          <td colSpan={4}>
            <p className="error">Failed to delete group; {errorMessage}</p>
          </td>
        </tr>
      )}

    </>
  )
}


export default GroupItem;
