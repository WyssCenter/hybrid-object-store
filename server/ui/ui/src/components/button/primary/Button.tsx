// vendor
import React, { FC, MouseEvent } from 'react';
// css
import './Button.scss';


interface Props {
  click: (event: MouseEvent) => void;
  disabled?: boolean;
  text: string,
}

const Button: FC<Props> = ({
  click,
  disabled = false,
  text,
}: Props) => {

  return (
    <button
      className="Button"
      disabled={disabled}
      onClick={click}
    >
      {text}
    </button>
  )
}


export default Button;
