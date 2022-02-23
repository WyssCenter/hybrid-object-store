// vendor
import React, { FC, useContext } from 'react';
// context
import AppContext from 'Src/AppContext';
// components
import UserTable from './table/UserTable';
// css
import './GroupSection.scss';


interface User {
  full_name: string;
  username: string;
  role: string;
}

interface Group {
  description: string;
  group_name: string;
  memberships: Array<Membership>;
}

interface Membership {
  user: User;
}

interface Props {
  group: Group;
}

const GroupSection: FC<Props> = ({ group }: Props) => {

  return (
    <div className="GroupSection flex">
      <div className="column-1-span-4 flex flex--column justify--center UserTable__text">
        <p>
          The table shows which users are members of this group. When a group is granted access to a dataset, all users in the group will receive the permissions granted to the group.
        </p>
        <p>
          For example, if a group is given read/write access to a dataset, all users in the group will be able to read and write data.
        </p>
        <p>
          To manage groups, you must have the &quot;admin&quot; or &quot;privileged&quot; role.
        </p>
      </div>
      <div>
        <UserTable group={group} />
      </div>
    </div>
  )
}


export default GroupSection;
