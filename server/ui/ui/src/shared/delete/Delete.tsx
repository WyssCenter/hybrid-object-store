// vendor
import React,
{
  FC,
  useContext,
  useRef,
  useState,
  Context,
} from 'react';
import { faTrash } from '@fortawesome/free-solid-svg-icons';
// environment
import { del } from 'Environment/createEnvironment';
// components
import { Tooltip, TooltipConfirm } from 'Components/tooltip/index';
import { IconButton } from 'Components/button/index';
// css
import './Delete.scss';


interface Props {
  context: Context<{ send: any; }> ;
  route: string;
  sectionType: string;
  tooltipText: string;
}

const AddPermission: FC<Props> = ({
  context,
  route,
  sectionType,
  tooltipText,
}:Props) => {
  // ref
  const tooltipRef = useRef(null);
  // state
  const [isTooltipVisible, setIsTooltipVisible] = useState(false);
  const [errorMessage, setErrorMessage] = useState(null);
  // context
  const { send } = useContext(context);


  /**
  * Method hides tooltip
  * @param {}
  * @return {Void}
  * @calls {state#setIsErrorTooltipVisible}
  *
  */
  const dismissErrorMessage = () => {
    setErrorMessage(null);
  }

  /**
  * Method hides tooltip
  * @param {}
  * @return {Void}
  * @calls {state#setIsTooltipVisible}
  *
  */
  const hideTooltipVisible = () => {
    setIsTooltipVisible(false);
  }


  /**
  * Method handles click on delete button
  * @param {MouseEvent} event
  * @return {void}
  * @calls {event#stopPropagation}
  * @calls {environment#del}
  * @calls {fetchDatasetData}
  * @calls {state#setErrorMessage}
  */
  const deleteItem = () => {
    del(route).then((response: Response) => {
      send('REFETCH');
    }).catch((error: Error) => {
      setErrorMessage(`Error: ${sectionType} may have not been deleted correctly, refresh to confirm`);
    });
  }

  return (
    <div className="Delete relative">
      <IconButton
        click={() => { setIsTooltipVisible(true)}}
        color="white"
        disabled={false}
        icon={faTrash}
      />
      <TooltipConfirm
        confirm={deleteItem}
        cancel={hideTooltipVisible}
        isVisible={isTooltipVisible}
        text={tooltipText}
        tooltipRef={tooltipRef}
      />

      <Tooltip
        isVisible={errorMessage !== null}
        updateIsVisible={dismissErrorMessage}
        text={errorMessage}
      />
    </div>
  );
}


export default AddPermission;
