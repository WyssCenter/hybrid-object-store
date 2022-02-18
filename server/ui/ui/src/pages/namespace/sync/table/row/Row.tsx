// vendor
import React,
  {
    FC,
    MouseEvent,
    useCallback,
    useContext,
    useRef,
    useState,
  } from 'react';
import { faTrash } from '@fortawesome/free-solid-svg-icons';
import { useParams } from 'react-router-dom';
// context
import NamespaceContext from '../../../NamespaceContext';
// environment
import { del } from 'Environment/createEnvironment';
// components
import { IconButton } from 'Components/button/index';
import { TooltipConfirm } from 'Components/tooltip/index';

interface ParamTypes {
  namespace: string;
}

interface Target {
  sync_type: string;
  target_core_service: string;
  target_namespace: string;
}

interface Props {
  target: Target;
}



const Row:FC<Props> = ({ target }: Props) => {
  // context
  const { send } = useContext(NamespaceContext)
  // params
  const { namespace } = useParams<ParamTypes>();
  // refs
  const tooltipRef = useRef(null);
  // state
  const [isTooltipVisible, setIsTooltipVisible] = useState(false);

  // events
  /**
  * Method updates tooltip visibility
  * @param {}
  */
  const shopwIsTooltipVisible = useCallback(() => {
    setIsTooltipVisible(true);
  }, []);
  /**
  * Method updates tooltip visibility
  * @param {}
  */
  const hideIsTooltipVisible = () => {
    setIsTooltipVisible(false);
  }

  /**
  * Method removes sync item from table
  * @param {event} Object
  */
  const removeSyncItem:any = useCallback((event: MouseEvent) => {
    del(`namespace/${namespace}/sync`, target).then((response) => {

      if(response.ok) {
        setIsTooltipVisible(false);
        send("REFETCH")
      }
    }).catch((error) => {
      console.log(error);
    });
  }, []);
  return (
    <tr>

      <td colSpan={2}>{target.target_core_service}</td>
      <td>{target.target_namespace}</td>
      <td>{target.sync_type === 'simplex' ? '1 Way' : '2 Way' }</td>

      <td className="SyncTable__actions">
        <div ref={tooltipRef}>
          <IconButton
            click={shopwIsTooltipVisible}
            color="white"
            icon={faTrash}
          />

          <TooltipConfirm
            cancel={hideIsTooltipVisible}
            confirm={removeSyncItem}
            isVisible={isTooltipVisible}
            tooltipRef={tooltipRef}
            text="Are you sure you want to remove this sync service?"
          />
        </div>
      </td>
    </tr>
  )
}

export default Row;
