// vendor
import React, { FC } from 'react';
// css
import './AccountLabel.scss';

interface Props {
  label: string;
  value: string;
}

const Account: FC<Props> = ({
  label,
  value,
}: Props) => {
  if (value === '*') {
    return null;
  }

  return (
    <div className="AccountLabel">
      <h6 className="AccountLabel__h6">
        {`${label}: `}
      </h6>
      <p className="AccountLabel__value">
        {value}
      </p>
    </div>
  )
}

export default Account;
