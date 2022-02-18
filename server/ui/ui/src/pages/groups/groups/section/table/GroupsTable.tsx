// vendor
import React,
  {
    FC,
    MouseEvent,
    useState,
  } from 'react';
import { faPlus } from '@fortawesome/free-solid-svg-icons';
// components
import CreateGroupModal from '../create/CreateGroupModal';
import { FlatIconTextButton } from 'Components/button/index';
import GroupItem from './item/GroupItem';
import { DivisionText } from 'Components/division/index';
// css
import './GroupsTable.scss';


interface Group {
  description: string;
  group_name: string;
  memberships: null;
}

interface Membership {
  group: Group;
}


interface User {
  memberships: Array<Membership>;
  role: string;
}

interface Props {
  user: User;
  authorized: boolean;
}

const GroupsTable: FC<Props> = ({ user, authorized }: Props) => {
  // state
  const [isModalVisible, setIsModalVisible] = useState(false);

  // filter out users by hoss-default group
  const memberships = user.memberships.filter((membership) => (
    membership.group.group_name.indexOf('hoss-default-group') === -1
  ))

  return (
    <>
      { authorized && (
        <CreateGroupModal
          handleClose={() => setIsModalVisible(false)}
          isVisible={isModalVisible}
        />
      )}
      <table className="GroupsTable">
        <thead>
          <tr>
            <th>Group</th>
            <th>Description</th>
            { authorized && (
              <th>Actions</th>
            )}
          </tr>
        </thead>
        <tbody>
          { authorized && (
            <tr>
              <td
                className="td--center td--actions"
                colSpan={3}
              >
                <FlatIconTextButton
                  click={() => setIsModalVisible(true)}
                  color="primary"
                  icon={faPlus}
                  text="Create a Group"
                />
              </td>
            </tr>
          )}
          {
            memberships.map((membership) => (
                <GroupItem
                  key={membership.group.group_name}
                  membership={membership}
                  authorized={authorized}
                />
            ))
          }
          { memberships && memberships.length === 0 && (
            <td className="GroupsTable__td--empty" colSpan={3}>
              <p>You are currently not a member of any groups.</p>
            </td>
          )
          }
        </tbody>
      </table>
    </>
  )
}

export default GroupsTable;
