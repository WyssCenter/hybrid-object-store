// vendor
import React, { FC, useCallback, useRef, useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCopy } from '@fortawesome/free-solid-svg-icons';
// components
import { Tooltip } from 'Components/tooltip/index';

// css
import './CopyText.scss';


interface Props {
  text: string;
}

const CopyText: FC<Props> = ({text}: Props) => {
  const inputRef = useRef(null);
  const [isTooltipVisible, setIsTooltipVisible] = useState(false);

  // refs
  const copyRef = useRef(null);

  /**
  * Method copies input value to clipboard
  * @param {}
  * @return {void}
  * @calls {document#execCommand}
  */
  const copyToClipboard = useCallback(() => {
    inputRef.current.select();

    document.execCommand('Copy');
    setIsTooltipVisible(true);

  }, []);

  return (
    <div className="CopyText">
      <input
        className="CopyText__input"
        readOnly
        ref={inputRef}
        value={text}
      />
      <button
        ref={copyRef}
        className="CopyText__btn CopyText__btn--copy"
        onClick={copyToClipboard}
      >
         <FontAwesomeIcon icon={faCopy} color="white" />
      </button>
      {
        isTooltipVisible
        && (
        <Tooltip
          isVisible={isTooltipVisible}
          text="Copied to clipboard"
        />
        )
      }

    </div>
  )
}


export default CopyText;
