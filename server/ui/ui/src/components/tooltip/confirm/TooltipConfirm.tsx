// vendor
import React, { FC, useRef, MutableRefObject } from 'react';
// Hooks
import { useEventListener,  } from 'Hooks/events/index';

// components
import { PrimaryButton, TernaryButton } from '../../button/index';
// css
import './TooltipConfirm.scss';

interface Props {
  cancel: () => void;
  confirm: () => void;
  isVisible: boolean;
  text: string;
  tooltipRef: MutableRefObject<HTMLElement>;
}

const TooltipConfirm: FC<Props> = ({
  cancel,
  confirm,
  isVisible,
  text,
  tooltipRef,
}:Props) => {
  /**
  * Method provides a way for child componts to update state
  * @param {Object} evt
  * @fires setMenuVisibility
  */
  const windowClickHandler = (event: Event) => {
    const element = event.target as HTMLInputElement;
    if (tooltipRef.current && !tooltipRef.current.contains(element)) {
      cancel();
    }
  }

  // Add event listener using our hook
  useEventListener('click', windowClickHandler);

  if (!isVisible) {
    return null;
  }
  return (
    <div
      className="TooltipConfirm"
    >
      <div className="TooltipConfirm__pointer" />
      <div className="TooltipConfirm__menu">
        <p>{text}</p>
        <div className="flex justify--space-around">
          <TernaryButton
            click={cancel}
            text="Cancel"
          />
          <PrimaryButton
            click={confirm}
            text="Confirm"
          />
        </div>
      </div>
    </div>
  )
}

export default TooltipConfirm;
