// vendor
import React, { FC } from 'react';
import { IconDefinition }from '@fortawesome/fontawesome-common-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faExclamationTriangle } from '@fortawesome/free-solid-svg-icons';
// css
import './Error.scss';


interface Props {
  errorMessage: string;
}


const Error: FC<Props> = ({ errorMessage }: Props) => {
  return (
    <div className="Error">
      <FontAwesomeIcon icon={faExclamationTriangle} color="orange" size="6x" />
      <h5 className="Error__h5">{errorMessage}</h5>
    </div>
  );
}


export default Error;
