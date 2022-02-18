// vendor
import React, { FC } from 'react';
// components
import Loader from 'Components/loader/index';
// css
import './Loading.scss';

const Loading: FC = () => {
  return (
    <div className="Loading">
      <h6>Fetching Namespaces</h6>
      <Loader nested={true} />
    </div>
  )
}

export default Loading;
