// vendor
import React, { FC } from 'react';
// Components
import Loader from 'Components/loader/index';
// css
import './Loading.scss';

const Loading: FC = () => {
  return (
    <div className="Loading">
      <Loader nested />
    </div>
  );
}


export default Loading;
