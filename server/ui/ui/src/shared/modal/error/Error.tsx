// vendor
import React, { FC, useEffect, useContext } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faExclamationTriangle } from '@fortawesome/free-solid-svg-icons';
import AppContext from 'Src/AppContext';
// components
import Button from "Components/button/index";
// css
import './Error.scss';

interface Props {
  errorMessage: string;
  name: string;
  send: (arg: string) => void;
}

// TODO add error icon

const Error: FC<Props> = ({
  errorMessage,
  name,
  send
}:Props) => {
  const { user } = useContext(AppContext);
  useEffect(() => {
    if (errorMessage) {
      if(user && user.profile) {
        const { email } = user.profile;
        const redirectObject = {
          [email]: window.location.pathname,
        }
        const oldRedirectObject = localStorage.getItem('redirect');
        if (oldRedirectObject) {
          const newRedirectObject = Object.assign({}, JSON.parse(oldRedirectObject), redirectObject);
          localStorage.setItem('redirect', JSON.stringify(newRedirectObject))
        } else {
          localStorage.setItem('redirect', JSON.stringify(redirectObject))
        }
      }
    }
    send('AUTH')
  }, [user.profile, errorMessage, send]);


  return (
    <div className="Error">
      <h4>Error Creating {name}</h4>
      <FontAwesomeIcon
        color="#EC940B"
        icon={faExclamationTriangle}
        size="8x"
      />
      <p className="Error__p--error">{errorMessage}</p>

      <Button
        click={()=> send("TRY_AGAIN")}
        disabled={false}
        text="Try Again"
      />
    </div>
  )
}

export default Error;
