// vendor
import React, { FC, MouseEvent } from 'react';
// css
import './WarningButton.scss';


interface Props {
  click: (event: MouseEvent) => void;
  disabled?: boolean;
  text: string,
}

const WarningButton: FC<Props> = ({
  click,
  disabled = false,
  text,
}: Props) => {

  return (
    <button
      className="WarningButton"
      disabled={disabled}
      onClick={click}
    >
      {text}
    </button>
  )
}


export default WarningButton;
