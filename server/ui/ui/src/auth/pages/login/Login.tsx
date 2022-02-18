// vendor
import React, { FC } from 'react';
// css
import './Login.scss';

interface State {
  [key: string]: any;
}

interface Props {
  send: any;
  machineState: State;
}

const Login: FC<Props> = ({ machineState, send }: Props ) => {
  return (
    <div className="Login">
      <h4>Welcome to the Hoss!</h4>
      <br />
      <br />
      <br />
      <button
        onClick={() => send('AUTH')}
      >
        Login
      </button>
    </div>
  );
}


export default Login;
