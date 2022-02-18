// vendor
import React, { FC, useState } from 'react';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
// layout
import Layout from '../layout/Layout';
// pages
import Account from './account/Account';
import GroupsRouter from './groups/GroupsRouter';
import PersonalAccessTokens from './tokens/PersonalAccessTokens';
import NamespaceListing from './namespaces/NamespaceListing';
import Namespace from './namespace/Namespace';
import Dataset from './dataset/Dataset';

interface Props {
  send: any;
  setSyncValue: any;
  syncValue: any;
  serverName: string;
}

const Routes: FC<Props> = ({ send, serverName }:Props) => {
  const [syncValue, setSyncValue] = useState(null);
  const [ datasetDescription, setDatasetDescription ] = useState(null);
  return (

    <BrowserRouter basename="ui">
      <Layout
        hasReactRouter={true}
        send={send}
        syncValue={syncValue}
        datasetDescription={datasetDescription}
        serverName={serverName}
        >
       <Switch>

          <Route
            exact
            path="/account"
          >
            <Account />
          </Route>


          <Route
            exact
            path={["/groups", "/groups/:groupName"]}
          >
            <GroupsRouter />
          </Route>


          <Route
            exact
            path="/tokens"
          >
            <PersonalAccessTokens />
          </Route>

          <Route
            exact
            path="/"
          >
            <NamespaceListing />
          </Route>

          <Route
            exact
            path={["/:namespace", "/:namespace/settings"]}
          >
            <Namespace
              setSyncValue={setSyncValue}
            />
          </Route>

          <Route
            exact
            path={["/:namespace/:datasetName", "/:namespace/:datasetName/settings"]}
            component={Dataset}
          >
            <Dataset
              setSyncValue={setSyncValue}
              setDatasetDescription={setDatasetDescription}
            />
          </Route>

       </Switch>
      </Layout>
   </BrowserRouter>
  )
}


export default Routes;
