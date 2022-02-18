
import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';

// import reportWebVitals from './reportWebVitals';

import 'Styles/critical.scss';

// machine
import authMachine from './auth/machine/AuthStateMachine';
// TODO add a dev wrapper for this
// import('@xstate/inspect').then(({ inspect }) => {
//    inspect({
//      iframe: false
//    })
//  })






ReactDOM.render(
  <React.StrictMode>
    <App machine={authMachine} />
  </React.StrictMode>,
  document.getElementById('root')
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
// reportWebVitals(console.log);
