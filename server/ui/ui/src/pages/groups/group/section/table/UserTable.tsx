// vendor
import React, { FC, useContext } from 'react';
// context
import AppContext from 'Src/AppContext'
// components
import UserItem from './item/UserItem';
import AddUser from '../add/AddUser';
// css
import './UserTable.scss';


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
  group: Group
}

const UserTable: FC<Props> = ({ group }: Props) => {
  // context
  const { user } = useContext(AppContext);


  return (
    <div>
      <table>
        <tbody>
          { (user.profile.role !== 'user') &&
            <AddUser />
          }
        </tbody>
      </table>
      <table className="UserTable">
        <thead>
          <tr>
            <th>Username</th>
            <th>Name</th>
            <th>Role</th>

            { (user.profile.role !== 'user') && (
              <th>Remove User</th>
            )}
          </tr>
        </thead>
        <tbody>
          {
            group.memberships && group.memberships.map((membership) => (
                <UserItem
                  key={membership.user.username}
                  userItem={membership.user}
                />
            ))
          }

          {
            !group.memberships && (
              <tr>
                <td
                  className="UserTable__td UserTable__td--empty"
                  colSpan={3}
                >
                  <p>Currently this group has no users assigned.</p>
                </td>
              </tr>
            )
          }
        </tbody>
      </table>
    </div>
  )
}

export default UserTable;
