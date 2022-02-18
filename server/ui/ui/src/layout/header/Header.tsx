// vendor
import React, { FC } from 'react';
// components
import Logo from './logo/Logo';
import Toolbar from './toolbar/Toolbar';
// css
import './Header.scss';

interface Props {
  hasReactRouter: boolean;
  send: any;
  syncValue: any;
  serverName: string;
  datasetDescription: string;
}

const Header: FC<Props> = ({ hasReactRouter, send, serverName, syncValue, datasetDescription }:Props) => {
  return (
    <header className="Header">
      { hasReactRouter &&
        <Logo send={send} serverName={serverName} />
      }

      { hasReactRouter &&
        <Toolbar syncValue={syncValue} datasetDescription={datasetDescription} />
      }
    </header>
  )
}


export default Header;
