// vendor
import React, { FC, ReactNode } from 'react';
// components
import Header from './header/Header';
// css
import './Layout.scss';



interface Props {
  children: ReactNode;
  hasReactRouter: boolean;
  send: any;
  syncValue: any;
  serverName: string;
  datasetDescription: string;
}

const Layout: FC<Props> = ({ children, hasReactRouter, send, serverName, syncValue, datasetDescription }: Props) => {
  return (
    <div className="Layout">
      <Header
        hasReactRouter={hasReactRouter}
        send={send}
        syncValue={syncValue}
        serverName={serverName}
        datasetDescription={datasetDescription}
      />

      <main>
        { children }
      </main>
    </div>
  )
}


export default Layout;
