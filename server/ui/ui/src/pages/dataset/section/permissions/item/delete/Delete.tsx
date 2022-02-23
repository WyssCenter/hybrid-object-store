// vendor
import React,
{
  FC,
  useContext,
  useRef,
  useState,
} from 'react';
import { faTrash } from '@fortawesome/free-solid-svg-icons';
// environment
import { del } from 'Environment/createEnvironment';
// components
import { TooltipConfirm } from 'Components/tooltip/index';
import { IconButton } from 'Components/button/index';
// context
import DatasetContext from '../../../../DatasetContext';
// css
import './Delete.scss';


interface Props {
  datasetName: string;
  name: string;
  namespace: string;
  sectionType: string;
  setErrorMessage: any;
  isOwner: boolean;
}

const AddPermission: FC<Props> = ({
  datasetName,
  namespace,
  name,
  sectionType,
  setErrorMessage,
  isOwner,
}:Props) => {
  // ref
  const tooltipRef = useRef(null);
  // state
  const [isTooltipVisible, setIsTooltipVisible] = useState(false);
  // context
  const { send } = useContext(DatasetContext);


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
    del(`namespace/${namespace}/dataset/${datasetName}/${sectionType}/${name}`)
    .then((response) => {
      if (response.ok) {
        send('REFETCH');
        return;
      }
      return response.json();
    })
    .then((data: any) => {
      if (data && data.error) {
        setErrorMessage(data.error);
        setTimeout(() => {
          setErrorMessage('');
        }, 5000);
        setIsTooltipVisible(false);
      }
    })
    .catch((error) => {
      setErrorMessage(`Error: ${sectionType} may have not been deleted correctly, refresh to confirm`);
    });
  }

  return (
    <div className="Delete relative">
      <IconButton
        click={() => { setIsTooltipVisible(true)}}
        color="white"
        disabled={isOwner}
        icon={faTrash}
      />
      <TooltipConfirm
        confirm={deleteItem}
        cancel={hideTooltipVisible}
        isVisible={isTooltipVisible}
        text={`Are you sure? This ${sectionType} will lose access to this dataset.`}
        tooltipRef={tooltipRef}
      />

    </div>
  );
}


export default AddPermission;
