// vendor
import React, { FC } from 'react';
// components
import Loader from 'Components/loader/index';
// css
import './Authenticating.scss';

const Authenticating: FC = () => {
  return (
    <div className="Authenticating">
      <h2>Authenticating</h2>
      <Loader nested />
    </div>
  );
}


export default Authenticating;
