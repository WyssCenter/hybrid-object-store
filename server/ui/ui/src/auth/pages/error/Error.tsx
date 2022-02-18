// vendor
import React, { FC, useCallback, useEffect } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPoo } from '@fortawesome/free-solid-svg-icons';
// components
import { SecondaryButton } from 'Components/button/index';
// css
import './Error.scss';

interface Data {
  message: string;
}

interface Error {
  data: Data;
}

interface Context {
  error?: Error;
}

interface Props {
  context?: any;
  send: any;
}

const Error: FC<Props> = ({ context, send }:Props) => {
  const errorMessage = context && context.error && context.error.data
   ? context.error.data.message
   : 'Unknown error occured';

  /**
  * Method triggers state machine TRY_AGAIN event and sets machine to loggedOut state
  * @param {}
  * @return {Void}
  * @calls {machine#send}
  */
  const tryAgain  = useCallback(()=> {
    send("TRY_AGAIN");
  }, [send])

  useEffect(() => {
    if (errorMessage) {
      if (context.user && context.user.profile) {
        const { email } = context.user.profile;
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
      send('AUTH')
    }
  }, [context.user, errorMessage, tryAgain]);

  return (
    <div className="Error">
      <h2 className="Error__h2">Authentication Error</h2>
      { errorMessage &&
        <p className="Error__p">{errorMessage}</p>
      }
      <FontAwesomeIcon
        color="rgb(236, 128, 11)"
        icon={faPoo}
        size="6x"
      />
      <br />
      <SecondaryButton
        click={tryAgain}
        text="Try Again"
      />


    </div>
  );
}


export default Error;
