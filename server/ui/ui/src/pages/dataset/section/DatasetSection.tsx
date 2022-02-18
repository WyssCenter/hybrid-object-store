// vendor
import React, { FC, useMemo } from 'react';
// components
import { SectionCard } from 'Components/card/index';
import PermissionsSection from './permissions/PermissionsSection';
// types
import {
  Dataset,
  FetchPermissionItem,
  Permission,
} from './DatasetSectionTypes';
// css
import './DatasetSection.scss';


interface Props {
  dataset: Dataset;
  datasetName: string;
  namespace: string;
}

/**
* Method sorts permissions into users and groups
* @param {Array<Object>} string
* @return {Object: {<Array>, <Array>}}
*/
const sortPermissions = (permissions: Array<FetchPermissionItem>) => {
  const groups: Array<Permission> = [];
  const users: Array<Permission> = [];

  permissions.forEach((permission) => {
    const item = { ...permission, permission: { name: permission.permission}}
    if(permission.group.group_name.indexOf('hoss-default-group') > -1) {
      users.push(item);
    } else {
      groups.push(item);
    }
  });

  return { groups, users };
}


const DatasetSection: FC<Props> = ({
  dataset,
}:Props) => {

  const {
    permissions,
  } = dataset;

  const { users, groups } = useMemo(() => sortPermissions(permissions), [permissions]);

  return (
    <div className="DatasetSection">
        <h3 className="DatasetSection__h3">Permissions</h3>
        <SectionCard verticalHeight="grid-v-auto">
          <div className="flex">
            <PermissionsSection
              list={users}
              sectionType="user"
              headerText="Users"
            />
            <PermissionsSection
              list={groups}
              sectionType="group"
              headerText="Groups"
            />
          </div>
        </SectionCard>
    </div>
  )
}

export default DatasetSection;
