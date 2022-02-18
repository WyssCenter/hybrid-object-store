// vendor
import React,
  {
    FC,
  } from 'react';
import {
  Link,
  useLocation,
} from "react-router-dom";
import { useMachine } from '@xstate/react';
import classNames from 'classnames';
// enivornment
import { get } from 'Environment/createEnvironment';
// components
import BreadCrumbs from './breadcrumbs/Breadcrumbs';
import SyncStatus from './syncStatus/SyncStatus';
// css
import { ReactComponent as Gear } from 'Images/icons/icon-gear.svg';
import './Toolbar.scss';

interface Props {
  syncValue: any;
  datasetDescription: string;
}

const Toolbar: FC<Props> = ({ syncValue, datasetDescription }:Props) => {
  const location = useLocation();

  const paths = location.pathname.split('/')
  const namespace = paths.length > 0 ? paths[1] : '';
  const settings = ['account', 'groups', 'tokens'];
  const isSettings = settings.includes(namespace);
  const datasetName = paths.length > 1 ? paths[2] : false;

  const showSettings = (namespace !== '')
    && !isSettings;
  let settingsLink = datasetName
    ? `/${namespace}/${datasetName}/settings`
    : `/${namespace}/settings`;

  const settingsSelected = (paths[2] === 'settings')
    || paths[3] === 'settings';

  if (settingsSelected) {
    settingsLink = '#';
  }

  const settingsCSS = classNames({
    'column-1-span-1 Toolbar__settings': true,
    'Toolbar__settings--selected': settingsSelected
  });

  return (
    <div className="Toolbar">
      <div className="column-1-span-1"></div>
      <BreadCrumbs datasetDescription={datasetDescription}/>
      {
        showSettings && (
          <SyncStatus
            sync={syncValue}
          />
        )
      }
      {
        showSettings && (
          <Link
          to={settingsLink}
          className={settingsCSS}>
            <Gear />
            Settings
          </Link>
        )
      }
    </div>
  )
}


export default Toolbar;
