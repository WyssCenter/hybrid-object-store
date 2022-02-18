// vendor
import React, { FC } from 'react';
// css
import './SecondaryButton.scss';


interface Props {
  click: () => void;
  disabled?: boolean;
  text: string,
}

const SecondaryButton: FC<Props> = ({
  click,
  disabled = false,
  text,
}: Props) => {

  return (
    <button
      className="SecondaryButton"
      disabled={disabled}
      onClick={click}
    >
      {text}
    </button>
  )
}


export default SecondaryButton;
