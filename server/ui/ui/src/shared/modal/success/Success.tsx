// vendor
import React, { FC } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCheckCircle } from '@fortawesome/free-solid-svg-icons';
// css
import './Success.scss';

interface Props {
  modalType: string;
  name: string;
}
// TODO add checkmark icon to body
const Success: FC<Props> = ({ name, modalType }: Props) => {
  const root = document.documentElement;
  const primaryHash = root.style.cssText ?
    root.style.cssText.split(';')[0].split(':')[1]
    : '#957299';
  return (
    <div className="Success">
      <h4>{modalType} {name} was successfully created</h4>

      <FontAwesomeIcon
        color={primaryHash}
        icon={faCheckCircle}
        size="8x"
      />
    </div>
  )
}

export default Success;
