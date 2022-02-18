// vendor
import React, { FC } from 'react';
// components
import Loader from 'Components/loader/index';
// css
import './Processing.scss';

interface Props {
  modalType: string;
  name: string;
}

const Processing: FC<Props> = ({ name, modalType }: Props) => {
  return (
    <div className="Processing">
      <h4>Creating new {modalType} {name}</h4>
      <Loader nested={true} />
    </div>
  )
}

export default Processing;
