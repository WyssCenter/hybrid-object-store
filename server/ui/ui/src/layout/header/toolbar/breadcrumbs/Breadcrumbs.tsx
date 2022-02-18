// vendor
import React, { FC } from 'react';
import {
  Link,
  useLocation,
} from "react-router-dom";
// css
import { ReactComponent as Server } from 'Images/icons/icon-server.svg';
import { ReactComponent as SectionDivider } from 'Images/icons/arrows/icon-nav-chevron.svg';
import './Breadcrumbs.scss';

interface settingsMap {
   [key: string]: string;
}

interface Props {
  datasetDescription: string;
}

const Breadcrumbs: FC<Props> = ({ datasetDescription }: Props) => {
  const location = useLocation();
  const paths = location.pathname.split('/')
  const namespace = paths.length > 0 ? paths[1] : '';
  const settings = ['account', 'groups', 'tokens'];
  const settingsMap:settingsMap =  {
    account: 'Account',
    groups: 'Groups',
    tokens: 'Personal Access Tokens',
  };
  const isSettings = settings.includes(namespace);
  const namespaceHeader = isSettings ? settingsMap[namespace] : namespace;
  const namespaceSubheader = isSettings ? '' : 'Namespace';
  let datasetSubheader = isSettings ? 'Group' : 'Dataset';
  const datasetName = paths.length > 1 ? paths[2] : false;
  let datasetHeader = datasetName;
  const datasetSettings = (paths[3] === 'settings')
  if (datasetName === 'settings') {
    datasetSubheader = '';
    datasetHeader = 'Settings';
  }
  return (
    <div className="Breadcrumbs">
      <Link
        to = "/"
      >
        <Server />
      </Link>
      {
        !(namespace || datasetName) && (
        <>
          <div className="Breadcrumbs__divider">
            <SectionDivider />
          </div>
          <div className="Breadcrumbs__name">
            <div className="Breadcrumbs__title">
              Available Namespaces
            </div>
            <div className="Breadcrumbs__sub-title">
            </div>
          </div>
        </>
        )
      }

      {  namespace && (
          <>
            <div className="Breadcrumbs__divider">
              <SectionDivider />
            </div>
            <div className="Breadcrumbs__name">
              <div className="Breadcrumbs__title">
                <Link className="Breadcrumbs__link" to={`/${namespace}`}>
                  {namespaceHeader}
                </Link>
              </div>
              <div className="Breadcrumbs__sub-title">
                {namespaceSubheader}
              </div>
            </div>
          </>
        )
      }

      {  datasetName && (
          <>
            <div className="Breadcrumbs__divider">
              <SectionDivider />
            </div>
            <div className="Breadcrumbs__name">
              <div className="Breadcrumbs__title">
                <Link className="Breadcrumbs__link" to={`/${namespace}/${datasetName}`}>
                  {datasetHeader}
                </Link>
              </div>
              <div className="Breadcrumbs__sub-title">
                {datasetSubheader}
              </div>
            </div>
          </>
        )
      }
      {
        datasetName && !datasetSettings && (
          <>
            <div className="Breadcrumbs__description">
              {datasetDescription}
            </div>
          </>
        )
      }
      {  datasetSettings && (
          <>
            <div className="Breadcrumbs__divider">
              <SectionDivider />
            </div>
            <div className="Breadcrumbs__name">
              <div className="Breadcrumbs__title">
                Settings
              </div>
              <div className="Breadcrumbs__sub-title">
              </div>
            </div>
          </>
        )
      }
    </div>
  )
}


export default Breadcrumbs;
