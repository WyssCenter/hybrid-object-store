// vendor
import React, { FC } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlug } from '@fortawesome/free-solid-svg-icons';

// css
import './SystemError.scss';



const SystemError: FC = () => {
  return (
    <div className="SystemError">
      <h2 className="SystemError__h2">System Error</h2>
      <p className="SystemError__p">Trouble finding the auth endpoint. Contact your system administrator for support.</p>

      <FontAwesomeIcon
        color="rgb(236, 128, 11)"
        icon={faPlug}
        size="6x"
      />
    </div>
  );
}


export default SystemError;
