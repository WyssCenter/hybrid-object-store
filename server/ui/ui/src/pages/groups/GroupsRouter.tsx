// vendor
import React,
{
  FC,
} from 'react';
import { Switch, Route } from 'react-router-dom';
// components
import Groups from './groups/Groups';
import Group from './group/Group';

const GroupsRouter: FC = () => {

  return (
      <Switch>
          <Route
            exact
            path="/groups"
          >
            <Groups />
          </Route>

          <Route
            exact
            path="/groups/:groupname"
          >
            <Group />
          </Route>
      </Switch>
  )
}

export default GroupsRouter;
