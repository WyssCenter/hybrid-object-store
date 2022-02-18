// vendor
import React,
{
  FC,
  useCallback,
  useRef,
  useState,
} from 'react';
import { faTrash } from '@fortawesome/free-solid-svg-icons'
// environment
import { del } from 'Environment/createEnvironment'
// components
import { IconButton } from 'Components/button/index';
import { TooltipConfirm } from 'Components/tooltip/index';
// css
import './PatItem.scss';

interface Pat {
  description: string;
  id: number;
}

interface ReponseError {
  error: string;
}

interface Props {
  pat: Pat;
  send: any;
}

const PatItem: FC<Props> = ({ pat, send }: Props) => {
  // vars
  const { description, id } = pat;
  // ref
  const tooltipRef = useRef<HTMLTableDataCellElement>(null);
  // state
  const [errorMessage, setErrorMessage] = useState(null);
  const [isTooltipVisible, setTooltipVisible] = useState(false);


  /**
  * Method deletes pat
  * @param {}
  * @return {void}
  * @call {del}
  * @call {updateFetchId}
  */
  const deletePat = useCallback(() => {
    del(
      `pat/${id}`,
      {},
      true,
    )
      .then((response) => response.text())
      .then((data) => {
        if (data) {
          const dataObj: ReponseError = JSON.parse(data);
          setErrorMessage(dataObj.error);
          return;
        }
        send('REFETCH');
      }).catch((error) => {
        const newErrorMessage = error.toString
          ? error.toString()
          : 'error deleting pat';
        setErrorMessage(newErrorMessage);
      })
  }, [send, setErrorMessage, id]);

  /**
  * Method sets cancels delete
  * @param {}
  * @return {void}
  * @call {state#setTooltipVisible}
  */
  const cancel = useCallback(() => {
    setTooltipVisible(false);
  }, []);

  /**
  * Method shows tooltip to confirm delete
  * @param {}
  * @return {void}
  * @call {state#setTooltipVisible}
  */
  const showTooltip = useCallback(() => {
    setTooltipVisible(true);
  }, []);

  return (
    <>
      <tr className="PatItem">
        <td className="column-1-span-5">
          <p>{description}</p>
        </td>
        <td ref={tooltipRef}>
          <IconButton
            click={showTooltip}
            icon={faTrash}
            color="white"
          />

          <TooltipConfirm
            confirm={deletePat}
            cancel={cancel}
            isVisible={isTooltipVisible}
            text="Are you sure? Any applications or scripts using this token will no longer work. This cannot be undone."
            tooltipRef={tooltipRef}
          />
        </td>
      </tr>

      { errorMessage &&
        <tr className="PatItem PatItem--error">
          <td
            aria-colspan={2}
            colSpan={2}
          >
            <p>Error revoking token</p>
          </td>
        </tr>
      }
    </>
  )
}


export default PatItem;
