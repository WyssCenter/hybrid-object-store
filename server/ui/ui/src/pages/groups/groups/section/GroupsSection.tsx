// vendor
import React, { FC } from 'react';
// components
import GroupsTable from './table/GroupsTable';


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
}

const GroupsSection: FC<Props> = ({ user }: Props) => {
  const authorized = user.role !== 'user';

  return (
    <section>

      <GroupsTable
        user={user}
        authorized={authorized}
      />
    </section>
  )
}

export default GroupsSection;
