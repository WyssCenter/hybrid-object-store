// vendor
import React, {
  FC,
  useContext,
  } from 'react';
import classNames from 'classnames';
// components
import AddSection from '../shared/add/AddSection';
import PermissionsItem from './item/PermissionsItem';
// context
import AppContext from 'Src/AppContext';
// css
import './PermissionsSection.scss'
// types
import {
  Permission,
} from '../DatasetSectionTypes';




interface Props {
  list: Array<Permission>;
  sectionType: string;
  headerText: string;
}


const PermissionsSection: FC<Props> = ({ list, sectionType, headerText }: Props) => {

  // context
  const { user } = useContext(AppContext);

  const authorized = user.profile.role !== 'user';

  const nameCSS = classNames({
    PermissionsSection__name: true,
    "PermissionsSection__name--unauthorized": !authorized,
  })

  return (
    <div className="PermissionsSection">
      <h3 className="DatasetSection__h3">{headerText}</h3>
      {
        authorized && (
          <AddSection
            sectionType={sectionType}
          />
        )
      }
      <table className="PermissionsSection__table">
        <thead>
          <tr>
            <th className={nameCSS}>
              {sectionType}
            </th>
            <th>
              Permissions
            </th>
            {
              authorized && (
                <>
                  <th>
                  </th>
                  <th>
                    Actions
                  </th>
                </>
              )
            }
          </tr>
        </thead>
        <tbody className="PermissionsSection__body">
          {
            list.map((item) => {

              return (
                <PermissionsItem
                  key={item.group.group_name}
                  item={item}
                  sectionType={sectionType}
                />
              )
            })
          }

        </tbody>
      </table>
    </div>
  );
}


export default PermissionsSection;
