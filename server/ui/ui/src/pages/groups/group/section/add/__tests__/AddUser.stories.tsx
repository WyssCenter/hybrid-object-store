// vendor
import React from 'react';
// context
import GroupContext from '../../../GroupContext'
// component
import AddUser from '../AddUser';

export default {
  title: 'Pages/groups/group/section/add/AddUser',
  component: AddUser,
  argTypes: {},
};

const send = () => null;

const groupname = 'group-1';

const Template = (args) => (
  <GroupContext.Provider value={{ send, groupname }}>
    <AddUser />
  </GroupContext.Provider>
);

export const AddUserComp = Template.bind({});

AddUserComp.args = {};

AddUserComp.parameters = {
  jest: ['AddUser.test.js'],
};
