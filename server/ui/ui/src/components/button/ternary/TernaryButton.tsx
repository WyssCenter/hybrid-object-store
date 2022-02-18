// vendor
import React, { FC } from 'react';
// css
import './TernaryButton.scss';


interface Props {
  click: () => void;
  disabled?: boolean;
  text: string,
}

const defaultProps = {
  disabled: false
};

const TernaryButton: FC<Props> = ({
  click,
  disabled = false,
  text,
}: Props) => {

  return (
    <button
      className="TernaryButton"
      disabled={disabled}
      onClick={click}
    >
      {text}
    </button>
  )
}

TernaryButton.defaultProps = defaultProps;


export default TernaryButton;
