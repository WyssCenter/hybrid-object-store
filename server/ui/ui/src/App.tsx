// vendor
import React, { useEffect, useRef, useState, FC } from 'react';
import { useMachine, asEffect } from '@xstate/react';
import hexRgb from 'hex-rgb';
// context
import AppContext from './AppContext';
// components
import Routes from './pages/Routes';
import Layout from './layout/Layout'
import Login from './auth/pages/login/Login';
import Error from './auth/pages/error/Error';
import SystemError from './auth/pages/system/SystemError';
import Authenticating from './auth/pages/authenticating/Authenticating';
import Loading from './auth/pages/loading/Loading';
// css
import './App.scss';

interface RenderMap {
  [key: string]: JSX.Element | undefined;
  loggedIn?: JSX.Element;
  loggedOut?: JSX.Element;
  error?: JSX.Element;
  loading?: JSX.Element;
  logginIn?: JSX.Element;
  authenticating?: JSX.Element;
}

interface Props {
  machine: any;
}

const App: FC<Props> = ({ machine }: Props ) => {
  // state
  const [state, send] = useMachine(machine);
  const [serverName, setServerName] = useState('localhost');
  // vars
  const stateValue: string = typeof state.value === 'string' ? state.value : 'authenticating';
  const [ syncValue, setSyncValue ] = useState({ section: null, data: null });

  useEffect(() => {

  const link: HTMLLinkElement  = document.querySelector("link[rel~='icon']");
  if (link) {
    link.href = `${window.location.protocol}//${window.location.hostname}/ui/favicon.png`;
  }

  const baseHeaders = {
    "Access-Control-Allow-Origin": "*",
    "Content-Type": 'application/json',
    "Origin": window.location.origin
  }
   fetch(`${window.location.protocol}//${window.location.hostname}/ui/config.json`, { headers: baseHeaders, method: 'GET' })
   .then(res => res.json())
   .then((data) => {
     const { colors } = data;
     setServerName(data.server_name);
     const rgb = hexRgb(colors.primary);
     document.title = `Hoss - ${data.server_name}`
      const root = document.documentElement;
      root.style.setProperty('--main-color', colors.primary);
      root.style.setProperty('--secondary-color', colors.secondary);
      root.style.setProperty('--main-rgb', `${rgb.red},${rgb.green},${rgb.blue}`)
    })
  }, [])

 useEffect(() => {
   if (stateValue === 'idle') {
     send("FETCH_CONFIG")
   }
 }, [send, stateValue]);


 const renderMap: RenderMap = {
   idle: (null),
   loggedIn: (
     <Routes
      send={send}
      syncValue={syncValue}
      setSyncValue={setSyncValue}
      serverName={serverName}
     />
   ),
   loggedOut: (
     <Login
       machineState={state}
       send={send}
     />
   ),
   error: (
     <Error
      context={state.context}
      send={send}
     />
    ),
   loading: (<Loading />),
   logginIn:(<Authenticating />),
   authenticating: (<Authenticating />),
   systemError: (<SystemError />),
   redirect: (<Loading />)
 };


 if (stateValue === 'loggedIn') {
   return (
     <AppContext.Provider value={state.context}>
      {renderMap[stateValue]}
     </AppContext.Provider>
   );
 }

 return (
   <AppContext.Provider value={state.context}>
     <Layout
      serverName={serverName}
      datasetDescription={''}
      send={() => {return}}
      hasReactRouter={false}
      syncValue={syncValue}
      >
      <div className="grid text-center">
        {renderMap[stateValue]}
      </div>
     </Layout>
   </AppContext.Provider>
 );
};



export default App;
