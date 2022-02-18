// @flow
// vendor
import React, { FC } from 'react';
import classNames from 'classnames';
// css
import './Loader.scss';

interface Props {
  nested: boolean;
}

const Loader: FC<Props> = ({
  nested,
}: Props) => {
  const loaderCSS = classNames({
    Loader: !nested,
    'Loader--nested': nested,
  });

  return (
    <div
      className={loaderCSS}
      data-testid="loader"
    />
  );
};

export default Loader;
